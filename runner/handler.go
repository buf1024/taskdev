package main

import (
	myproto "taskdev/proto"
	"taskdev/util"

	"github.com/golang/protobuf/proto"
)

func (b *runner) handleTimeout(sid string, v interface{}) {
	msg := v.(*myproto.Message)
	switch msg.Head.Command {
	case myproto.KCmdRegisterReq:
		b.log.Error("register timeout, try to register again\n")
		if b.conn != nil && b.conn.conStatus && !b.conn.regStatus {
			b.register()
		}
	}
}

func (b *runner) handleHeartBeatReq(reqMsg *myproto.Message) {
	reqPb := reqMsg.Body.(*myproto.HeartBeatReq)

	rspMsg := &myproto.Message{}
	rspMsg.Head.Command = reqMsg.Head.Command + 1
	// beartbeat
	rspPb, err := myproto.GetMessage(rspMsg.Head.Command)
	if err != nil {
		b.log.Error("gen heartbeat failed, err = %s\n", err)
		return
	}
	pb := rspPb.(*myproto.HeartBeatRsp)
	pb.SID = proto.String(reqPb.GetSID())

	rspMsg.Body = rspPb
	b.log.Info("RSP：\n%s\n", myproto.Debug(rspMsg))
	err = b.net.SendData(b.conn.conn, rspMsg)
	if err != nil {
		b.log.Error("send hearbeat, err = %s\n", err)
	}

}

func (b *runner) handleRegisterRsp(reqMsg *myproto.Message) {
	reqPb := reqMsg.Body.(*myproto.RegisterRsp)
	if reqPb.GetRetCode() != util.KESuccess {
		b.log.Error("register failed, try to register again\n")
		b.register()
	} else {
		if b.conn != nil && b.conn.conStatus {
			b.conn.regStatus = true
		}
	}

	b.sess.Del(reqPb.GetSID())
	b.log.Info("success handle register\n")
}
