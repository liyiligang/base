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
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/grpclog"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/peer"
	"io"
	"log"
	"net"
	"time"
)

//metadata如果要传输二进制数据，key必须以bin结尾
const ConstRpcHeader = "rpc-header-bin"


type RpcBaseConfig struct {
	Addr           string
	PublicKeyPath  string
}

type RpcServerConfig struct {
	RpcBaseConfig
	PrivateKeyPath 		string
	RegisterCall   		func(*grpc.Server)
	LogWrite 	   		io.Writer
	ServerOption 		[]grpc.ServerOption
}

type RpcClientConfig struct {
	RpcBaseConfig
	CertName       	 string
	Header           []byte
	ConnectTimeOut 	 time.Duration
	ClientOption 	 []grpc.DialOption
}

type RpcContext struct {
	RpcHeader       		[]byte
	RpcStreamClientHeader	[]byte
	RpcStreamServerTrailer 	[]byte
	RpcStreamServerHeader 	[]byte
	RpcClientAddr      		string
}

func GrpcServerInit(config RpcServerConfig) (*grpc.Server, error) {
	lis, err := net.Listen("tcp", config.Addr)
	if err != nil {
		return nil, err
	}

	if config.LogWrite != nil {
		grpclog.SetLoggerV2(grpclog.NewLoggerV2(config.LogWrite, config.LogWrite, config.LogWrite))
	}

	creds, err := credentials.NewServerTLSFromFile(config.PublicKeyPath, config.PrivateKeyPath)
	if err != nil {
		return nil, err
	}

	config.ServerOption = append(config.ServerOption, grpc.Creds(creds))
	s := grpc.NewServer(config.ServerOption...)
	config.RegisterCall(s)
	go func() {
		if err := s.Serve(lis); err != nil {
			log.Panic("rpc server startup failed: ", err)
		}
	}()
	return s, nil
}

func GrpcClientInit(config RpcClientConfig) (*grpc.ClientConn, error) {
	cred, err := credentials.NewClientTLSFromFile(config.PublicKeyPath, config.CertName)
	if err != nil {
		return nil, err
	}
	ctx := context.Background()
	if config.ConnectTimeOut != 0 {
		ctx, _ = context.WithTimeout(ctx, config.ConnectTimeOut)
	}
	config.ClientOption = append(config.ClientOption, grpc.WithTransportCredentials(cred),
		grpc.WithPerRPCCredentials(&rpcHeader{header: config.Header}))
	conn, err := grpc.DialContext(ctx, config.Addr, config.ClientOption...)
	if err != nil {
		return nil, err
	}
	return conn, nil
}

//获取连接参数
func ParseRpcContext(ctx context.Context) (RpcContext, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return RpcContext{}, errors.New("参数解析错误")
	}
	addr, ok := peer.FromContext(ctx)
	if !ok {
		return RpcContext{}, errors.New("获取对端IP错误")
	}
	return RpcContext{RpcHeader: getRpcHeader(md), RpcStreamClientHeader: getRpcStreamClientHeader(md), RpcClientAddr: addr.Addr.String()}, nil
}

//获取元数据
func getRpcHeader(md metadata.MD) []byte {
	data, ok := md[ConstRpcHeader]
	if ok {
		return []byte(data[0])
	}
	return []byte("")
}

//grpc认证
type rpcHeader struct {
	header []byte
}

//写入元数据
func (grpc *rpcHeader) GetRequestMetadata(context.Context, ...string) (map[string]string, error) {
	return map[string]string{
		ConstRpcHeader: string(grpc.header),
	}, nil
}

//允许元数据传输
func (grpc *rpcHeader) RequireTransportSecurity() bool {
	return true
}

