// Copyright 2021 The Authors. All rights reserved.
// Author: liyiligang
// Date: 2021/06/28 11:11
// Description:

package Jrpc

import (
	"context"
	"errors"
	"github.com/liyiligang/base/commonConst"
	"google.golang.org/grpc"
	"google.golang.org/grpc/connectivity"
	"time"
)

type RpcKeepaliveCall interface {
	RpcServeConnected(client *RpcKeepalive, isReConnect bool)
	RpcServeDisconnected(client *RpcKeepalive, isCloseByUser bool)
	RpcKeepaliveError(text string, err error)
}

type RpcKeepalive struct {
	ServerNode      	 *commonConst.CommonNodeData
	Conn                 *grpc.ClientConn
	KeepaliveTime 	     time.Duration
	call 			     RpcKeepaliveCall
	cancel  			 context.CancelFunc
}

func RegisterRpcKeepalive(rpcKeepalive *RpcKeepalive, call RpcKeepaliveCall) error {
	if rpcKeepalive == nil {
		return errors.New("rpcKeepalive must not be nil")
	}
	if rpcKeepalive.ServerNode == nil {
		return errors.New("rpcKeepalive.ServerNode must not be nil")
	}
	if rpcKeepalive.Conn == nil {
		return errors.New("rpcKeepalive.Conn must not be nil")
	}
	if call == nil {
		return errors.New("call must not be nil")
	}
	if rpcKeepalive.KeepaliveTime <= 0 {
		return errors.New("rpcKeepalive.KeepaliveTime need more than 0")
	}
	rpcKeepalive.call = call
	go func (){
		rpcKeepalive.runRpcKeepalive()
	}()
	return nil
}

func (rpc *RpcKeepalive) Close() {
	rpc.cancel()
}

func (rpc *RpcKeepalive) runRpcKeepalive() {
	isConnected := false
	isReConnect := false
	ctx, cancel := context.WithCancel(context.Background())
	rpc.cancel = cancel
	for {
		select {
		case <-ctx.Done():
			rpc.call.RpcServeDisconnected(rpc, true)
			err := rpc.Conn.Close()
			if err != nil {
				rpc.call.RpcKeepaliveError("rpc close error: ", err)
			}
			return
		default:
			state := rpc.Conn.GetState()
			if state == connectivity.Ready {
				if isConnected == false {
					rpc.call.RpcServeConnected(rpc, isReConnect)
					isConnected = true
					isReConnect = true
				}
			} else {
				if isConnected == true {
					rpc.call.RpcServeDisconnected(rpc, false)
					isConnected = false
				}
			}
			time.Sleep(rpc.KeepaliveTime)
		}
	}
}
