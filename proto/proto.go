package proto

import (
	"bytes"
	"encoding/binary"
	"fmt"

	mynet "github.com/buf1024/golib/net"
	pb "github.com/golang/protobuf/proto"
)

// MessageHead 协议头
type MessageHead struct {
	Length  uint32 // 4
	Command uint32 // 4
}

// MessageBody 协议体
type MessageBody pb.Message

// Message 协议定义
type Message struct {
	Head MessageHead
	Body MessageBody

	/*
		ReqSid  string // 16
		RspSid  string // 16
		RetCode uint32 // 4
	*/
}

type TaskProto struct {
}

const (
	// KMessageHeadLen 协议头长度
	KMessageHeadLen uint32 = 8
)

const (
	KCmdHeartBeatReq uint32 = 0x00010001
	KCmdHeartBeatRsp uint32 = 0x00010002
	KCmdRegisterReq  uint32 = 0x00010003
	KCmdRegisterRsp  uint32 = 0x00010004
	KCmdTaskBuildReq uint32 = 0x00010005
	KCmdTaskBuildRsp uint32 = 0x00010006
	KCmdTaskStateReq uint32 = 0x00010007
	KCmdTaskStateRsp uint32 = 0x00010008
)

var message map[uint32]pb.Message

func (p *TaskProto) FilterAccept(conn *mynet.Connection) bool {
	return true
}
func (p *TaskProto) HeadLen() uint32 {
	return KMessageHeadLen
}
func (p *TaskProto) BodyLen(head []byte) (interface{}, uint32, error) {
	if (uint32)(len(head)) != KMessageHeadLen {
		return nil, 0, fmt.Errorf("head size not right")
	}
	h := &MessageHead{}

	buf := head[:4]
	reader := bytes.NewReader(buf)
	err := binary.Read(reader, binary.BigEndian, &h.Length)
	if err != nil {
		return nil, 0, err
	}

	buf = head[4:8]
	reader = bytes.NewReader(buf)
	err = binary.Read(reader, binary.BigEndian, &h.Command)
	if err != nil {
		return nil, 0, err
	}

	return h, h.Length, nil
}
func (p *TaskProto) Parse(head interface{}, body []byte) (interface{}, error) {
	m := &Message{
		Head: *(head.(*MessageHead)),
	}
	pbm, err := GetMessage(m.Head.Command)
	if err != nil {
		return nil, err
	}
	err = pb.Unmarshal(body, pbm)
	if err != nil {
		return nil, err
	}

	m.Body = pbm

	return m, nil
}
func (p *TaskProto) Serialize(data interface{}) ([]byte, error) {
	m := data.(*Message)

	body, err := pb.Marshal(m.Body)
	if err != nil {
		return nil, err
	}
	m.Head.Length = (uint32)(len(body))

	buf := new(bytes.Buffer)
	if err = binary.Write(buf, binary.BigEndian, m.Head.Length); err != nil {
		return nil, err
	}
	if err = binary.Write(buf, binary.BigEndian, m.Head.Command); err != nil {
		return nil, err
	}

	head := buf.Bytes()

	return append(head, body...), nil

}

func NewProto() *TaskProto {
	p := &TaskProto{}
	return p
}

func Debug(msg *Message) string {
	return fmt.Sprintf("command:0x%x\nlength :%d\n==========\n%s",
		msg.Head.Command, msg.Head.Length, pb.MarshalTextString(msg.Body))
}

func GetMessage(command uint32) (pb.Message, error) {
	if m, ok := message[command]; ok {
		pb.Clone(m)
		return m, nil
	}
	return nil, fmt.Errorf("command %d not found", command)
}

func init() {
	message = make(map[uint32]pb.Message)

	message[KCmdHeartBeatReq] = &HeartBeatReq{} //0x00010001 // 心跳请求
	message[KCmdHeartBeatRsp] = &HeartBeatRsp{} //0x00010002
	message[KCmdRegisterReq] = &RegisterReq{}   //0x00010003
	message[KCmdRegisterRsp] = &RegisterRsp{}   //0x00010004
	message[KCmdTaskBuildReq] = &TaskBuildReq{} //0x00010005
	message[KCmdTaskBuildRsp] = &TaskBuildRsp{} //0x00010006
	message[KCmdTaskStateReq] = &TaskStateReq{} //0x00010007
	message[KCmdTaskStateRsp] = &TaskStateRsp{} //0x00010008
}
