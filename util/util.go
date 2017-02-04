package util

import (
	"math/rand"
	"time"
)

// SID generate len random string
func SID(len int32) string {
	sid := make([]byte, len)
	rand.Read(sid)
	return string(sid)
}

func init() {
	rand.NewSource(time.Now().UnixNano())
}
