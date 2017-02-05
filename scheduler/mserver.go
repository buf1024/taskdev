package main

import (
	"flag"
	"fmt"
	"net"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"

	"taskdev/logging"
)

type scheMsg struct {
	command int64
	message []byte
}

type clientConn struct {
	conn      *net.TCPConn
	conStatus bool
	regStatus bool

	recvChan chan *scheMsg
	sendChan chan *scheMsg

	uptime int64
}

// Server represents the the client connect to scheduler
type Server struct {
	Host  string
	Debug bool
	Dir   string

	listen *net.TCPListener

	Log *logging.Log

	sessChan  chan struct{}
	timerChan chan struct{}

	cmdChan  chan string
	exitChan chan struct{}

	clientQueue []*clientConn
}

const (
	statusInited = iota
	statusStarted
	statusReady
	statusStoped
)

func (m *Server) scheHandleMsg() {
END:
	for {
		select {
		// case msg := <-m.exch.recvChan:
		// 	{
		// 	}
		// case msg := <-m.exch.sendChan:
		// 	{

		// 	}
		default:
			// chan close
			break END

		}
	}
}

func (m *Server) scheRegister() error {
	// reg packet
	// m.Log.Info("reg to exch\n")
	// msg, err := proto.Message(proto.CMD_SVR_REG_REQ)
	// if err != nil {
	// 	m.Log.Critical("create message failed, ERR=%s\n", err.Error())
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
	// 	m.Log.Critical("serialize message failed, ERR=%s\n", err.Error())
	// 	return err
	// }

	// return err
	return nil
}

func (m *Server) timerOnce(to int64, typ string, ch interface{}, cmd interface{}) {
	t := time.NewTimer((time.Duration)((int64)(time.Second) * to))
	<-t.C

	switch {
	case typ == "schereconnect":
		exchChan := (ch).(chan string)
		scheMsg := (cmd).(string)

		exchChan <- scheMsg
	}
}

func (m *Server) handleConnRead(client *clientConn) {
	for {
		// read package
		// client.conn.Read
	}
}

func (m *Server) handleConnWrite(client *clientConn) {
	for {
		select {
		case msg := <-client.sendChan:
			// send msg
		}
	}
}

func (m *Server) handleConn(client *clientConn) {
	go m.handleConnRead(client)
	go m.handleConnWrite(client)
}

func (m *Server) listening() {
	for {
		conn, err := m.listen.AcceptTCP()
		if err != nil {
			m.Log.Error("accept failed, err = %s\n", err)
			continue
		}
		addr := conn.RemoteAddr()
		m.Log.Info("accept new client, addr=%s\n", addr)
		client := &clientConn{}

		client.conn = conn
		client.conStatus = true
		client.regStatus = false
		client.uptime = time.Now().UnixNano() / (1000 * 1000 * 1000)

		m.clientQueue = append(m.clientQueue, client)

		go m.handleConn(client)
	}
}

func (m *Server) startListen() error {
	listen, err := net.Listen("tcp", m.Host)
	if err != nil {
		return err
	}
	m.listen = (listen).(*net.TCPListener)

	go m.listening()

	return nil
}

// task represents exchange connect go routine
func (m *Server) cmdTask() {
	for {
		cmd := <-m.cmdChan
		switch cmd {
		case "listen":
			{
				m.Log.Info("start to listen, addr = %s\n", m.Host)
				err := m.startListen()
				if err != nil {
					m.Log.Critical("listen failed, err = %s\n", err)
					m.Stop()
				} else {
					m.Log.Info("listen success, waiting for connection\n")
				}
			}
		case "stop":
			{
				fmt.Printf("stop")
			}
		}
	}
}

func (m *Server) sigTask() {
	ch := make(chan os.Signal)
	signal.Notify(ch, syscall.SIGHUP, syscall.SIGTERM, syscall.SIGUSR2, syscall.SIGUSR1)
END:
	for {
		select {
		case sig := <-ch:
			switch sig {
			case syscall.SIGHUP:
				m.Log.Info("catch SIGHUP, stop the server.")
				m.Stop()
				break END
			case syscall.SIGTERM:
				m.Log.Info("catch SIGTERM, stop the server.")
				m.Stop()
				break END
			case syscall.SIGUSR1:
				m.Log.Info("catch SIGUSR1.")
			case syscall.SIGUSR2:
				m.Log.Info("catch SIGUSR2.")
				m.Log.Sync()

			}
		}

	}

}

// initLog Setup Logger
func (m *Server) initLog(prefix string) (err error) {
	m.Log, err = logging.NewLogging()
	if err != nil {
		return err
	}
	level := logging.LevelAll
	if !m.Debug {
		level = logging.LevelInformational
	}
	_, err = logging.SetupLog("file",
		fmt.Sprintf(`{"prefix":"%s", "filedir":"%s", "level":%d, "switchsize":%d, "switchtime":%d}`,
			prefix, m.Dir, level, -1, 0))
	if err != nil {
		return err
	}
	_, err = logging.SetupLog("console",
		fmt.Sprintf(`{"level":%d}`, level))
	if err != nil {
		return err
	}
	m.Log.Start()

	m.Log.Info("log started.\n")

	return nil
}

// parseArgs parse the commandline arguments
func (m *Server) parseArgs() {
	host := flag.String("l", "127.0.0.1:4421", "the scheduler server listen address")
	debug := flag.Bool("d", false, "run in debug model")
	help := flag.Bool("h", false, "show the help message")

	flag.Parse()
	if *help {
		flag.PrintDefaults()
		os.Exit(0)
	}
	if len(*host) == 0 {
		fmt.Printf("connect host is empty.\n")
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

	m.Host = *host
	m.Debug = (*debug)
	m.Dir = dir

}

// InitServer init the client
func (m *Server) InitServer() {
	m.parseArgs()

	err := m.initLog("scheduler")
	if err != nil {
		fmt.Printf("init log failed. err=%s\n", err)
		os.Exit(-1)
	}

	m.cmdChan = make(chan string, 1024)
	m.exitChan = make(chan struct{})

	m.clientQueue = make([]*clientConn, 1024)

	go m.sigTask()

	m.Log.Info("start task go routine\n")
	go m.cmdTask()
}

// ServerForever the client
func (m *Server) ServerForever() {

	m.cmdChan <- "listen"

	m.Wait()
}

// Stop the client
func (m *Server) Stop() {
	m.Log.Stop()

	m.exitChan <- struct{}{}
}

// Wait wait the client to stop
func (m *Server) Wait() {
	<-m.exitChan
}
