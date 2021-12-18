/*
 * Copyright 2021 liyiligang.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package Jrpc

import (
	"context"
	"errors"
	"fmt"
	"github.com/gogo/protobuf/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"io"
	"sync/atomic"
	"time"
)

//metadata如果要传输二进制数据，key必须以bin结尾
const ConstRpcStreamClientHeader = "rpc-stream-client-header-bin"
const ConstRpcStreamServerHeader = "rpc-stream-server-header-bin"
const ConstRpcStreamServerTrailer = "rpc-stream-server-trailer-bin"

type RpcStreamCall struct {
	RpcStreamConnect  		func(conn *RpcStream) (interface{}, error)
	RpcStreamConnected  	func(conn *RpcStream) error
	RpcStreamClosed 		func(conn *RpcStream)
	RpcStreamReceiver   	func(conn *RpcStream, recv interface{})
	RpcStreamError 			func(text string, err error)
}

type RpcStream struct {
	send         chan interface{}
	sendCnt		 int64
	cancel       context.CancelFunc
	context      context.Context
	rpcContext   RpcContext
	rpcBindVal   interface{}
	recvMsgProto proto.Message
	call         RpcStreamCall
}

// 发送rpc数据
func (rpc *RpcStream) SendData(data interface{}) {
	atomic.AddInt64(&rpc.sendCnt, 1)
	rpc.send <- data
}

func (rpc *RpcStream) Close(immediately bool) {
	if !immediately {
		rpc.sendWait()
	}
	if rpc.cancel != nil {
		rpc.cancel()
	}
}

// 获取绑定值
func (rpc *RpcStream) GetBindVal() interface{} {
	return rpc.rpcBindVal
}

// 获取参数值
func (rpc *RpcStream) GetRpcContext() RpcContext {
	return rpc.rpcContext
}

// 设置rpc流客户端数据
func (rpc *RpcStream) WriteRpcStreamClientHeader(data []byte) {
	rpc.rpcContext.RpcStreamClientHeader = data
}

// 设置rpc流客户端数据
func (rpc *RpcStream) WriteRpcStreamServerHeader(data []byte) {
	rpc.rpcContext.RpcStreamServerHeader = data
}

// 初始化rpc流服务
func GrpcStreamServerInit(stream grpc.ServerStream, recvMsgProto proto.Message, call RpcStreamCall) (*RpcStream, error) {
	rpcContext, err := ParseRpcContext(stream.Context())
	if err != nil {
		return nil, err
	}
	childCtx, childCancel := context.WithCancel(stream.Context())
	conn := RpcStream{
		send:         make(chan interface{}, 256),
		cancel:       childCancel,
		context:      childCtx,
		rpcContext:   rpcContext,
		recvMsgProto: recvMsgProto,
		call:         call,
	}
	if conn.call.RpcStreamConnect != nil {
		id, err := conn.call.RpcStreamConnect(&conn)
		if err != nil {
			stream.SetTrailer(setRpcStreamServerTrailer([]byte(err.Error())))
			return nil, err
		}
		conn.rpcBindVal = id
	}
	err = stream.SendHeader(setRpcStreamServerHeader(conn.rpcContext.RpcStreamServerHeader))
	if err != nil {
		return nil, err
	}
	return &conn, nil
}

// 运行rpc流服务
func (rpc *RpcStream) GrpcStreamServerRun(stream grpc.ServerStream) error {
	go rpc.readMessage(rpc.context, stream.RecvMsg)
	go rpc.writeMessage(rpc.context, stream.SendMsg)

	var err error
	if rpc.call.RpcStreamConnected != nil {
		err = rpc.call.RpcStreamConnected(rpc)
		if err !=nil {
			stream.SetTrailer(setRpcStreamServerTrailer([]byte(err.Error())))
			rpc.Close(true)
		}
	}
	<-rpc.context.Done()
	if rpc.call.RpcStreamClosed != nil {
		rpc.call.RpcStreamClosed(rpc)
	}
	return err
}

//初始化rpc流客户端
func GrpcStreamClientInit(recvMsgProto proto.Message, call RpcStreamCall) (*RpcStream, error) {
	conn := &RpcStream{
		send:         make(chan interface{}, 256),
		rpcContext:   RpcContext{},
		recvMsgProto: recvMsgProto,
		call:         call,
	}
	if conn.call.RpcStreamConnect != nil {
		id, connErr := conn.call.RpcStreamConnect(conn)
		if connErr !=  nil  {
			return nil, connErr
		}
		conn.rpcBindVal = id
	}
	return conn, nil
}

//运行rpc流客户端
func (rpc *RpcStream) GrpcStreamClientRun(stream grpc.ClientStream) error {
	rpc.context, rpc.cancel = context.WithCancel(stream.Context())
	go rpc.readMessage(rpc.context, stream.RecvMsg)
	go rpc.writeMessage(rpc.context, stream.SendMsg)
	go func(){
		<-rpc.context.Done()
		err := stream.CloseSend()
		if err != nil {
			rpc.rpcStreamError("rpc close send error", err)
		}
		rpc.rpcContext.RpcStreamServerTrailer = getRpcStreamServerTrailer(stream.Trailer())
		if rpc.call.RpcStreamClosed != nil {
			rpc.call.RpcStreamClosed(rpc)
		}
	}()
	header, err := stream.Header()
	if err != nil {
		rpc.Close(true)
		return err
	}
	rpc.rpcContext.RpcStreamServerHeader = getRpcStreamServerHeader(header)
	if rpc.call.RpcStreamConnected != nil {
		err = rpc.call.RpcStreamConnected(rpc)
		if err != nil {
			rpc.Close(true)
			return err
		}
	}
	return nil
}

func (rpc *RpcStream) readMessage(ctx context.Context, recv func(m interface{}) error) {
	for {
		select {
		case <-ctx.Done():
			return
		default:
			err := recv(rpc.recvMsgProto)
			if err != nil {
				rpc.rpcStreamError("rpc read error", err)
				return
			}
			if rpc.call.RpcStreamReceiver != nil {
				rpc.call.RpcStreamReceiver(rpc, rpc.recvMsgProto)
			}
		}
	}
}

func (rpc *RpcStream) writeMessage(ctx context.Context, send func(m interface{}) error) {
	for {
		select {
		case <-ctx.Done():
			return
		case message, ok := <-rpc.send:
			if !ok {
				rpc.rpcStreamError("rpc write error", errors.New("chan is closed"))
				return
			}
			err := send(message)
			atomic.AddInt64(&rpc.sendCnt, -1)
			if err != nil {
				rpc.rpcStreamError("rpc write error", err)
				return
			}
		}
	}
}

//设置Rpc流元数据
func SetRpcStreamClientHeader(data []byte) context.Context {
	return metadata.AppendToOutgoingContext(context.Background(), ConstRpcStreamClientHeader, string(data))
}

//获取Rpc流元数据
func getRpcStreamClientHeader(md metadata.MD) []byte {
	s := md.Get(ConstRpcStreamClientHeader)
	if s != nil {
		return []byte(s[0])
	}
	return []byte("")
}

//设置Rpc流Header
func setRpcStreamServerHeader(data []byte) metadata.MD {
	return metadata.Pairs(ConstRpcStreamServerHeader, string(data))
}

//获取Rpc流Header
func getRpcStreamServerHeader(md metadata.MD) []byte {
	s := md.Get(ConstRpcStreamServerHeader)
	if s != nil {
		return []byte(s[0])
	}
	return []byte("")
}

//设置Rpc流Trailer
func setRpcStreamServerTrailer(data []byte) metadata.MD {
	return metadata.Pairs(ConstRpcStreamServerTrailer, string(data))
}

//获取Rpc流Trailer
func getRpcStreamServerTrailer(md metadata.MD) []byte {
	s := md.Get(ConstRpcStreamServerTrailer)
	if s != nil {
		return []byte(s[0])
	}
	return []byte("")
}

func (rpc *RpcStream) rpcStreamError(text string, err error) {
	if errors.Is(err, io.EOF){
		return
	}
	if rpc.call.RpcStreamError != nil {
		rpc.call.RpcStreamError(text, err)
	}else {
		fmt.Println(text + ": ", err)
	}
}

func (rpc *RpcStream) sendWait() {
	for {
		select {
		case <-rpc.context.Done():
			return
		default:
			if atomic.LoadInt64(&rpc.sendCnt) > 0 {
				time.Sleep(5 * time.Millisecond)
			}else{
				return
			}
		}
	}
}