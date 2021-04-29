// Copyright 2019 The Authors. All rights reserved.
// Author: liyiligang
// Date: 2019/4/2 10:10
// Description: rpc客户端

package Jrpc

import (
	"context"
	"errors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"io"
)

//metadata如果要传输二进制数据，key必须以bin结尾
const ConstRpcStreamClientMetadata = "rpc-stream-metadata-bin"
const ConstRpcStreamServerHeader = "rpc-stream-header-bin"
const ConstRpcStreamServerTrailer = "rpc-stream-trailer-bin"

type RpcStreamFunc interface {
	RpcStreamConnect(conn *RpcStream) (interface{}, error)
	RpcStreamConnected(conn *RpcStream) error
	RpcStreamClose(conn *RpcStream)
	RpcStreamReceiver(conn *RpcStream, recv interface{})
	RpcStreamError(text string, err error)
}

type RpcStream struct {
	send         chan interface{}
	cancel       context.CancelFunc
	context      context.Context
	rpcParm      RpcParm
	rpcBindVal   interface{}
	recvMsgProto interface{}
	call         RpcStreamFunc
}

// 发送rpc数据
func (rpc *RpcStream) SendData(data interface{}) {
	rpc.send <- data
}

func (rpc *RpcStream) Close() {
	if rpc.cancel != nil {
		rpc.cancel()
	}
}

// 获取绑定值
func (rpc *RpcStream) GetBindVal() interface{} {
	return rpc.rpcBindVal
}

// 获取参数值
func (rpc *RpcStream) GetParm() RpcParm {
	return rpc.rpcParm
}

// 设置rpc流客户端数据
func (rpc *RpcStream) SetRpcStreamClientMsg(data []byte) {
	rpc.rpcParm.RpcStreamClientMsg = data
}

// 设置rpc流客户端数据
func (rpc *RpcStream) SetRpcStreamServerHeader(data []byte) {
	rpc.rpcParm.RpcStreamServerHeader = data
}

// 获取Context
func (rpc *RpcStream) GetContext() context.Context {
	return rpc.context
}

// 初始化rpc流服务
func GrpcStreamServerInit(stream grpc.ServerStream, recvMsgProto interface{}, call RpcStreamFunc) (*RpcStream, error) {
	rpcParm, pErr := GetConnectParm(stream.Context())
	if pErr != nil {
		return nil, pErr
	}
	childCtx, cancel := context.WithCancel(stream.Context())
	conn := RpcStream{
		send:         make(chan interface{}, 256),
		cancel:       cancel,
		context:      childCtx,
		rpcParm:      rpcParm,
		recvMsgProto: recvMsgProto,
		call:         call,
	}
	id, connErr := conn.call.RpcStreamConnect(&conn)
	if connErr != nil {
		stream.SetTrailer(SetRpcStreamServerTrailer([]byte(connErr.Error())))
		return nil, connErr
	}
	conn.rpcBindVal = id
	stream.SendHeader(SetRpcStreamServerHeader(conn.rpcParm.RpcStreamServerHeader))
	return &conn, nil
}

// 运行rpc流服务
func (rpc *RpcStream) GrpcStreamServerRun(stream grpc.ServerStream) error {
	go rpc.readMessage(stream.Context(), stream.RecvMsg)
	go rpc.writeMessage(stream.Context(), stream.SendMsg)
	err := rpc.call.RpcStreamConnected(rpc)
	if err !=nil {
		stream.SetTrailer(SetRpcStreamServerTrailer([]byte(err.Error())))
		rpc.Close()
	}
	<-rpc.context.Done()
	rpc.call.RpcStreamClose(rpc)
	return err
}

//初始化rpc流客户端
func GrpcStreamClientInit(recvMsgProto interface{}, call RpcStreamFunc) (*RpcStream, error) {
	conn := &RpcStream{
		send:         make(chan interface{}, 256),
		rpcParm:      RpcParm{},
		recvMsgProto: recvMsgProto,
		call:         call,
	}
	id, connErr := conn.call.RpcStreamConnect(conn)
	if connErr !=  nil  {
		return nil, connErr
	}
	conn.rpcBindVal = id
	conn.context, conn.cancel = context.WithCancel(context.Background())
	conn.context = SetRpcStreamClientMetadata(conn.context, conn.rpcParm.RpcStreamClientMsg)
	return conn, nil
}

//运行rpc流客户端
func (rpc *RpcStream) GrpcStreamClientRun(stream grpc.ClientStream) error {
	go rpc.readMessage(stream.Context(), stream.RecvMsg)
	go rpc.writeMessage(stream.Context(), stream.SendMsg)
	go func(){
		<-stream.Context().Done()
		stream.CloseSend()
		rpc.rpcParm.RpcStreamServerTrailer = GetRpcStreamServerTrailer(stream.Trailer())
		rpc.call.RpcStreamClose(rpc)
	}()
	header, err := stream.Header()
	if err != nil {
		rpc.Close()
		return err
	}
	rpc.rpcParm.RpcStreamServerHeader = GetRpcStreamServerHeader(header)
	err = rpc.call.RpcStreamConnected(rpc)
	if err != nil {
		rpc.Close()
		return err
	}
	return nil
}

func (rpc *RpcStream) readMessage(ctx context.Context, recv func(m interface{}) error) {
	for {
		select {
		case <-ctx.Done(): //rpc流结束
			return
		default:
			err := recv(rpc.recvMsgProto)
			if err != nil {
				if !errors.Is(err, io.EOF){
					rpc.call.RpcStreamError("rpc readMessage 关闭", err)
				}
				return
			}
			rpc.call.RpcStreamReceiver(rpc, rpc.recvMsgProto)
		}
	}
}

func (rpc *RpcStream) writeMessage(ctx context.Context, send func(m interface{}) error) {
	for {
		select {
		case <-ctx.Done(): //rpc流结束
			return
		case message, ok := <-rpc.send:
			if !ok {
				rpc.call.RpcStreamError("rpc send读取错误", nil)
				return
			}
			err := send(message)
			if err != nil {
				if !errors.Is(err, io.EOF){
					rpc.call.RpcStreamError("rpc sendMessage 关闭",  err)
				}
				return
			}
		}
	}
}

//设置Rpc流元数据
func SetRpcStreamClientMetadata(ctx context.Context, data []byte) context.Context {
	return metadata.AppendToOutgoingContext(ctx, ConstRpcStreamClientMetadata, string(data))
}

//获取Rpc流元数据
func GetRpcStreamMetadata(md metadata.MD) []byte {
	s := md.Get(ConstRpcStreamClientMetadata)
	if s != nil {
		return []byte(s[0])
	}
	return []byte("")
}

//设置Rpc流Header
func SetRpcStreamServerHeader(data []byte) metadata.MD {
	return metadata.Pairs(ConstRpcStreamServerHeader, string(data))
}

//获取Rpc流Header
func GetRpcStreamServerHeader(md metadata.MD) []byte {
	s := md.Get(ConstRpcStreamServerHeader)
	if s != nil {
		return []byte(s[0])
	}
	return []byte("")
}

//设置Rpc流Trailer
func SetRpcStreamServerTrailer(data []byte) metadata.MD {
	return metadata.Pairs(ConstRpcStreamServerTrailer, string(data))
}

//获取Rpc流Trailer
func GetRpcStreamServerTrailer(md metadata.MD) []byte {
	s := md.Get(ConstRpcStreamServerTrailer)
	if s != nil {
		return []byte(s[0])
	}
	return []byte("")
}

