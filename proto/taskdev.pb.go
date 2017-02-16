// Code generated by protoc-gen-go.
// source: taskdev.proto
// DO NOT EDIT!

/*
Package proto is a generated protocol buffer package.

It is generated from these files:
	taskdev.proto

It has these top-level messages:
	HeartBeatReq
	HeartBeatRsp
	RegisterReq
	RegisterRsp
	TaskBuildReq
	TaskBuildRsp
	TaskStateReq
	TaskStateRsp
*/
package proto

import proto1 "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto1.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto1.ProtoPackageIsVersion2 // please upgrade the proto package

// 心跳请求 0x00010001
type HeartBeatReq struct {
	SID              *string `protobuf:"bytes,1,req,name=SID" json:"SID,omitempty"`
	XXX_unrecognized []byte  `json:"-"`
}

func (m *HeartBeatReq) Reset()                    { *m = HeartBeatReq{} }
func (m *HeartBeatReq) String() string            { return proto1.CompactTextString(m) }
func (*HeartBeatReq) ProtoMessage()               {}
func (*HeartBeatReq) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{0} }

func (m *HeartBeatReq) GetSID() string {
	if m != nil && m.SID != nil {
		return *m.SID
	}
	return ""
}

// 心跳回应 0x00010002
type HeartBeatRsp struct {
	SID              *string `protobuf:"bytes,1,req,name=SID" json:"SID,omitempty"`
	XXX_unrecognized []byte  `json:"-"`
}

func (m *HeartBeatRsp) Reset()                    { *m = HeartBeatRsp{} }
func (m *HeartBeatRsp) String() string            { return proto1.CompactTextString(m) }
func (*HeartBeatRsp) ProtoMessage()               {}
func (*HeartBeatRsp) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{1} }

func (m *HeartBeatRsp) GetSID() string {
	if m != nil && m.SID != nil {
		return *m.SID
	}
	return ""
}

// 注册请求 0x00010003
type RegisterReq struct {
	SID              *string  `protobuf:"bytes,1,req,name=SID" json:"SID,omitempty"`
	Type             *string  `protobuf:"bytes,2,req,name=Type" json:"Type,omitempty"`
	OS               *string  `protobuf:"bytes,3,req,name=OS" json:"OS,omitempty"`
	Arch             *string  `protobuf:"bytes,4,req,name=Arch" json:"Arch,omitempty"`
	Host             *string  `protobuf:"bytes,5,req,name=Host" json:"Host,omitempty"`
	Adress           []string `protobuf:"bytes,6,rep,name=Adress" json:"Adress,omitempty"`
	XXX_unrecognized []byte   `json:"-"`
}

func (m *RegisterReq) Reset()                    { *m = RegisterReq{} }
func (m *RegisterReq) String() string            { return proto1.CompactTextString(m) }
func (*RegisterReq) ProtoMessage()               {}
func (*RegisterReq) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{2} }

func (m *RegisterReq) GetSID() string {
	if m != nil && m.SID != nil {
		return *m.SID
	}
	return ""
}

func (m *RegisterReq) GetType() string {
	if m != nil && m.Type != nil {
		return *m.Type
	}
	return ""
}

func (m *RegisterReq) GetOS() string {
	if m != nil && m.OS != nil {
		return *m.OS
	}
	return ""
}

func (m *RegisterReq) GetArch() string {
	if m != nil && m.Arch != nil {
		return *m.Arch
	}
	return ""
}

func (m *RegisterReq) GetHost() string {
	if m != nil && m.Host != nil {
		return *m.Host
	}
	return ""
}

func (m *RegisterReq) GetAdress() []string {
	if m != nil {
		return m.Adress
	}
	return nil
}

// 注册回应 0x00010004
type RegisterRsp struct {
	RetCode          *int32                  `protobuf:"varint,1,req,name=RetCode" json:"RetCode,omitempty"`
	SID              *string                 `protobuf:"bytes,2,req,name=SID" json:"SID,omitempty"`
	Info             *RegisterRsp_ServerInfo `protobuf:"bytes,3,opt,name=info" json:"info,omitempty"`
	XXX_unrecognized []byte                  `json:"-"`
}

func (m *RegisterRsp) Reset()                    { *m = RegisterRsp{} }
func (m *RegisterRsp) String() string            { return proto1.CompactTextString(m) }
func (*RegisterRsp) ProtoMessage()               {}
func (*RegisterRsp) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{3} }

func (m *RegisterRsp) GetRetCode() int32 {
	if m != nil && m.RetCode != nil {
		return *m.RetCode
	}
	return 0
}

func (m *RegisterRsp) GetSID() string {
	if m != nil && m.SID != nil {
		return *m.SID
	}
	return ""
}

func (m *RegisterRsp) GetInfo() *RegisterRsp_ServerInfo {
	if m != nil {
		return m.Info
	}
	return nil
}

type RegisterRsp_ServerInfo struct {
	OS               *string  `protobuf:"bytes,1,req,name=OS" json:"OS,omitempty"`
	Arch             *string  `protobuf:"bytes,2,req,name=Arch" json:"Arch,omitempty"`
	Host             *string  `protobuf:"bytes,3,req,name=Host" json:"Host,omitempty"`
	Adress           []string `protobuf:"bytes,4,rep,name=Adress" json:"Adress,omitempty"`
	BuildFtpHost     *string  `protobuf:"bytes,5,opt,name=BuildFtpHost" json:"BuildFtpHost,omitempty"`
	BuildFtpUser     *string  `protobuf:"bytes,6,opt,name=BuildFtpUser" json:"BuildFtpUser,omitempty"`
	BuildFtpPass     *string  `protobuf:"bytes,7,opt,name=BuildFtpPass" json:"BuildFtpPass,omitempty"`
	XXX_unrecognized []byte   `json:"-"`
}

func (m *RegisterRsp_ServerInfo) Reset()                    { *m = RegisterRsp_ServerInfo{} }
func (m *RegisterRsp_ServerInfo) String() string            { return proto1.CompactTextString(m) }
func (*RegisterRsp_ServerInfo) ProtoMessage()               {}
func (*RegisterRsp_ServerInfo) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{3, 0} }

func (m *RegisterRsp_ServerInfo) GetOS() string {
	if m != nil && m.OS != nil {
		return *m.OS
	}
	return ""
}

func (m *RegisterRsp_ServerInfo) GetArch() string {
	if m != nil && m.Arch != nil {
		return *m.Arch
	}
	return ""
}

func (m *RegisterRsp_ServerInfo) GetHost() string {
	if m != nil && m.Host != nil {
		return *m.Host
	}
	return ""
}

func (m *RegisterRsp_ServerInfo) GetAdress() []string {
	if m != nil {
		return m.Adress
	}
	return nil
}

func (m *RegisterRsp_ServerInfo) GetBuildFtpHost() string {
	if m != nil && m.BuildFtpHost != nil {
		return *m.BuildFtpHost
	}
	return ""
}

func (m *RegisterRsp_ServerInfo) GetBuildFtpUser() string {
	if m != nil && m.BuildFtpUser != nil {
		return *m.BuildFtpUser
	}
	return ""
}

func (m *RegisterRsp_ServerInfo) GetBuildFtpPass() string {
	if m != nil && m.BuildFtpPass != nil {
		return *m.BuildFtpPass
	}
	return ""
}

// 编译请求 0x00010005
type TaskBuildReq struct {
	SID              *string                   `protobuf:"bytes,1,req,name=SID" json:"SID,omitempty"`
	Task             []*TaskBuildReq_TaskBuild `protobuf:"bytes,2,rep,name=Task" json:"Task,omitempty"`
	XXX_unrecognized []byte                    `json:"-"`
}

func (m *TaskBuildReq) Reset()                    { *m = TaskBuildReq{} }
func (m *TaskBuildReq) String() string            { return proto1.CompactTextString(m) }
func (*TaskBuildReq) ProtoMessage()               {}
func (*TaskBuildReq) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{4} }

func (m *TaskBuildReq) GetSID() string {
	if m != nil && m.SID != nil {
		return *m.SID
	}
	return ""
}

func (m *TaskBuildReq) GetTask() []*TaskBuildReq_TaskBuild {
	if m != nil {
		return m.Task
	}
	return nil
}

type TaskBuildReq_TaskBuild struct {
	Name             *string  `protobuf:"bytes,1,req,name=Name" json:"Name,omitempty"`
	Version          *string  `protobuf:"bytes,2,req,name=Version" json:"Version,omitempty"`
	GroupName        *string  `protobuf:"bytes,3,opt,name=GroupName" json:"GroupName,omitempty"`
	Vcs              *string  `protobuf:"bytes,4,req,name=Vcs" json:"Vcs,omitempty"`
	Repos            *string  `protobuf:"bytes,5,req,name=Repos" json:"Repos,omitempty"`
	ReposDir         *string  `protobuf:"bytes,6,opt,name=ReposDir" json:"ReposDir,omitempty"`
	User             *string  `protobuf:"bytes,7,opt,name=User" json:"User,omitempty"`
	Pass             *string  `protobuf:"bytes,8,opt,name=Pass" json:"Pass,omitempty"`
	PreBuildScript   *string  `protobuf:"bytes,9,opt,name=PreBuildScript" json:"PreBuildScript,omitempty"`
	BuildScript      *string  `protobuf:"bytes,10,opt,name=BuildScript" json:"BuildScript,omitempty"`
	PostBuildScript  *string  `protobuf:"bytes,11,opt,name=PostBuildScript" json:"PostBuildScript,omitempty"`
	OutBin           []string `protobuf:"bytes,12,rep,name=OutBin" json:"OutBin,omitempty"`
	OutFtpHost       *string  `protobuf:"bytes,13,opt,name=OutFtpHost" json:"OutFtpHost,omitempty"`
	OutFtpUser       *string  `protobuf:"bytes,14,opt,name=OutFtpUser" json:"OutFtpUser,omitempty"`
	OutFtpPass       *string  `protobuf:"bytes,15,opt,name=OutFtpPass" json:"OutFtpPass,omitempty"`
	TestScript       *string  `protobuf:"bytes,16,opt,name=TestScript" json:"TestScript,omitempty"`
	XXX_unrecognized []byte   `json:"-"`
}

func (m *TaskBuildReq_TaskBuild) Reset()                    { *m = TaskBuildReq_TaskBuild{} }
func (m *TaskBuildReq_TaskBuild) String() string            { return proto1.CompactTextString(m) }
func (*TaskBuildReq_TaskBuild) ProtoMessage()               {}
func (*TaskBuildReq_TaskBuild) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{4, 0} }

func (m *TaskBuildReq_TaskBuild) GetName() string {
	if m != nil && m.Name != nil {
		return *m.Name
	}
	return ""
}

func (m *TaskBuildReq_TaskBuild) GetVersion() string {
	if m != nil && m.Version != nil {
		return *m.Version
	}
	return ""
}

func (m *TaskBuildReq_TaskBuild) GetGroupName() string {
	if m != nil && m.GroupName != nil {
		return *m.GroupName
	}
	return ""
}

func (m *TaskBuildReq_TaskBuild) GetVcs() string {
	if m != nil && m.Vcs != nil {
		return *m.Vcs
	}
	return ""
}

func (m *TaskBuildReq_TaskBuild) GetRepos() string {
	if m != nil && m.Repos != nil {
		return *m.Repos
	}
	return ""
}

func (m *TaskBuildReq_TaskBuild) GetReposDir() string {
	if m != nil && m.ReposDir != nil {
		return *m.ReposDir
	}
	return ""
}

func (m *TaskBuildReq_TaskBuild) GetUser() string {
	if m != nil && m.User != nil {
		return *m.User
	}
	return ""
}

func (m *TaskBuildReq_TaskBuild) GetPass() string {
	if m != nil && m.Pass != nil {
		return *m.Pass
	}
	return ""
}

func (m *TaskBuildReq_TaskBuild) GetPreBuildScript() string {
	if m != nil && m.PreBuildScript != nil {
		return *m.PreBuildScript
	}
	return ""
}

func (m *TaskBuildReq_TaskBuild) GetBuildScript() string {
	if m != nil && m.BuildScript != nil {
		return *m.BuildScript
	}
	return ""
}

func (m *TaskBuildReq_TaskBuild) GetPostBuildScript() string {
	if m != nil && m.PostBuildScript != nil {
		return *m.PostBuildScript
	}
	return ""
}

func (m *TaskBuildReq_TaskBuild) GetOutBin() []string {
	if m != nil {
		return m.OutBin
	}
	return nil
}

func (m *TaskBuildReq_TaskBuild) GetOutFtpHost() string {
	if m != nil && m.OutFtpHost != nil {
		return *m.OutFtpHost
	}
	return ""
}

func (m *TaskBuildReq_TaskBuild) GetOutFtpUser() string {
	if m != nil && m.OutFtpUser != nil {
		return *m.OutFtpUser
	}
	return ""
}

func (m *TaskBuildReq_TaskBuild) GetOutFtpPass() string {
	if m != nil && m.OutFtpPass != nil {
		return *m.OutFtpPass
	}
	return ""
}

func (m *TaskBuildReq_TaskBuild) GetTestScript() string {
	if m != nil && m.TestScript != nil {
		return *m.TestScript
	}
	return ""
}

// 编译应答 0x00010006
type TaskBuildRsp struct {
	RetCode          *int32  `protobuf:"varint,1,req,name=RetCode" json:"RetCode,omitempty"`
	SID              *string `protobuf:"bytes,2,req,name=SID" json:"SID,omitempty"`
	XXX_unrecognized []byte  `json:"-"`
}

func (m *TaskBuildRsp) Reset()                    { *m = TaskBuildRsp{} }
func (m *TaskBuildRsp) String() string            { return proto1.CompactTextString(m) }
func (*TaskBuildRsp) ProtoMessage()               {}
func (*TaskBuildRsp) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{5} }

func (m *TaskBuildRsp) GetRetCode() int32 {
	if m != nil && m.RetCode != nil {
		return *m.RetCode
	}
	return 0
}

func (m *TaskBuildRsp) GetSID() string {
	if m != nil && m.SID != nil {
		return *m.SID
	}
	return ""
}

// 编译状态请求 0x00010007
type TaskStateReq struct {
	SID              *string                `protobuf:"bytes,1,req,name=SID" json:"SID,omitempty"`
	Task             []*TaskStateReq_TaskID `protobuf:"bytes,2,rep,name=Task" json:"Task,omitempty"`
	XXX_unrecognized []byte                 `json:"-"`
}

func (m *TaskStateReq) Reset()                    { *m = TaskStateReq{} }
func (m *TaskStateReq) String() string            { return proto1.CompactTextString(m) }
func (*TaskStateReq) ProtoMessage()               {}
func (*TaskStateReq) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{6} }

func (m *TaskStateReq) GetSID() string {
	if m != nil && m.SID != nil {
		return *m.SID
	}
	return ""
}

func (m *TaskStateReq) GetTask() []*TaskStateReq_TaskID {
	if m != nil {
		return m.Task
	}
	return nil
}

type TaskStateReq_TaskID struct {
	TaskName         *string `protobuf:"bytes,1,req,name=TaskName" json:"TaskName,omitempty"`
	TaskGroupName    *string `protobuf:"bytes,2,opt,name=TaskGroupName" json:"TaskGroupName,omitempty"`
	XXX_unrecognized []byte  `json:"-"`
}

func (m *TaskStateReq_TaskID) Reset()                    { *m = TaskStateReq_TaskID{} }
func (m *TaskStateReq_TaskID) String() string            { return proto1.CompactTextString(m) }
func (*TaskStateReq_TaskID) ProtoMessage()               {}
func (*TaskStateReq_TaskID) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{6, 0} }

func (m *TaskStateReq_TaskID) GetTaskName() string {
	if m != nil && m.TaskName != nil {
		return *m.TaskName
	}
	return ""
}

func (m *TaskStateReq_TaskID) GetTaskGroupName() string {
	if m != nil && m.TaskGroupName != nil {
		return *m.TaskGroupName
	}
	return ""
}

// 编译状态应答 0x00010008
type TaskStateRsp struct {
	RetCode          *int32                    `protobuf:"varint,1,req,name=RetCode" json:"RetCode,omitempty"`
	SID              *string                   `protobuf:"bytes,2,req,name=SID" json:"SID,omitempty"`
	State            []*TaskStateRsp_TaskState `protobuf:"bytes,3,rep,name=State" json:"State,omitempty"`
	XXX_unrecognized []byte                    `json:"-"`
}

func (m *TaskStateRsp) Reset()                    { *m = TaskStateRsp{} }
func (m *TaskStateRsp) String() string            { return proto1.CompactTextString(m) }
func (*TaskStateRsp) ProtoMessage()               {}
func (*TaskStateRsp) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{7} }

func (m *TaskStateRsp) GetRetCode() int32 {
	if m != nil && m.RetCode != nil {
		return *m.RetCode
	}
	return 0
}

func (m *TaskStateRsp) GetSID() string {
	if m != nil && m.SID != nil {
		return *m.SID
	}
	return ""
}

func (m *TaskStateRsp) GetState() []*TaskStateRsp_TaskState {
	if m != nil {
		return m.State
	}
	return nil
}

type TaskStateRsp_TaskState struct {
	TaskName         *string  `protobuf:"bytes,1,req,name=TaskName" json:"TaskName,omitempty"`
	TaskGroupName    *string  `protobuf:"bytes,2,opt,name=TaskGroupName" json:"TaskGroupName,omitempty"`
	BuildState       *int32   `protobuf:"varint,3,req,name=BuildState" json:"BuildState,omitempty"`
	BuildLog         *string  `protobuf:"bytes,4,req,name=BuildLog" json:"BuildLog,omitempty"`
	OutBin           []string `protobuf:"bytes,5,rep,name=OutBin" json:"OutBin,omitempty"`
	TestLog          *string  `protobuf:"bytes,6,opt,name=TestLog" json:"TestLog,omitempty"`
	XXX_unrecognized []byte   `json:"-"`
}

func (m *TaskStateRsp_TaskState) Reset()                    { *m = TaskStateRsp_TaskState{} }
func (m *TaskStateRsp_TaskState) String() string            { return proto1.CompactTextString(m) }
func (*TaskStateRsp_TaskState) ProtoMessage()               {}
func (*TaskStateRsp_TaskState) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{7, 0} }

func (m *TaskStateRsp_TaskState) GetTaskName() string {
	if m != nil && m.TaskName != nil {
		return *m.TaskName
	}
	return ""
}

func (m *TaskStateRsp_TaskState) GetTaskGroupName() string {
	if m != nil && m.TaskGroupName != nil {
		return *m.TaskGroupName
	}
	return ""
}

func (m *TaskStateRsp_TaskState) GetBuildState() int32 {
	if m != nil && m.BuildState != nil {
		return *m.BuildState
	}
	return 0
}

func (m *TaskStateRsp_TaskState) GetBuildLog() string {
	if m != nil && m.BuildLog != nil {
		return *m.BuildLog
	}
	return ""
}

func (m *TaskStateRsp_TaskState) GetOutBin() []string {
	if m != nil {
		return m.OutBin
	}
	return nil
}

func (m *TaskStateRsp_TaskState) GetTestLog() string {
	if m != nil && m.TestLog != nil {
		return *m.TestLog
	}
	return ""
}

func init() {
	proto1.RegisterType((*HeartBeatReq)(nil), "proto.HeartBeatReq")
	proto1.RegisterType((*HeartBeatRsp)(nil), "proto.HeartBeatRsp")
	proto1.RegisterType((*RegisterReq)(nil), "proto.RegisterReq")
	proto1.RegisterType((*RegisterRsp)(nil), "proto.RegisterRsp")
	proto1.RegisterType((*RegisterRsp_ServerInfo)(nil), "proto.RegisterRsp.ServerInfo")
	proto1.RegisterType((*TaskBuildReq)(nil), "proto.TaskBuildReq")
	proto1.RegisterType((*TaskBuildReq_TaskBuild)(nil), "proto.TaskBuildReq.TaskBuild")
	proto1.RegisterType((*TaskBuildRsp)(nil), "proto.TaskBuildRsp")
	proto1.RegisterType((*TaskStateReq)(nil), "proto.TaskStateReq")
	proto1.RegisterType((*TaskStateReq_TaskID)(nil), "proto.TaskStateReq.TaskID")
	proto1.RegisterType((*TaskStateRsp)(nil), "proto.TaskStateRsp")
	proto1.RegisterType((*TaskStateRsp_TaskState)(nil), "proto.TaskStateRsp.TaskState")
}

func init() { proto1.RegisterFile("taskdev.proto", fileDescriptor0) }

var fileDescriptor0 = []byte{
	// 514 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x09, 0x6e, 0x88, 0x02, 0xff, 0x8c, 0x53, 0x4d, 0x6f, 0xd3, 0x40,
	0x10, 0x95, 0xed, 0x38, 0x69, 0x26, 0x4e, 0x52, 0xc2, 0x97, 0x65, 0x84, 0x54, 0xe5, 0x14, 0x89,
	0x2a, 0x12, 0xfc, 0x83, 0x86, 0x08, 0x1a, 0x09, 0xd1, 0x28, 0x0e, 0x3d, 0x63, 0x25, 0x43, 0xb1,
	0x80, 0xd8, 0xec, 0x6e, 0x2a, 0x10, 0x77, 0xfe, 0x0d, 0xbf, 0x87, 0x33, 0xff, 0x82, 0x23, 0xb3,
	0xe3, 0xdd, 0x7a, 0x9d, 0x1e, 0xca, 0x29, 0x7e, 0x6f, 0x66, 0x77, 0xde, 0xbc, 0xb7, 0x81, 0xbe,
	0xca, 0xe4, 0xa7, 0x2d, 0x5e, 0x4f, 0x4b, 0x51, 0xa8, 0x62, 0x14, 0xf2, 0xcf, 0xf8, 0x09, 0x44,
	0xe7, 0x98, 0x09, 0x35, 0xc3, 0x4c, 0xad, 0xf0, 0xeb, 0xa8, 0x07, 0x41, 0xba, 0x98, 0xc7, 0xde,
	0x89, 0x3f, 0xe9, 0x36, 0x8b, 0xb2, 0x6c, 0x16, 0xdf, 0x43, 0x6f, 0x85, 0x57, 0xb9, 0x54, 0x28,
	0x0e, 0x0f, 0x8e, 0x22, 0x68, 0xad, 0xbf, 0x97, 0x18, 0xfb, 0x8c, 0x00, 0xfc, 0x8b, 0x34, 0x0e,
	0x6c, 0xe5, 0x4c, 0x6c, 0x3e, 0xc6, 0x2d, 0x8b, 0xce, 0x0b, 0xa9, 0xe2, 0x90, 0xd1, 0x00, 0xda,
	0x67, 0x5b, 0x81, 0x52, 0xc6, 0xed, 0x93, 0x80, 0x26, 0xfc, 0xf1, 0x9c, 0x11, 0x34, 0x7e, 0x08,
	0x9d, 0x15, 0xaa, 0x97, 0xc5, 0x16, 0x79, 0x4c, 0x68, 0x67, 0x56, 0x53, 0x9e, 0x41, 0x2b, 0xdf,
	0x7d, 0x28, 0x68, 0x8e, 0x37, 0xe9, 0xbd, 0x78, 0x5a, 0xad, 0x39, 0x75, 0xce, 0x4f, 0x53, 0x14,
	0xd7, 0x28, 0x16, 0xd4, 0x94, 0xfc, 0xf4, 0x00, 0x6a, 0x68, 0x14, 0x7a, 0x0d, 0x85, 0x7e, 0x43,
	0x61, 0x70, 0xa0, 0xb0, 0xa5, 0x15, 0x8e, 0x1e, 0x40, 0x34, 0xdb, 0xe7, 0x9f, 0xb7, 0xaf, 0x54,
	0x69, 0xf6, 0xf0, 0x9a, 0xec, 0x3b, 0x89, 0x82, 0xb6, 0x39, 0x60, 0x97, 0x19, 0xdd, 0xd0, 0xd1,
	0xec, 0xf8, 0xaf, 0x0f, 0xd1, 0x9a, 0x82, 0xe1, 0xd2, 0x2d, 0x1f, 0x69, 0x27, 0x5d, 0x24, 0x2d,
	0x81, 0xb3, 0x93, 0xdb, 0x5f, 0x83, 0xe4, 0x97, 0x0f, 0xdd, 0x1b, 0xa4, 0x85, 0xbf, 0xcd, 0xbe,
	0xa0, 0xb9, 0x88, 0xac, 0xbb, 0x44, 0x21, 0xf3, 0x62, 0x67, 0xf6, 0xba, 0x07, 0xdd, 0xd7, 0xa2,
	0xd8, 0x97, 0xdc, 0x13, 0xb0, 0x40, 0x9a, 0x7c, 0xb9, 0x91, 0x26, 0x99, 0x3e, 0x84, 0x2b, 0x2c,
	0x0b, 0x69, 0xa2, 0x39, 0x86, 0x23, 0x86, 0xf3, 0xdc, 0xae, 0x43, 0xf7, 0xf3, 0x72, 0x1d, 0x8b,
	0x78, 0xa9, 0x23, 0x46, 0x8f, 0x60, 0xb0, 0x14, 0xc8, 0x3a, 0xd2, 0x8d, 0xc8, 0x4b, 0x15, 0x77,
	0x99, 0xbf, 0x0f, 0x3d, 0x97, 0x04, 0x26, 0x1f, 0xc3, 0x70, 0x49, 0xde, 0xb9, 0x85, 0x1e, 0x17,
	0xc8, 0xec, 0x8b, 0xbd, 0x9a, 0xe5, 0xbb, 0x38, 0x62, 0xb3, 0x29, 0x25, 0xc2, 0xd6, 0xea, 0x3e,
	0xf7, 0xdc, 0x70, 0xac, 0x65, 0xd0, 0xe4, 0x58, 0xd1, 0xd0, 0x72, 0x6b, 0x94, 0xca, 0xdc, 0x7f,
	0xcc, 0xd6, 0x9f, 0xba, 0xce, 0xdf, 0xf5, 0xbc, 0xc6, 0xdf, 0xaa, 0xee, 0x54, 0x65, 0x0a, 0x6f,
	0xe5, 0x34, 0x69, 0xe4, 0x94, 0x38, 0x39, 0xd9, 0x7e, 0x06, 0x8b, 0x79, 0xf2, 0x1c, 0xda, 0xd5,
	0x97, 0xb6, 0x54, 0x7f, 0x39, 0x21, 0x3d, 0x84, 0xbe, 0x66, 0xea, 0x5c, 0x7c, 0xd6, 0xf9, 0xdb,
	0x73, 0x47, 0xdf, 0xf9, 0x3f, 0x38, 0x85, 0x90, 0x3b, 0x29, 0xd5, 0xc3, 0x47, 0x63, 0x6f, 0xa8,
	0x41, 0xf2, 0xa3, 0x7a, 0x33, 0x0c, 0xfe, 0x5b, 0x92, 0xb6, 0xb3, 0xca, 0xcb, 0x0c, 0xd2, 0x22,
	0xe8, 0x30, 0x73, 0x6f, 0x8a, 0x2b, 0xf3, 0x86, 0xea, 0x00, 0x43, 0x0e, 0x90, 0x74, 0xeb, 0x10,
	0x74, 0x03, 0xbf, 0xa1, 0x7f, 0x01, 0x00, 0x00, 0xff, 0xff, 0xf2, 0xaa, 0xa7, 0xb0, 0x93, 0x04,
	0x00, 0x00,
}
