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

package Jdiscovery

import (
	"context"
	"errors"
	"go.etcd.io/etcd/client/v3"
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
		return errors.New("config key is empty")
	}
	if config.ConfigCall == nil {
		return errors.New("config call is nil")
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
	watch, ok := discovery.DiscoveryWatchConfigMap[configKey]
	if !ok {
		return errors.New("watch config " + configKey + " is not found")
	}
	watch.configCancel()
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
		case clientv3.EventTypePut:
			var preData []byte
			if ev.PrevKv != nil {
				preData = ev.PrevKv.Value
			}
			config.ConfigCall(preData, ev.Kv.Value)
		case clientv3.EventTypeDelete:
			config.ConfigCall(ev.PrevKv.Value, nil)
		}
	})
}
