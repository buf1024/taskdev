package util

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"fmt"
)

// MsgHead the protocol header
type MsgHead struct {
	Len     uint32 // 4
	Cmd     uint32 // 4
	ReqSid  string // 8
	RspSid  string // 8
	RetCode uint32 // 4
}

// MsgBody represents the protocol body
type MsgBody interface{}

// Message protocol definition
type Message struct {
	Head MsgHead
	Body MsgBody
}

// MsgHeadLen the head length
const MsgHeadLen uint32 = 28

func newMsgBody(cmd uint32) MsgBody {
	body := MsgBody(nil)
	switch cmd {
	case CMD_HEARTBEAT_REQ:
		body = &MsgHeartBeatReq{}
	case CMD_HEARBEAT_RSP:
		body = &MsgHeartBeatRsp{}
	case CMD_REGISTER_REQ:
		body = &MsgRegReq{}
	case CMD_REGISTER_RSP:
		body = &MsgRegRsp{}
	default:
		return nil
	}
	return body
}

func NewMessage(cmd uint32) *Message {
	m := &Message{}
	m.Head.Cmd = cmd
	if cmd%2 == 0 {
		m.Head.RspSid = SID(8)
	} else {
		m.Head.ReqSid = SID(8)
	}
	m.Body = newMsgBody(cmd)

	return m
}

func ParseHead(head []byte) (*MsgHead, error) {
	if uint32(len(head)) < MsgHeadLen {
		return nil, fmt.Errorf("head len is less than %d", MsgHeadLen)
	}

	h := &MsgHead{}

	buf := bytes.NewReader(head[0:4])
	err := binary.Read(buf, binary.BigEndian, &h.Len)
	if err != nil {
		return nil, err
	}
	buf = bytes.NewReader(head[4:8])
	err = binary.Read(buf, binary.BigEndian, &h.Cmd)
	if err != nil {
		return nil, err
	}
	buf = bytes.NewReader(head[8:16])
	err = binary.Read(buf, binary.BigEndian, &h.ReqSid)
	if err != nil {
		return nil, err
	}
	buf = bytes.NewReader(head[16:24])
	err = binary.Read(buf, binary.BigEndian, &h.RspSid)
	if err != nil {
		return nil, err
	}
	buf = bytes.NewReader(head[24:28])
	err = binary.Read(buf, binary.BigEndian, &h.RetCode)
	if err != nil {
		return nil, err
	}

	return h, nil
}
func ParseBody(cmd uint32, body []byte) (MsgBody, error) {
	if len(body) == 0 {
		return nil, nil
	}
	b := make([]byte, len(body))

	buf := bytes.NewReader(body)
	err := binary.Read(buf, binary.BigEndian, b)
	if err != nil {
		return nil, err
	}

	mb := newMsgBody(cmd)
	if mb == nil {
		return nil, nil
	}

	if err = json.Unmarshal(b, mb); err != nil {
		return nil, err
	}

	return mb, nil
}

// Serialize serialize proto to send
func Serialize(m *Message) ([]byte, error) {
	buf := new(bytes.Buffer)

	err := binary.Write(buf, binary.BigEndian, m.Body)
	if err != nil {
		return nil, err
	}

	bytebody := buf.Bytes()
	m.Head.Len = uint32(len(bytebody)) + MsgHeadLen

	buf = new(bytes.Buffer)

	err = binary.Write(buf, binary.BigEndian, m.Head.Len)
	if err != nil {
		return nil, err
	}
	err = binary.Write(buf, binary.BigEndian, m.Head.Cmd)
	if err != nil {
		return nil, err
	}
	err = binary.Write(buf, binary.BigEndian, m.Head.ReqSid)
	if err != nil {
		return nil, err
	}
	err = binary.Write(buf, binary.BigEndian, m.Head.RspSid)
	if err != nil {
		return nil, err
	}
	err = binary.Write(buf, binary.BigEndian, m.Head.RetCode)
	if err != nil {
		return nil, err
	}
	bytehead := buf.Bytes()

	msgbyte := append(bytehead, bytebody...)

	return msgbyte, nil

}

//====================//
const (
	CMD_HEARTBEAT_REQ uint32 = 0x00000001
	CMD_HEARBEAT_RSP  uint32 = 0x00000002

	CMD_REGISTER_REQ uint32 = 0x00000003
	CMD_REGISTER_RSP uint32 = 0x00000004
)

type MsgHeartBeatReq struct {
}
type MsgHeartBeatRsp struct {
}

type MsgRegReq struct {
	Type string   `json:"type"`
	Host string   `json:"host"`
	OS   string   `json:"os"`
	ARCH string   `json:"arch"`
	IP   []string `json:"ip"`
}
type MsgRegRsp struct {
}
