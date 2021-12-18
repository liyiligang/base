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
	"google.golang.org/grpc"
	"google.golang.org/grpc/connectivity"
	"time"
)

type RpcKeepaliveCall interface {
	RpcServeConnected(client *RpcKeepalive, isReConnect bool)
	RpcServeDisconnected(client *RpcKeepalive, isCloseByUser bool)
}

type RpcKeepalive struct {
	Data 				 interface{}
	Conn                 *grpc.ClientConn
	KeepaliveTime 	     time.Duration
	call 			     RpcKeepaliveCall
	cancel  			 context.CancelFunc
}

func RegisterRpcKeepalive(rpcKeepalive *RpcKeepalive, call RpcKeepaliveCall) error {
	if rpcKeepalive == nil {
		return errors.New("rpcKeepalive must not be nil")
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
