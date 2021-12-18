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

package Jsync

import (
	"context"
	"errors"
	uuid "github.com/satori/go.uuid"
	"sync"
	"time"
)

//SyncResult 实现了在异步操作中进行阻塞, 并等待结果返回的功能
//用于tcp, websocket等网络通信中需要等待返回结果的需求, 实现类似于一个http请求的效果
type SyncResult struct {
	uuidMap 	sync.Map
}

type syncResultVal struct {
	resChan 	chan interface{}
	context		context.Context
	cancel	    context.CancelFunc
}

// NewSyncResultID 为同步等待设置一个超时时间, 同时返回一个uuid, 作为参数传递的标记
// 设置的超时时间在调用该方法后立即生效, 确保在超时前完成一次同步结果的所有流程, 否则调用其他方法将直接返回错误
func (sr *SyncResult) NewSyncResultID(timeout time.Duration) string {
	key := uuid.NewV4().String()
	val := syncResultVal{resChan: make(chan interface{}, 1)}
	val.context, val.cancel = context.WithTimeout(context.Background(), timeout)
	sr.uuidMap.Store(key, &val)
	return key
}

// GetSyncResultWithID 阻塞调用该方法的goroutine, 直到读取到值或者超时
// 在调用 NewSyncResultID 后, 需确保调用一次该方法, 否则将可能导致内存泄漏
func (sr *SyncResult) GetSyncResultWithID(uid string) (interface{}, error) {
	defer sr.uuidMap.Delete(uid)
	res, err := sr.getSyncResultValWithID(uid)
	if err != nil {
		return nil, err
	}
	if res.context.Err() != nil {
		return nil, err
	}
	select {
	case <-res.context.Done():
		return nil, res.context.Err()
	case data := <- res.resChan:
		return data, nil
	}
}

// SetSyncResultWithID 写入同步数据, 无论该数据是否被读取, 都将立即返回
// 如果写入参数时已经超时, 将返回错误
// 确保一个uuid只写入一次数据, 否则可能会阻塞调用该方法的goroutine
func (sr *SyncResult) SetSyncResultWithID(uid string, data interface{}) error {
	res, err := sr.getSyncResultValWithID(uid)
	if err != nil {
		return err
	}
	if res.context.Err() != nil {
		return res.context.Err()
	}
	select {
	case <-res.context.Done():
		return res.context.Err()
	case res.resChan <- data:
		return nil
	}
}

// CancelSyncResultWithID 立即失效该uuid对应的所有流程, 释放所有方法的阻塞行为, 这将会导致它们返回错误
func (sr *SyncResult) CancelSyncResultWithID(uid string) error {
	res, err := sr.getSyncResultValWithID(uid)
	if err != nil {
		return err
	}
	res.cancel()
	return nil
}

func (sr *SyncResult) getSyncResultValWithID(uid string) (*syncResultVal, error) {
	val, ok := sr.uuidMap.Load(uid)
	if !ok {
		return nil, errors.New("uuid " + uid + "is not found")
	}
	c, ok := val.(*syncResultVal)
	if !ok {
		return nil, errors.New("val assert fail with *syncResultVal")
	}
	return c, nil
}