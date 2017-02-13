package main

import myproto "taskdev/proto"

func (b *builder) handleHeartBeatReq(m *myproto.Message) {
	b.log.Info("handleHeartBeatReq\n")

}

func (b *builder) handleRegisterRsp(m *myproto.Message) {

	b.log.Info("handleRegisterRsp\n")

}
