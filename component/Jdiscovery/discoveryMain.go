package Jdiscovery

import (
	"context"
	"errors"
	"go.etcd.io/etcd/client/v3"
	"sync"
	"time"
)

type DiscoveryInitConfig struct {
	EtcdAddr       string
	ConnectTimeout int
	RequestTimeout int
}

type Discovery struct {
	Client                  *clientv3.Client
	Config                  DiscoveryInitConfig
	discoveryNodeMap 		sync.Map
	DiscoveryWatchConfigMap map[string]*DiscoveryConfig
	DiscoveryWatchNodeMap   map[string]*DiscoveryWatchNode
}

func DiscoveryInit(config DiscoveryInitConfig) (*Discovery, error) {
	discovery := Discovery{Config: config}
	client, err := clientv3.New(clientv3.Config{
		Endpoints:   []string{discovery.Config.EtcdAddr},
		DialTimeout: time.Duration(discovery.Config.ConnectTimeout) * time.Second,
	})
	if err != nil {
		return nil, err
	}
	discovery.DiscoveryWatchConfigMap = make(map[string]*DiscoveryConfig)
	discovery.DiscoveryWatchNodeMap = make(map[string]*DiscoveryWatchNode)
	discovery.Client = client
	return &discovery, nil
}

func (discovery *Discovery) getRequestContext() context.Context {
	ctx := context.Background()
	if discovery.Config.RequestTimeout != 0 {
		ctx, _ = context.WithTimeout(ctx, time.Duration(discovery.Config.RequestTimeout)*time.Second)
	}
	return ctx
}

func (discovery *Discovery) SetData(key string, data string, opts ...clientv3.OpOption) error {
	_, err := discovery.Client.Put(discovery.getRequestContext(), key, data, opts...)
	if err != nil {
		return err
	}
	return nil
}

func (discovery *Discovery) GetData(key string) ([]byte, error) {
	resp, err := discovery.Client.Get(discovery.getRequestContext(), key)
	if err != nil {
		return nil, err
	}
	if resp.Kvs == nil {
		return nil, errors.New("找不到Key: " + key + " 对应数据")
	}
	return resp.Kvs[0].Value, nil
}

func (discovery *Discovery) WatchData(ctx context.Context, key string, call func(e *clientv3.Event)) {
	for {
		watch := discovery.Client.Watch(ctx, key, clientv3.WithPrevKV())
		for res := range watch {
			for _, ev := range res.Events {
				call(ev)
			}
		}
	}
}
