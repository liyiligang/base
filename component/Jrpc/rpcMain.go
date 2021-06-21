// Copyright 2019 The Authors. All rights reserved.
// Author: liyiligang
// Date: 2019/4/1 17:41
// Description: rpc主服务

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
	"net"
	"time"
)

const ConstRpcMetadata = "auth-bin"

type RpcBaseConfig struct {
	Addr           string
	PublicKeyPath  string
}

type RpcServerConfig struct {
	RpcBaseConfig
	PrivateKeyPath string
	RegisterCall   func(*grpc.Server)
	LogWrite 	   io.Writer
	ErrorCall      func(str string, keysAndValues ...interface{})
}

type RpcClientConfig struct {
	RpcBaseConfig
	CertName       string
	Auth           []byte
	ConnectTimeOut time.Duration
}

type RpcParm struct {
	RpcClientMsg       		[]byte
	RpcStreamClientMsg 		[]byte
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

	s := grpc.NewServer(grpc.Creds(creds))
	config.RegisterCall(s)
	go func() {
		if err := s.Serve(lis); err != nil {
			config.ErrorCall("Rpc服务初始化失败", "err", err)
		}
	}()
	return s, nil
}

func GrpcClientInit(config RpcClientConfig) (*grpc.ClientConn, error) {
	creds, err := credentials.NewClientTLSFromFile(config.PublicKeyPath, config.CertName)
	if err != nil {
		return nil, err
	}
	ctx := context.Background()
	if config.ConnectTimeOut != 0 {
		ctx, _ = context.WithTimeout(ctx, config.ConnectTimeOut)
	}
	conn, err := grpc.DialContext(ctx, config.Addr, grpc.WithBlock(), grpc.WithTransportCredentials(creds), grpc.WithPerRPCCredentials(&grpcAuth{auth: config.Auth}))
	if err != nil {
		return nil, err
	}
	return conn, nil
}

//grpc认证
type grpcAuth struct {
	auth []byte
}

//获取连接参数
func GetConnectParm(ctx context.Context) (RpcParm, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return RpcParm{}, errors.New("参数解析错误")
	}

	addr, ok := peer.FromContext(ctx)
	if !ok {
		return RpcParm{}, errors.New("获取对端IP错误")
	}
	return RpcParm{RpcClientMsg: GetRpcMetadata(md), RpcStreamClientMsg: GetRpcStreamMetadata(md), RpcClientAddr: addr.Addr.String()}, nil
}

//获取元数据
func GetRpcMetadata(md metadata.MD) []byte {
	data, ok := md[ConstRpcMetadata]
	if ok {
		return []byte(data[0])
	}
	return []byte("")
}

//写入元数据
func (grpc *grpcAuth) GetRequestMetadata(context.Context, ...string) (map[string]string, error) {
	return map[string]string{
		ConstRpcMetadata: string(grpc.auth),
	}, nil
}

//允许元数据传输
func (grpc *grpcAuth) RequireTransportSecurity() bool {
	return true
}

