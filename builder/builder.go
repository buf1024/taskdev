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

	mylog "github.com/buf1024/golib/logging"
	mynet "github.com/buf1024/golib/net"
)

type builder struct {
	net      *mynet.SimpleNet
	proto    *myproto.TaskProto
	log      *mylog.Log
	handler  map[uint32]runnerHandler
	conn     *connState
	hostInfo *util.HostInfo

	sessChan  chan struct{}
	timerChan chan struct{}
	cmdChan   chan string
	exitChan  chan struct{}

	host   string
	logDir string
	debug  bool
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
}

func (b *builder) handle(p *myproto.Message) error {
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
			b.cmdChan <- "stop"
			return
		}
		conn := evt.Conn
		switch {
		case evt.EventType == mynet.EventConnectionError:
			{
				fmt.Printf("event error: local = %s, remote = %s\n",
					conn.LocalAddress(), conn.RemoteAddress())
				b.cmdChan <- "reconnect"
				return
			}
		case evt.EventType == mynet.EventConnectionClosed:
			{
				fmt.Printf("event close: local = %s, remote = %s\n",
					conn.LocalAddress(), conn.RemoteAddress())
				b.cmdChan <- "reconnect"
				return
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

	dir := filepath.Dir(os.Args[0])
	dir = dir + string(filepath.Separator) + "logs" + string(filepath.Separator)

	if _, err := os.Stat(dir); err != nil {
		if !os.IsNotExist(err) {
			fmt.Printf("log dir %s don't have permission\n", dir)
			os.Exit(-1)
		}
		os.Mkdir(dir, 0777)
	}

	b.host = *host
	b.debug = *debug
	b.logDir = dir

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

	err := b.initLog("builder")
	if err != nil {
		return fmt.Errorf("init log failed, err = %s", err)
	}

	b.net = mynet.NewSimpleNet()
	b.proto = myproto.NewProto()
	b.hostInfo = util.GetHostInfo()
	b.handler = make(map[uint32]runnerHandler)

	b.cmdChan = make(chan string, 1024)
	b.exitChan = make(chan struct{})

	b.addHandler()

	go b.sigTask()

	go b.eventTask()

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
					b.log.Info("reconnected %s failed, try after 30 second\n",
						b.host)
					go b.timerOnce(30, "reconnect", b.cmdChan, "reconnect")
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
	ch := make(chan os.Signal)
	signal.Notify(ch, syscall.SIGHUP, syscall.SIGTERM, syscall.SIGUSR2, syscall.SIGUSR1)
END:
	for {
		select {
		case sig := <-ch:
			switch sig {
			case syscall.SIGHUP:
				b.log.Info("catch SIGHUP, stop the server.")
				b.cmdChan <- "stop"
				break END
			case syscall.SIGTERM:
				b.log.Info("catch SIGTERM, stop the server.")
				b.cmdChan <- "stop"
				break END
			case syscall.SIGUSR1:
				b.log.Info("catch SIGUSR1.")
			case syscall.SIGUSR2:
				b.log.Info("catch SIGUSR2.")
				b.log.Sync()

			}
		}

	}

}

func (b *builder) register() error {
	// reg packet
	// b.log.Info("reg to exch\n")
	// msg, err := proto.Message(proto.CMD_SVR_REG_REQ)
	// if err != nil {
	// 	b.log.Critical("create message failed, ERR=%s\n", err.Error())
	// 	return err
	// }

	// req := (*bankmsg.SvrRegReq)(msg)
	// *req.SID = "123"
	// *req.SvrType = bankid
	// *req.SvrId = bankid

	// reqMsg := &scheMsg{}
	// reqMsg.conn = m.exchCtx.conn
	// reqMsg.command = proto.CMD_SVR_REG_REQ
	// reqMsg.message, err = proto.Serialize(req)
	// if err != nil {
	// 	b.log.Critical("serialize message failed, ERR=%s\n", err.Error())
	// 	return err
	// }

	// return err
	return nil
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
