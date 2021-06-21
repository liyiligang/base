package Jdiscovery

import (
	"context"
	"errors"
	"go.etcd.io/etcd/client/v3"
)

type DiscoveryNode struct {
	NodeKey      string
	NodeData     []byte
	NodeKeepLive int64
}

type discoveryWatchNodeCall struct {
	NodeConnect    func(nodeData []byte)
	NodeDisconnect func(nodeData []byte)
}

type DiscoveryWatchNode struct {
	NodeKey    string
	NodeData   []byte
	NodeCall   discoveryWatchNodeCall
	nodeCtx    context.Context
	nodeCancel context.CancelFunc
}

func (discovery *Discovery) RegisterNode(node *DiscoveryNode) error {
	leaseID, err := discovery.getNodeStateGrantID(node.NodeKeepLive)
	if err != nil {
		return err
	}
	return discovery.SetData(node.NodeKey, string(node.NodeData), clientv3.WithLease(leaseID))
}

func (discovery *Discovery) UnRegisterNode(nodeKey string) error {
	_, err := discovery.Client.Delete(discovery.getRequestContext(), nodeKey)
	return err
}

func (discovery *Discovery) getNodeStateGrantID(nodeKeepLive int64) (clientv3.LeaseID, error) {
	resp, err := discovery.Client.Grant(discovery.getRequestContext(), nodeKeepLive)
	if err != nil {
		return resp.ID, err
	}
	_, err = discovery.Client.KeepAlive(context.Background(), resp.ID)
	if err != nil {
		return resp.ID, err
	}
	return resp.ID, nil
}

func (discovery *Discovery) RegisterNodeWatch(watchNode *DiscoveryWatchNode) error {
	watchNode.nodeCtx, watchNode.nodeCancel = context.WithCancel(context.Background())
	err := discovery.getNodeWatch(watchNode)
	if err != nil {
		return err
	}
	go discovery.startNodeWatch(watchNode)
	return nil
}

func (discovery *Discovery) UnRegisterNodeWatch(watchKey string) error {
	watch, ok := discovery.DiscoveryWatchNodeMap[watchKey]
	if !ok {
		return errors.New("找不到" + watchKey + "对应配置")
	}
	watch.nodeCancel()
	delete(discovery.DiscoveryWatchNodeMap, watchKey)
	return nil
}

func (discovery *Discovery) getNodeWatch(watchNode *DiscoveryWatchNode) error {
	val, err := discovery.GetData(watchNode.NodeKey)
	if err != nil {
		return err
	}
	watchNode.NodeCall.NodeConnect(val)
	return nil
}

func (discovery *Discovery) startNodeWatch(watchNode *DiscoveryWatchNode) {
	discovery.WatchData(watchNode.nodeCtx, watchNode.NodeKey, func(ev *clientv3.Event) {
		switch ev.Type {
		case clientv3.EventTypePut:
			watchNode.NodeCall.NodeConnect(ev.Kv.Value)
		case clientv3.EventTypeDelete:
			watchNode.NodeCall.NodeDisconnect(ev.PrevKv.Value)
		}
	})
}