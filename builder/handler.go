package main

import (
	myproto "taskdev/proto"
	"taskdev/util"

	"fmt"
	"io/ioutil"

	"github.com/golang/protobuf/proto"
)

func (b *builder) handleTimeout(sid string, v interface{}) {
	msg := v.(*myproto.Message)
	switch msg.Head.Command {
	case myproto.KCmdRegisterReq:
		b.log.Error("register timeout, try to register again\n")
		if b.conn != nil && b.conn.conStatus && !b.conn.regStatus {
			b.register()
		}
	}
}

func (b *builder) handleHeartBeatReq(reqMsg *myproto.Message) {
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

func (b *builder) handleRegisterRsp(reqMsg *myproto.Message) {
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

func (b *builder) handleTaskBuildReq(reqMsg *myproto.Message) {
	reqPb := reqMsg.Body.(*myproto.TaskBuildReq)
	sid := reqPb.GetSID()
	if b.sess.Get(sid) != nil {
		return
	}
	b.sess.Add(sid, reqMsg)
	defer b.sess.Del(sid)

	tasks := reqPb.GetTask()
	retcode := util.KESuccess
	for _, task := range tasks {
		//go b.build(task)
		key := fmt.Sprintf("%s-%s", task.GetGroupName(), task.GetName())
		if b.checkBuild(key) {
			retcode = util.KEBuildExist
		}
	}
	if retcode == util.KESuccess {
		for _, task := range tasks {

			go b.build(task)
		}
	}
	rspMsg := &myproto.Message{}
	rspMsg.Head.Command = reqMsg.Head.Command + 1
	pb, err := myproto.GetMessage(rspMsg.Head.Command)
	if err != nil {
		b.log.Error("gen taskbuildrsp failed, err = %s\n", err)
		return
	}
	rspPb := pb.(*myproto.TaskBuildRsp)
	rspPb.SID = proto.String(reqPb.GetSID())
	rspPb.RetCode = proto.Int32(retcode)

	rspMsg.Body = rspPb
	b.log.Info("RSP：\n%s\n", myproto.Debug(rspMsg))
	err = b.net.SendData(b.conn.conn, rspMsg)
	if err != nil {
		b.log.Error("send taskbuildrsp, err = %s\n", err)
	}
}
func (b *builder) handleBuildStateReq(reqMsg *myproto.Message) {
	reqPb := reqMsg.Body.(*myproto.TaskStateReq)
	sid := reqPb.GetSID()
	if b.sess.Get(sid) != nil {
		return
	}
	b.sess.Add(sid, reqMsg)
	defer b.sess.Del(sid)

	var taskRsp []*myproto.TaskStateRsp_TaskState

	tasks := reqPb.GetTask()
	for _, task := range tasks {
		groupname := task.GetGroupName()
		name := task.GetName()
		fetchall := task.GetFetchAll()

		if fetchall {
			states, err := b.db.gets(groupname, name)
			if err != nil {
				b.log.Error("db.gets err, g=%s, n=%s, err=%s\n", groupname, name, err)
				taskState := &myproto.TaskStateRsp_TaskState{
					GroupName:  proto.String(groupname),
					Name:       proto.String(name),
					BuildState: proto.Int32(util.KENotfound),
				}
				taskRsp = append(taskRsp, taskState)
				continue
			}
			for _, state := range states {
				taskState := &myproto.TaskStateRsp_TaskState{
					GroupName:  proto.String(groupname),
					Name:       proto.String(name),
					BuildState: proto.Int32(state.state),
					OutBin:     proto.String(state.outbin),
				}

				up, err := state.uptime.MarshalText()
				if err == nil {
					taskState.BuildTime = proto.String(string(up))
				}
				cont, err := ioutil.ReadFile(state.log)
				if err == nil {
					taskState.BuildLog = proto.String(string(cont))
				} else {
					taskState.BuildLog = proto.String(fmt.Sprintf("open log failed, path=%s", state.log))
				}
				taskRsp = append(taskRsp, taskState)
			}
		} else {
			state, err := b.db.get(groupname, name)
			if err != nil {
				b.log.Error("db.get err, g=%s, n=%s, err=%s\n", groupname, name, err)
				taskState := &myproto.TaskStateRsp_TaskState{
					GroupName:  proto.String(groupname),
					Name:       proto.String(name),
					BuildState: proto.Int32(util.KENotfound),
				}
				taskRsp = append(taskRsp, taskState)
				continue
			}
			taskState := &myproto.TaskStateRsp_TaskState{
				GroupName:  proto.String(groupname),
				Name:       proto.String(name),
				BuildState: proto.Int32(state.state),
				OutBin:     proto.String(state.outbin),
			}

			up, err := state.uptime.MarshalText()
			if err == nil {
				taskState.BuildTime = proto.String(string(up))
			}
			cont, err := ioutil.ReadFile(state.log)
			if err == nil {
				taskState.BuildLog = proto.String(string(cont))
			} else {
				taskState.BuildLog = proto.String(fmt.Sprintf("open log failed, path=%s", state.log))
			}
			taskRsp = append(taskRsp, taskState)
		}
	}

	rspMsg := &myproto.Message{}
	rspMsg.Head.Command = reqMsg.Head.Command + 1
	pb, err := myproto.GetMessage(rspMsg.Head.Command)
	if err != nil {
		b.log.Error("gen taskstatersp failed, err = %s\n", err)
		return
	}
	rspPb := pb.(*myproto.TaskStateRsp)
	rspPb.SID = proto.String(reqPb.GetSID())
	rspPb.RetCode = proto.Int32(util.KESuccess)
	rspPb.State = taskRsp

	rspMsg.Body = rspPb
	b.log.Info("RSP：\n%s\n", myproto.Debug(rspMsg))
	err = b.net.SendData(b.conn.conn, rspMsg)
	if err != nil {
		b.log.Error("send taskstatersp, err = %s\n", err)
	}
}
