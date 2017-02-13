package util

import (
	"fmt"
	"math/rand"
	"net"
	"os"
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
		s = fmt.Sprintf("%s%x", s, v)
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

func init() {
	rand.NewSource(time.Now().UnixNano())
}
