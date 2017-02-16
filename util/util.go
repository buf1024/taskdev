package util

import (
	"fmt"
	"math/rand"
	"net"
	"os"
	"os/exec"
	"runtime"
	"time"
)

type HostInfo struct {
	Host   string
	OS     string
	ARCH   string
	Adress []string
}

func GetSID(len int32) string {
	sid := make([]byte, len/2)
	rand.Read(sid)
	s := ""
	for _, v := range sid {
		s = fmt.Sprintf("%s%02x", s, v)
	}
	return s
}

func GetHostInfo() *HostInfo {
	info := &HostInfo{}

	info.Host, _ = os.Hostname()
	info.OS = runtime.GOOS
	info.ARCH = runtime.GOARCH

	addrs, err := net.InterfaceAddrs()
	if err != nil {
		info.Adress = nil
		return info
	}
	for _, address := range addrs {
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ip := ipnet.IP.To4(); ip != nil {
				info.Adress = append(info.Adress, ip.String())
			}

		}
	}

	return info
}

func Run(exe string, args ...string) (errout string, out string, err error) {
	cmd := exec.Command(exe, args...)
	stderr, err := cmd.StderrPipe()
	if err != nil {
		return "", "", err

	}
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return "", "", err
	}
	err = cmd.Start()
	if err != nil {
		return "", "", err
	}
	for {
		buf := make([]byte, 1024)
		n, err := stderr.Read(buf)
		if n > 0 {
			errout = fmt.Sprintf("%s%s", errout, (string)(buf))
		}
		if err != nil {
			break
		}
	}
	for {
		buf := make([]byte, 1024)
		n, err := stdout.Read(buf)
		if n > 0 {
			out = fmt.Sprintf("%s%s", out, (string)(buf))
		}
		if err != nil {
			break
		}
	}
	err = cmd.Wait()
	if err != nil {
		return "", "", err
	}
	return errout, out, nil
}

func init() {
	rand.NewSource(time.Now().UnixNano())
}
