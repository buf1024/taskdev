package main

import (
	myproto "taskdev/proto"
	"taskdev/util"

	"github.com/golang/protobuf/proto"
)

func (s *scheServer) handleTimeout(sid string, v interface{}) {

}

func (s *scheServer) handleHeartBeatReq(conn *connState, reqMsg *myproto.Message) {
	reqPb := reqMsg.Body.(*myproto.HeartBeatReq)

	rspMsg := &myproto.Message{}
	rspMsg.Head.Command = reqMsg.Head.Command + 1
	// beartbeat
	rspPb, err := myproto.GetMessage(rspMsg.Head.Command)
	if err != nil {
		s.log.Error("gen register failed, err = %s\n", err)
		return
	}
	pb := rspPb.(*myproto.HeartBeatRsp)
	pb.SID = proto.String(reqPb.GetSID())

	rspMsg.Body = rspPb

	s.log.Info("RSP：\n%s\n", myproto.Debug(rspMsg))
	err = s.net.SendData(conn.conn, rspMsg)
	if err != nil {
		s.log.Error("send hearbeat, err = %s\n", err)
	}

}

func (s *scheServer) handleRegisterReq(conn *connState, reqMsg *myproto.Message) {
	conn.regStatus = true
	reqPb := reqMsg.Body.(*myproto.RegisterReq)

	rspMsg := &myproto.Message{}
	rspMsg.Head.Command = reqMsg.Head.Command + 1

	rspPb, err := myproto.GetMessage(rspMsg.Head.Command)
	if err != nil {
		s.log.Error("gen register failed, err = %s\n", err)
		return
	}
	pb := rspPb.(*myproto.RegisterRsp)
	pb.SID = proto.String(reqPb.GetSID())
	pb.RetCode = proto.Int32(util.KESuccess)

	rspMsg.Body = rspPb
	s.log.Info("RSP：\n%s\n", myproto.Debug(rspMsg))
	err = s.net.SendData(conn.conn, rspMsg)
	if err != nil {
		s.log.Error("send server register, err = %s\n", err)
	}
}
