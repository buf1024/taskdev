package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"taskdev/util"
	"time"

	myproto "taskdev/proto"

	"sync"

	mylog "github.com/buf1024/golib/logging"
	mynet "github.com/buf1024/golib/net"
	"github.com/golang/protobuf/proto"
)

// ServerTaskHandler the server handler
type scheHandler func(conn *connState, msg *myproto.Message)

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

	connQueue []*connState
	queueLock sync.Locker

	sessChan  chan struct{}
	timerChan chan struct{}

	cmdChan  chan string
	exitChan chan struct{}

	host   string
	logDir string
	debug  bool
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

func (s *scheServer) addHandler() {

	//	b.handler[myproto.KCmdHeartBeatReq] = b.handleHeartBeatReq
	//	b.handler[myproto.KCmdRegisterRsp] = b.handleRegisterRsp
}

func (s *scheServer) handle(conn *connState, p *myproto.Message) error {
	handler, ok := s.handler[p.Head.Command]
	if !ok {
		return fmt.Errorf("command 0x%x handler not found", p.Head.Command)
	}
	go handler(conn, p)

	return nil
}

func (s *scheServer) hearbeat() {
END:
	for {
		t := time.After(time.Second * 30)
		select {
		case <-t:
			{
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
			fmt.Printf("poll event error!")
			s.cmdChan <- "stop"
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
	ch := make(chan os.Signal)
	signal.Notify(ch, syscall.SIGHUP, syscall.SIGTERM, syscall.SIGUSR2, syscall.SIGUSR1)
END:
	for {
		select {
		case sig := <-ch:
			switch sig {
			case syscall.SIGHUP:
				s.log.Info("catch SIGHUP, stop the server.")
				s.cmdChan <- "stop"
				break END
			case syscall.SIGTERM:
				s.log.Info("catch SIGTERM, stop the server.")
				s.cmdChan <- "stop"
				break END
			case syscall.SIGUSR1:
				s.log.Info("catch SIGUSR1.")
			case syscall.SIGUSR2:
				s.log.Info("catch SIGUSR2.")
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

	dir := filepath.Dir(os.Args[0])
	dir = dir + string(filepath.Separator) + "logs" + string(filepath.Separator)

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

	s.net = mynet.NewSimpleNet()
	s.proto = myproto.NewProto()
	s.queueLock = &sync.Mutex{}

	s.cmdChan = make(chan string, 1024)
	s.exitChan = make(chan struct{})

	go s.sigTask()
	go s.eventTask()
	go s.hearbeat()

	s.log.Info("start command task go routine\n")
	go s.cmdTask()

	return nil
}

func (s *scheServer) server() {

	s.cmdChan <- "listen"

	s.wait()
}

func (s *scheServer) stop() {
	close(s.cmdChan)
	mynet.SimpleNetDestroy(s.net)

	s.log.Stop()

	s.exitChan <- struct{}{}
}

func (s *scheServer) wait() {
	<-s.exitChan
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
