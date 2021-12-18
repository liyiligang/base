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

package Jtool

import (
	"math/rand"
	"sync"
	"time"
)

type ConnManage struct {
	serverConnMap  sync.Map
	serverConnList []interface{}
	connListLock   sync.RWMutex
	pollingCnt     int
}

func (connManage *ConnManage) Load(id interface{}) (interface{}, bool) {
	return connManage.serverConnMap.Load(id)
}

func (connManage *ConnManage) GetLen() int {
	len := 0
	connManage.serverConnMap.Range(func(key, value interface{}) bool {
		len++
		return true
	})
	return len
}

func (connManage *ConnManage) AddConnList(id interface{}) {
	connManage.connListLock.Lock()
	hasVal := false
	for _, val := range connManage.serverConnList {
		if val == id {
			hasVal = true
			break
		}
	}
	if !hasVal {
		connManage.serverConnList = append(connManage.serverConnList, id)
	}
	connManage.connListLock.Unlock()
}

func (connManage *ConnManage) DelConnList(id interface{}) {
	connManage.connListLock.Lock()
	for index, val := range connManage.serverConnList {
		if val == id {
			connManage.serverConnList = append(connManage.serverConnList[:index], connManage.serverConnList[index+1:]...)
			//break
		}
	}
	connManage.connListLock.Unlock()
}

func (connManage *ConnManage) Rand() (interface{}, bool) {
	connManage.connListLock.RLock()
	len := len(connManage.serverConnList)
	if len == 0 {
		connManage.connListLock.RUnlock()
		return nil, false
	}
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	pos := r.Intn(len)
	id := connManage.serverConnList[pos]
	connManage.connListLock.RUnlock()
	return connManage.Load(id)
}

func (connManage *ConnManage) Polling() (interface{}, bool) {
	connManage.connListLock.RLock()
	len := len(connManage.serverConnList)
	if len == 0 {
		connManage.connListLock.RUnlock()
		return nil, false
	}

	connManage.pollingCnt = (connManage.pollingCnt + 1) % len
	id := connManage.serverConnList[connManage.pollingCnt]
	connManage.connListLock.RUnlock()
	return connManage.Load(id)
}

func (connManage *ConnManage) Store(id interface{}, conn interface{}) {
	connManage.serverConnMap.Store(id, conn)
	connManage.AddConnList(id)
}

func (connManage *ConnManage) Delete(id interface{}) {
	connManage.serverConnMap.Delete(id)
	connManage.DelConnList(id)
}

func (connManage *ConnManage) IsExist(id interface{}) bool {
	_, ok := connManage.serverConnMap.Load(id)
	return ok
}

func (connManage *ConnManage) IsExistDelayCheck(id interface{}, interval time.Duration, repeatNum int) bool {
	for i:=0; i<repeatNum; i++ {
		if !connManage.IsExist(id) {
			return false
		}
		time.Sleep(interval)
	}
	return connManage.IsExist(id)
}

func (connManage *ConnManage) LoadAll(f func(key, value interface{}) bool) {
	connManage.serverConnMap.Range(f)
}
