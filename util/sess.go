package util

import (
	"fmt"
	"sync"
	"time"
)

type sessionData struct {
	v      interface{}
	uptime time.Time
}

type Session struct {
	sess map[string]*sessionData
	lock sync.Locker
}

func NewSession() *Session {
	s := &Session{
		sess: make(map[string]*sessionData),
		lock: &sync.Mutex{},
	}
	return s
}

func (s *Session) Add(sid string, v interface{}) error {
	s.lock.Lock()
	defer s.lock.Unlock()

	if _, ok := s.sess[sid]; ok {
		return fmt.Errorf("sid %s exists", sid)
	}
	data := &sessionData{
		v:      v,
		uptime: time.Now(),
	}
	s.sess[sid] = data
	return nil
}
func (s *Session) Del(sid string) error {
	s.lock.Lock()
	defer s.lock.Unlock()

	if _, ok := s.sess[sid]; !ok {
		return fmt.Errorf("sid %s not exists", sid)
	}

	delete(s.sess, sid)
	return nil
}
func (s *Session) Get(sid string) interface{} {
	s.lock.Lock()
	defer s.lock.Unlock()

	v, ok := s.sess[sid]
	if !ok {
		return nil
	}
	return v.v
}

func (s *Session) GetTimeout(timeout int) map[string]interface{} {
	s.lock.Lock()
	defer s.lock.Unlock()
	m := make(map[string]interface{})
	n := time.Now().Unix()
	for k, v := range s.sess {
		diff := n - v.uptime.Unix()
		if diff >= (int64)(timeout) {
			m[k] = v.v
		}
	}
	for k := range m {
		delete(s.sess, k)
	}
	if len(m) > 0 {
		return m
	}
	return nil
}
