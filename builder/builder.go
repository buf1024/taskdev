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

type builder struct {
	net      *mynet.SimpleNet
	proto    *myproto.TaskProto
	log      *mylog.Log
	handler  map[uint32]runnerHandler
	conn     *connState
	hostInfo *util.HostInfo

	sess      *util.Session
	timerChan chan struct{}
	cmdChan   chan string
	exitChan  chan struct{}

	buildQueue     map[string]string
	buildQueueLock sync.Locker

	host     string
	logDir   string
	debug    bool
	buildDir string
	dbFile   string
	db       *buildDb
}

type runnerHandler func(*myproto.Message)

type connState struct {
	conn      *mynet.Connection
	conStatus bool
	regStatus bool
}

func (b *builder) addHandler() {

	b.handler[myproto.KCmdHeartBeatReq] = b.handleHeartBeatReq
	b.handler[myproto.KCmdRegisterRsp] = b.handleRegisterRsp
	b.handler[myproto.KCmdTaskBuildReq] = b.handleTaskBuildReq
	b.handler[myproto.KCmdTaskStateReq] = b.handleBuildStateReq

}

func (b *builder) handle(p *myproto.Message) error {
	b.log.Info("handle:\n%s\n", myproto.Debug(p))
	handler, ok := b.handler[p.Head.Command]
	if !ok {
		return fmt.Errorf("command 0x%x handler not found", p.Head.Command)
	}
	go handler(p)

	return nil
}

func (b *builder) eventTask() {
	n := b.net
	for {
		evt, err := n.PollEvent(1000 * 60)
		if err != nil {
			b.log.Error("poll event error!")
			_, ok := <-b.cmdChan
			if ok {
				b.cmdChan <- "stop"
			}
			return
		}
		conn := evt.Conn
		switch {
		case evt.EventType == mynet.EventConnectionError:
			{
				b.log.Info("event error: local = %s, remote = %s\n",
					conn.LocalAddress(), conn.RemoteAddress())
				b.cmdChan <- "reconnect"
			}
		case evt.EventType == mynet.EventConnectionClosed:
			{
				b.log.Info("event close: local = %s, remote = %s\n",
					conn.LocalAddress(), conn.RemoteAddress())
				b.cmdChan <- "reconnect"
			}
		case evt.EventType == mynet.EventNewConnectionData:
			{
				data := evt.Data.(*myproto.Message)
				if err := b.handle(data); err != nil {
					b.log.Error("handle message failed. err = %s\n", err)
					continue
				}
			}
		case evt.EventType == mynet.EventTimeout:
			{
				b.log.Info("poll timeout now = %d\n", time.Now().Unix())
			}
		}
	}
}

func (b *builder) parseArgs() {
	host := flag.String("host", "", "the server adress, format: ip:port")
	debug := flag.Bool("debug", false, "run in debug model")
	help := flag.Bool("help", false, "show the help message")

	flag.Parse()
	if *help {
		flag.PrintDefaults()
		os.Exit(0)
	}
	if len(*host) == 0 {
		fmt.Printf("connect host is empty. Usage:\n")
		flag.PrintDefaults()
		os.Exit(-1)
	}

	pwd, _ := filepath.Abs(filepath.Dir(os.Args[0]))
	logDir := pwd + string(filepath.Separator) + "logs" + string(filepath.Separator)
	buildDir := pwd + string(filepath.Separator) + "build" + string(filepath.Separator)
	dbDir := pwd + string(filepath.Separator) + "db" + string(filepath.Separator)

	if _, err := os.Stat(logDir); err != nil {
		if !os.IsNotExist(err) {
			fmt.Printf("log dir %s don't have permission\n", logDir)
			os.Exit(-1)
		}
		os.Mkdir(logDir, 0777)
	}

	if _, err := os.Stat(buildDir); err != nil {
		if !os.IsNotExist(err) {
			fmt.Printf("build dir %s don't have permission\n", buildDir)
			os.Exit(-1)
		}
		os.Mkdir(buildDir, 0777)
	}
	if _, err := os.Stat(dbDir); err != nil {
		if !os.IsNotExist(err) {
			fmt.Printf("build dir %s don't have permission\n", buildDir)
			os.Exit(-1)
		}
		os.Mkdir(dbDir, 0777)
	}

	b.host = *host
	b.debug = *debug
	b.logDir = logDir
	b.buildDir = buildDir
	b.dbFile = dbDir + "buildtask.db"
}

func (b *builder) initLog(prefix string) (err error) {
	b.log, err = mylog.NewLogging()
	if err != nil {
		return err
	}

	level := mylog.LevelAll
	if !b.debug {
		level = mylog.LevelInformational
	}
	_, err = mylog.SetupLog("file",
		fmt.Sprintf(`{"prefix":"%s", "filedir":"%s", "level":%d, "switchsize":%d, "switchtime":%d}`,
			prefix, b.logDir, level, -1, 0))
	if err != nil {
		return err
	}
	_, err = mylog.SetupLog("console",
		fmt.Sprintf(`{"level":%d}`, level))
	if err != nil {
		return err
	}
	b.log.Start()

	b.log.Info("log started.\n")

	return nil
}

func (b *builder) init() error {
	b.parseArgs()

	var err error

	b.db, err = NewBuildDb(b.dbFile)
	if err != nil {
		return err
	}

	err = b.initLog("builder")
	if err != nil {
		return fmt.Errorf("init log failed, err = %s", err)
	}

	b.net = mynet.NewSimpleNet(b.log)
	b.proto = myproto.NewProto()
	b.hostInfo = util.GetHostInfo()
	b.handler = make(map[uint32]runnerHandler)
	b.sess = util.NewSession()

	b.cmdChan = make(chan string, 1024)
	b.exitChan = make(chan struct{})

	b.buildQueue = make(map[string]string)
	b.buildQueueLock = &sync.Mutex{}

	b.addHandler()

	go b.sigTask()

	go b.eventTask()
	go b.sessionCheck()

	b.log.Info("start command task go routine\n")

	go b.cmdTask()

	return nil
}

func (b *builder) cmdTask() {
	for {
		cmd, ok := <-b.cmdChan
		if !ok {
			b.log.Error("command chan closed\n")
			return
		}
		switch cmd {
		case "connect":
			{
				conn, err := b.net.Connect(b.host, b.proto)
				if err != nil {
					b.cmdChan <- "stop"
					continue
				}
				b.conn = &connState{
					conn:      conn,
					conStatus: true,
					regStatus: false,
				}
				b.log.Info("connect to %s success!\n", b.host)
				b.cmdChan <- "register"
			}
		case "reconnect":
			{
				b.conn = &connState{}
				conn, err := b.net.Connect(b.host, b.proto)
				if err != nil {
					b.log.Info("reconnect %s failed, err = %s, try after 15 second\n",
						b.host, err)
					go b.timerOnce(15, "reconnect", b.cmdChan, "reconnect")
					continue
				}
				b.conn = &connState{
					conn:      conn,
					conStatus: true,
					regStatus: false,
				}
				b.log.Info("reconnect to %s success!\n", b.host)
				b.cmdChan <- "register"
			}
		case "register":
			{
				err := b.register()
				if err != nil {
					b.log.Error("construct register failed\n")
					b.cmdChan <- "stop"
					continue
				}
			}
		case "stop":
			{
				b.log.Info("command chan receive stop command\n")
				b.stop()
				return
			}
		}
	}
}

func (b *builder) sigTask() {
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGQUIT, syscall.SIGHUP, syscall.SIGTERM, syscall.SIGUSR2, syscall.SIGUSR1)
END:
	for {
		select {
		case sig := <-ch:
			switch sig {
			case syscall.SIGINT:
				b.log.Info("catch SIGINT, stop the builder.\n")
				b.cmdChan <- "stop"
				break END
			case syscall.SIGQUIT:
				b.log.Info("catch SIGQUIT, stop the builder.\n")
				b.cmdChan <- "stop"
				break END
			case syscall.SIGHUP:
				b.log.Info("catch SIGHUP, stop the builder.\n")
				b.cmdChan <- "stop"
				break END
			case syscall.SIGTERM:
				b.log.Info("catch SIGTERM, stop the builder.\n")
				b.cmdChan <- "stop"
				break END
			case syscall.SIGUSR1:
				b.log.Info("catch SIGUSR1.\n")
			case syscall.SIGUSR2:
				b.log.Info("catch SIGUSR2.\n")
				b.log.Sync()

			}
		}

	}
}

func (b *builder) register() error {
	m := &myproto.Message{}
	m.Head.Command = myproto.KCmdRegisterReq
	// beartbeat
	p, err := myproto.GetMessage(m.Head.Command)
	if err != nil {
		b.log.Error("gen register failed, err = %s\n", err)
		return err
	}
	reg := p.(*myproto.RegisterReq)
	reg.SID = proto.String(util.GetSID(16))
	reg.Type = proto.String("builder")
	reg.OS = proto.String(b.hostInfo.OS)
	reg.Host = proto.String(b.hostInfo.Host)
	reg.Arch = proto.String(b.hostInfo.ARCH)
	reg.Adress = b.hostInfo.Adress

	m.Body = reg

	b.log.Info("REQ：\n%s\n", myproto.Debug(m))
	err = b.net.SendData(b.conn.conn, m)
	if err != nil {
		b.log.Error("send register error, err = %s\n", err)
		return err
	}
	b.sess.Add(reg.GetSID(), m)
	return nil
}

func (b *builder) sessionCheck() {
	for {
		t := time.NewTimer((time.Duration)((int64)(time.Second) * 30))
		<-t.C
		b.log.Debug("session check\n")
		to := b.sess.GetTimeout(60)
		for k, v := range to {
			b.log.Warning("timeout sid = %s\n", k)
			go b.handleTimeout(k, v)
		}
	}

}

func (b *builder) timerOnce(to int64, typ string, ch interface{}, cmd interface{}) {
	t := time.NewTimer((time.Duration)((int64)(time.Second) * to))
	<-t.C

	switch {
	case typ == "reconnect":
		cmdChan := (ch).(chan string)
		msg := (cmd).(string)

		cmdChan <- msg
	}
}

// Start the client
func (b *builder) start() {
	b.cmdChan <- "connect"
}

// Stop the client
func (b *builder) stop() {
	close(b.cmdChan)
	mynet.SimpleNetDestroy(b.net)
	BuildDbDestroy(b.db)
	//		t := time.NewTimer((time.Duration)((int64)(time.Second) * 5))
	//		<-t.C

	b.log.Stop()

	b.exitChan <- struct{}{}
}

// Wait wait the client to stop
func (b *builder) wait() {
	<-b.exitChan
}

func main() {
	b := &builder{}

	err := b.init()
	if err != nil {
		fmt.Printf("%s\n", err)
		os.Exit(-1)
	}

	b.start()
	b.wait()

}
