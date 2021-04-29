package Jdiscovery

import (
	"context"
	"errors"
	"github.com/coreos/etcd/clientv3"
	"github.com/coreos/etcd/mvcc/mvccpb"
)

type discoveryConfigCall func(oldConfig []byte, newConfig []byte)
type DiscoveryConfig struct {
	ConfigKey    string
	ConfigCall   discoveryConfigCall
	configCtx    context.Context
	configCancel context.CancelFunc
}

func (discovery *Discovery) RegisterConfigWatch(config *DiscoveryConfig) error {
	if config.ConfigKey == "" {
		return errors.New("配置名不能为空")
	}
	if config.ConfigCall == nil {
		return errors.New("配置回调函数不能为空")
	}
	config.configCtx, config.configCancel = context.WithCancel(context.Background())
	data, err := discovery.GetConfig(config.ConfigKey)
	if err != nil {
		return err
	}
	config.ConfigCall(nil, data)
	go discovery.startConfigWatch(config)
	discovery.DiscoveryWatchConfigMap[config.ConfigKey] = config
	return nil
}

func (discovery *Discovery) UnRegisterConfigWatch(configKey string) error {
	config, ok := discovery.DiscoveryWatchConfigMap[configKey]
	if !ok {
		return errors.New("找不到" + configKey + "对应配置")
	}
	config.configCancel()
	delete(discovery.DiscoveryWatchConfigMap, configKey)
	return nil
}

func (discovery *Discovery) GetConfig(configKey string) ([]byte, error) {
	return discovery.GetData(configKey)
}

func (discovery *Discovery) SetConfig(configKey string, data string) error {
	return discovery.SetData(configKey, data)
}

func (discovery *Discovery) startConfigWatch(config *DiscoveryConfig) {
	discovery.WatchData(config.configCtx, config.ConfigKey, func(ev *clientv3.Event) {
		switch ev.Type {
		case mvccpb.PUT:
			var preData []byte
			if ev.PrevKv != nil {
				preData = ev.PrevKv.Value
			}
			config.ConfigCall(preData, ev.Kv.Value)
		case mvccpb.DELETE:
			config.ConfigCall(ev.PrevKv.Value, nil)
		}
	})
}
