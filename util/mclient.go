package util

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

// TaskHandler represents the businiess handler
type TaskHandler func(*Message) (*Message, error)

// ClientInfo the client Info
type ClientInfo struct {
	Name    string
	Handler map[int32]TaskHandler

	Type string
	Host string
	OS   string
	ARCH string
	IP   []string
}

// ClientHook the real client interface
type ClientHook interface {
	HookInfo() *ClientInfo
}

type scheMsg struct {
	command int64
	message []byte
}

type scheConn struct {
	conn      *net.TCPConn
	conStatus bool
	regStatus bool

	recvChan chan *scheMsg
	sendChan chan *scheMsg
}

// Client represents the the client connect to scheduler
type Client struct {
	Host  string
	Debug bool
	Dir   string

	Log  *logging.Log
	sche scheConn

	sessChan  chan struct{}
	timerChan chan struct{}

	cmdChan  chan string
	exitChan chan struct{}

	clientInfo *ClientInfo

	status int32
}

const (
	statusInited = iota
	statusStarted
	statusReady
	statusStoped
)

func (m *Client) scheHandleMsg() {
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

func (m *Client) scheConnect() error {
	m.Log.Info("connect to scheduler, addr = %s\n", m.Host)
	conn, err := net.DialTimeout("tcp", m.Host, time.Second*60)
	if err != nil {
		m.Log.Error("connect to scheduler failed, addr = %s, err=%s\n",
			m.Host, err)
		return err
	}

	(conn).(*net.TCPConn).SetNoDelay(true)
	(conn).(*net.TCPConn).SetKeepAlive(true)

	m.sche.conn = conn.(*net.TCPConn)
	m.sche.conStatus = true
	m.sche.regStatus = false

	if m.sche.recvChan != nil {
		close(m.sche.recvChan)
	}
	if m.sche.sendChan != nil {
		close(m.sche.sendChan)
	}

	m.sche.recvChan = make(chan *scheMsg, 1024)
	m.sche.sendChan = make(chan *scheMsg, 1024)

	go m.scheHandleMsg()

	m.Log.Info("scheduler connected\n")

	return nil
}
func (m *Client) scheRegister() error {
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

func (m *Client) timerOnce(to int64, typ string, ch interface{}, cmd interface{}) {
	t := time.NewTimer((time.Duration)((int64)(time.Second) * to))
	<-t.C

	switch {
	case typ == "schereconnect":
		exchChan := (ch).(chan string)
		scheMsg := (cmd).(string)

		exchChan <- scheMsg
	}
}

// task represents exchange connect go routine
func (m *Client) cmdTask() {
	for {
		cmd := <-m.cmdChan
		switch cmd {
		case "connect":
			{
				err := m.scheConnect()
				if err != nil {
					m.Stop()
					return
				}
				m.cmdChan <- "register"
			}
		case "reconnect":
			{
				err := m.scheConnect()
				if err != nil {
					m.Log.Info("reconnect to addr=%s\n after 30 second",
						m.Host)
					go m.timerOnce(30, "schereconnect", m.cmdChan, "reconnect")
					continue
				}
				m.cmdChan <- "register"
			}
		case "register":
			{
				err := m.scheRegister()
				if err != nil {
					// err
				} else {
					// post
				}
			}
		case "stop":
			{
				fmt.Printf("stop")
			}
		}
	}
}

func (m *Client) sigTask() {
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

// InitLog Setup Logger
func (m *Client) InitLog(prefix string) (err error) {
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

// ParseArgs parse the commandline arguments
func (m *Client) ParseArgs() {
	host := flag.String("p", "", "the scheduler server host, format: ip:port")
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

// InitClient init the client
func (m *Client) InitClient() {
	m.cmdChan = make(chan string, 1024)
	m.exitChan = make(chan struct{})

	go m.sigTask()

	m.Log.Info("start task go routine\n")
	go m.cmdTask()
}

// Start the client
func (m *Client) Start() {

	m.cmdChan <- "connect"
}

// Stop the client
func (m *Client) Stop() {
	m.Log.Stop()

	m.exitChan <- struct{}{}
}

// Wait wait the client to stop
func (m *Client) Wait() {
	<-m.exitChan
}

// SetHook adds the command handler
func (m *Client) SetHook(hook ClientHook) {
	if hook != nil {
		m.clientInfo = hook.HookInfo()
	}
}
