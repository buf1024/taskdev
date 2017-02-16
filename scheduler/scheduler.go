package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"path/filepath"
	"sync"
	"syscall"
	"taskdev/util"
	"time"

	myproto "taskdev/proto"

	mylog "github.com/buf1024/golib/logging"
	mynet "github.com/buf1024/golib/net"
	"github.com/golang/protobuf/proto"
)

// ServerTaskHandler the server handler
type scheHandler func(conn *connState, msg *myproto.Message)

type sessData struct {
	conn *connState
	msg  *myproto.Message
}

const (
	buildDaily = iota
	buildSomeday
	buildToday
	buildNow
)

type buildScheData struct {
	scheType int
	someday  int
	daySecs  int
	stop     bool

	conn  *connState
	task  *util.Task
	tasks *util.TaskGroup
}

type connState struct {
	conn      *mynet.Connection
	conStatus bool
	regStatus bool
}

type scheServer struct {
	net     *mynet.SimpleNet
	proto   *myproto.TaskProto
	listen  *mynet.Listener
	log     *mylog.Log
	handler map[uint32]scheHandler
	sess    *util.Session

	connQueue []*connState
	queueLock sync.Locker

	cmdChan       chan string
	exitChan      chan struct{}
	buildScheChan chan *buildScheData

	host   string
	logDir string
	debug  bool

	running bool
}

func (s *scheServer) addHandler() {

	s.handler[myproto.KCmdHeartBeatReq] = s.handleHeartBeatReq
	s.handler[myproto.KCmdRegisterReq] = s.handleRegisterReq
	s.handler[myproto.KCmdTaskBuildRsp] = s.handleTaskBuildRsp
	s.handler[myproto.KCmdTaskStateRsp] = s.handleTaskStateRsp

}

func (s *scheServer) handle(conn *connState, p *myproto.Message) error {
	s.log.Info("handle:\n%s\n", myproto.Debug(p))
	handler, ok := s.handler[p.Head.Command]
	if !ok {
		return fmt.Errorf("command 0x%x handler not found", p.Head.Command)
	}
	go handler(conn, p)

	return nil
}
func (s *scheServer) add(con *connState) {
	s.queueLock.Lock()
	defer s.queueLock.Unlock()

	s.connQueue = append(s.connQueue, con)
}
func (s *scheServer) del(con *connState) {
	for i, v := range s.connQueue {
		if con == v {
			if i == len(s.connQueue)-1 {
				s.connQueue = s.connQueue[:i]
			} else {
				s.connQueue = append(s.connQueue[:i], s.connQueue[i+1:]...)
			}
			break
		}
	}
}

func (s *scheServer) sessionCheck() {
	for {
		t := time.NewTimer((time.Duration)((int64)(time.Second) * 30))
		<-t.C
		s.log.Debug("session check\n")
		to := s.sess.GetTimeout(60)
		for k, v := range to {
			s.log.Warning("timeout sid = %s\n", k)
			go s.handleTimeout(k, v)
		}
	}

}

func (s *scheServer) hearbeat() {
END:
	for {
		t := time.After(time.Second * 30)
		select {
		case <-t:
			{
				s.log.Debug("timer check hearbeat\n")
				if s.connQueue == nil || len(s.connQueue) == 0 {
					continue END
				}
				n := time.Now()
				for _, v := range s.connQueue {
					conn := v.conn
					if conn != nil && conn.Status() != mynet.StatusBroken {
						up := conn.UpdateTime()
						diff := n.Unix() - up.Unix()
						if diff > 30 {
							if !v.regStatus {
								// 30秒没注册，主动断开
								s.log.Info("not reigster for 30s, disconnect, remote = %s\n", conn.RemoteAddress())
								s.net.CloseConn(conn)
								s.del(v)
								continue
							}
							// 心跳
							m := &myproto.Message{}
							m.Head.Command = myproto.KCmdHeartBeatReq
							// beartbeat
							p, err := myproto.GetMessage(m.Head.Command)
							if err != nil {
								s.log.Error("gen hearbeat failed, err = %s\n", err)
								continue
							}
							pheart := p.(*myproto.HeartBeatReq)
							pheart.SID = proto.String(util.GetSID(16))
							m.Body = pheart

							s.log.Info("REQ：\n%s\n", myproto.Debug(m))
							err = s.net.SendData(v.conn, m)
							if err != nil {
								s.log.Error("send hearbeat, err = %s\n", err)
								continue
							}
						}

					}
				}
			}

		}
	}

}

func (s *scheServer) eventTask() {
	n := s.net
	for {
		evt, err := n.PollEvent(1000 * 60)
		if err != nil {
			s.log.Error("poll event error!\n")
			_, ok := <-s.cmdChan
			if ok {
				s.cmdChan <- "stop"
			}
			return
		}
		conn := evt.Conn
		switch {
		case evt.EventType == mynet.EventConnectionError:
			{
				s.log.Info("event error: local = %s, remote = %s\n",
					conn.LocalAddress(), conn.RemoteAddress())
				myconn := conn.UserData.(*connState)
				s.del(myconn)
			}
		case evt.EventType == mynet.EventConnectionClosed:
			{
				s.log.Info("event close: local = %s, remote = %s\n",
					conn.LocalAddress(), conn.RemoteAddress())
				myconn := conn.UserData.(*connState)
				s.del(myconn)
			}
		case evt.EventType == mynet.EventNewConnectionData:
			{
				data := evt.Data.(*myproto.Message)
				userData := conn.UserData.(*connState)
				err := s.handle(userData, data)
				if err != nil {
					s.log.Error("handle message error, err = %s\n", err)
				}
			}
		case evt.EventType == mynet.EventNewConnection:
			{
				s.log.Info("client connected: local = %s, remote = %s\n",
					conn.LocalAddress(), conn.RemoteAddress())
				myconn := &connState{
					conn: conn,
				}
				conn.UserData = myconn
				s.add(myconn)
			}
		case evt.EventType == mynet.EventTimeout:
			{
				s.log.Info("poll timeout now = %d\n", time.Now().Unix())
			}
		}
	}
}

func (s *scheServer) cmdTask() {
	for {
		cmd := <-s.cmdChan
		switch cmd {
		case "listen":
			{
				s.log.Info("start to listen, addr = %s\n", s.host)
				listen, err := s.net.Listen(s.host, s.proto)
				if err != nil {
					s.log.Critical("listen failed, err = %s\n", err)
					s.cmdChan <- "stop"
				} else {
					s.log.Info("listen success, waiting for connection\n")
				}
				s.listen = listen
			}
		case "stop":
			{
				s.log.Info("command chan receive stop command\n")
				s.stop()
				return
			}
		}
	}
}

func (s *scheServer) sigTask() {
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGQUIT, syscall.SIGHUP, syscall.SIGTERM, syscall.SIGUSR2, syscall.SIGUSR1)
END:
	for {
		select {
		case sig := <-ch:
			switch sig {
			case syscall.SIGINT:
				s.log.Info("catch SIGINT, stop the server.\n")
				//s.cmdChan <- "stop"
				//break END
				t := &util.Task{
					Name:            "bank",
					Vcs:             util.KVCSSvn,
					Version:         "129624",
					Repos:           "https://10.200.1.29:20443/JRTZXM/01-%E5%B9%BF%E8%B4%B5%E8%BF%90%E8%90%A5/02-%E4%BA%8C%E6%9C%9F%E4%BA%A4%E6%98%93%E7%B3%BB%E7%BB%9F/02-%E4%BB%A3%E7%A0%81/08-%E6%B8%85%E7%AE%97%E4%B8%AD%E5%BF%83%E4%BB%A3%E7%A0%81/03-%E9%93%B6%E8%A1%8C%E6%8E%A5%E5%8F%A3(%E6%B8%85%E7%AE%97%E4%B8%AD%E5%BF%83)/01-%E9%93%B6%E8%A1%8C%E6%9C%8D%E5%8A%A1/bank_svc",
					ReposDir:        "src/bankcc_svronline",
					User:            "luoguochun",
					Pass:            "virtual2.",
					PreBuildScript:  "#!/bin/sh\nls\necho PreBuildScript",
					BuildScript:     "#!/bin/sh\npwd\nautoconf\nconfigure\ncd tsbase\nmake install",
					PostBuildScript: "#!/bin/sh\nls\necho PostBuildScript",
				}
				b := &buildScheData{}
				b.conn = s.connQueue[0]
				b.task = t

				s.addBuildTask(b)

			case syscall.SIGQUIT:
				s.log.Info("catch SIGQUIT, stop the server.\n")
				s.cmdChan <- "stop"
				break END
			case syscall.SIGHUP:
				s.log.Info("catch SIGHUP, stop the server.\n")
				s.cmdChan <- "stop"
				break END
			case syscall.SIGTERM:
				s.log.Info("catch SIGTERM, stop the server.\n")
				s.cmdChan <- "stop"
				break END
			case syscall.SIGUSR1:
				s.log.Info("catch SIGUSR1.\n")
			case syscall.SIGUSR2:
				s.log.Info("catch SIGUSR2.\n")
				s.log.Sync()

			}
		}

	}

}

// initLog Setup Logger
func (s *scheServer) initLog(prefix string) (err error) {
	s.log, err = mylog.NewLogging()
	if err != nil {
		return err
	}
	level := mylog.LevelAll
	if !s.debug {
		level = mylog.LevelInformational
	}
	_, err = mylog.SetupLog("file",
		fmt.Sprintf(`{"prefix":"%s", "filedir":"%s", "level":%d, "switchsize":%d, "switchtime":%d}`,
			prefix, s.logDir, level, -1, 0))
	if err != nil {
		return err
	}
	_, err = mylog.SetupLog("console",
		fmt.Sprintf(`{"level":%d}`, level))
	if err != nil {
		return err
	}
	s.log.Start()

	s.log.Info("log started.\n")

	return nil
}

// parseArgs parse the commandline arguments
func (s *scheServer) parseArgs() {
	host := flag.String("host", "", "the server listen address")
	debug := flag.Bool("debug", false, "run in debug model")
	help := flag.Bool("help", false, "show the help message")

	flag.Parse()
	if *help {
		flag.PrintDefaults()
		os.Exit(0)
	}
	if len(*host) == 0 {
		fmt.Printf("listen host is empty. Usage:\n")
		flag.PrintDefaults()
		os.Exit(-1)
	}

	pwd, _ := filepath.Abs(filepath.Dir(os.Args[0]))
	dir := pwd + string(filepath.Separator) + "logs" + string(filepath.Separator)

	if _, err := os.Stat(dir); err != nil {
		if !os.IsNotExist(err) {
			fmt.Printf("log dir %s don't have permission\n", dir)
			os.Exit(-1)
		}
		os.Mkdir(dir, 0777)
	}

	s.host = *host
	s.debug = *debug
	s.logDir = dir

}

func (s *scheServer) init() error {
	s.parseArgs()

	err := s.initLog("scheduler")
	if err != nil {
		return fmt.Errorf("init log failed, err = %s", err)
	}

	s.net = mynet.NewSimpleNet(s.log)
	s.proto = myproto.NewProto()
	s.queueLock = &sync.Mutex{}
	s.handler = make(map[uint32]scheHandler)
	s.sess = util.NewSession()

	s.cmdChan = make(chan string, 1024)
	s.exitChan = make(chan struct{})
	s.buildScheChan = make(chan *buildScheData, 1024)

	s.addHandler()

	go s.sigTask()
	go s.eventTask()
	go s.hearbeat()
	go s.sessionCheck()
	go s.schedule()

	s.log.Info("start command task go routine\n")
	go s.cmdTask()

	return nil
}

func (s *scheServer) server() {

	s.cmdChan <- "listen"
	s.running = true

	s.wait()
}

func (s *scheServer) stop() {
	if s.running {

		close(s.cmdChan)
		close(s.buildScheChan)
		mynet.SimpleNetDestroy(s.net)

		//		t := time.NewTimer((time.Duration)((int64)(time.Second) * 5))
		//		<-t.C

		s.log.Stop()

		s.exitChan <- struct{}{}
	}
}

func (s *scheServer) wait() {
	<-s.exitChan
	s.running = false
}

func main() {
	s := &scheServer{}
	err := s.init()
	if err != nil {
		fmt.Printf("server init failed\n")
		os.Exit(-1)
	}
	s.server()

}
