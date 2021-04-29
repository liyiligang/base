// Copyright 2019 The Authors. All rights reserved.
// Author: liyiligang
// Date: 2020/2/13 19:28
// Description: 公共定义文件

package commonConst

const ConstBroadcast = -1
const ConstRandSend = -2
const ConstTimeFormat = "2006-01-02 15:04:05"
const ManageNodeID = 0
const ManageNodeTypeID = 0

type OrderType int32

const (
	UnKnowOrder     OrderType = 0
	WebsocketOrder  OrderType = 1
	RpcServerOrder  OrderType = 2
	HttpServerOrder OrderType = 3
	RpcClientOrder  OrderType = 4
)

type NodeTypeName string

const (
	ManageServerName  NodeTypeName = "管控服务"
	GatewayServerName NodeTypeName = "网关服务"
	LoginServerName   NodeTypeName = "登录服务"
	AppServerName     NodeTypeName = "应用服务"
	CommonServerName  NodeTypeName = "公共服务"
)

type CommonNodeData struct {
	NodeID       int32
	NodeTypeID   int32
	NodeTypeName NodeTypeName
	NodeName     string
	NodeGroup    string
	NodeState    int32
	PrivateAddr  string
	PublicAddr   string
	GrpcPort     string
	Data         interface{}
}

type GetwayNodeData struct {
	WebPort string
}
