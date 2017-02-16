package main

import (
	myproto "taskdev/proto"
	"taskdev/util"

	"github.com/golang/protobuf/proto"
)

func (s *scheServer) addBuildTask(b *buildScheData) {
	s.buildScheChan <- b
}
func (s *scheServer) stopBuildTask(b *buildScheData) {
	b.stop = true
}
func (s *scheServer) sendBuildTask(b *buildScheData) error {
	reqMsg := &myproto.Message{}
	reqMsg.Head.Command = myproto.KCmdTaskBuildReq
	reqPb, err := myproto.GetMessage(reqMsg.Head.Command)
	if err != nil {
		s.log.Error("gen KCmdTaskBuildReq failed, err = %s\n", err)
		return err
	}
	pb := reqPb.(*myproto.TaskBuildReq)
	pb.SID = proto.String(util.GetSID(16))

	var tasks []*myproto.TaskBuildReq_TaskBuild
	if b.tasks != nil {
		for _, t := range b.tasks.Tasks {
			tpb := &myproto.TaskBuildReq_TaskBuild{}
			tpb.Name = proto.String(t.Name)
			tpb.Version = proto.String(t.Version)
			tpb.GroupName = proto.String(b.tasks.Name)
			tpb.Vcs = proto.String(t.Vcs)
			tpb.Repos = proto.String(t.Repos)
			tpb.ReposDir = proto.String(t.ReposDir)
			tpb.User = proto.String(t.User)
			tpb.Pass = proto.String(t.Pass)
			tpb.PreBuildScript = proto.String(t.PreBuildScript)
			tpb.BuildScript = proto.String(t.BuildScript)
			tpb.PostBuildScript = proto.String(t.PostBuildScript)
			tpb.OutBin = t.OutBin
			tpb.OutFtpHost = proto.String(t.OutFtpHost)
			tpb.OutFtpUser = proto.String(t.OutFtpUser)
			tpb.OutFtpPass = proto.String(t.OutFtpPass)
			tpb.TestScript = proto.String(t.TestScript)

			tasks = append(tasks, tpb)
		}
	}
	if b.task != nil {
		t := b.task
		tpb := &myproto.TaskBuildReq_TaskBuild{}
		tpb.Name = proto.String(t.Name)
		tpb.Version = proto.String(t.Version)
		if b.task.Group != nil {
			tpb.GroupName = proto.String(b.task.Group.Name)
		}
		tpb.Vcs = proto.String(t.Vcs)
		tpb.Repos = proto.String(t.Repos)
		tpb.ReposDir = proto.String(t.ReposDir)
		tpb.User = proto.String(t.User)
		tpb.Pass = proto.String(t.Pass)
		tpb.PreBuildScript = proto.String(t.PreBuildScript)
		tpb.BuildScript = proto.String(t.BuildScript)
		tpb.PostBuildScript = proto.String(t.PostBuildScript)
		tpb.OutBin = t.OutBin
		tpb.OutFtpHost = proto.String(t.OutFtpHost)
		tpb.OutFtpUser = proto.String(t.OutFtpUser)
		tpb.OutFtpPass = proto.String(t.OutFtpPass)
		tpb.TestScript = proto.String(t.TestScript)

		tasks = append(tasks, tpb)
	}

	pb.Task = tasks

	reqMsg.Body = reqPb
	s.log.Info("REQ：\n%s\n", myproto.Debug(reqMsg))
	err = s.net.SendData(b.conn.conn, reqMsg)
	if err != nil {
		s.log.Error("send TaskBuildReq, err = %s\n", err)
		return err
	}
	return nil
}

func (s *scheServer) schedule() {
	for {
		select {
		case data, ok := <-s.buildScheChan:
			if !ok {
				return
			}

			s.sendBuildTask(data)

		}

	}
}
