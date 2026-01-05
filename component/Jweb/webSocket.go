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

package Jweb

import (
	"context"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"net/http"
	"sync/atomic"
	"time"
)

const WebsocketCloseByServer = 4000

// Websocket接口配置
type WebsocketCall struct {
	WebsocketConnect   func(conn *WebsocketConn) (interface{}, error)
	WebsocketConnected func(conn *WebsocketConn) error
	WebsocketClosed    func(conn *WebsocketConn, code int, text string) error
	WebsocketReceiver  func(conn *WebsocketConn, data *[]byte) error
	WebsocketPong      func(conn *WebsocketConn, pingData string) (string, error)
	WebsocketError     func(text string, err error)
}

// Websocket配置参数
type WebsocketConfig struct {
	WriteWaitTime time.Duration
	ReadWaitTime  time.Duration
	PingWaitTime  time.Duration
	PongWaitTime  time.Duration
	Call          WebsocketCall
	UserConfig    interface{}
}

type WebsocketParm struct {
	WsClientMsg  string
	WsClientAddr string
}

// Websocket连接
type WebsocketConn struct {
	conn        *websocket.Conn
	send        chan []byte
	sendPre     chan *websocket.PreparedMessage
	context     context.Context
	cancel      context.CancelFunc
	config      WebsocketConfig
	wsParm      WebsocketParm
	sendCnt     int64
	connBindVal interface{}
}

// Websocket握手配置
var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	//WriteBufferPool:&sync.Pool{},

	//允许webSocket跨域访问
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func (config *WebsocketConfig) WsHandle(c *gin.Context) {
	login := c.Query("parameter")
	wsConnect(c, *config, login, c.ClientIP())
}

// 握手
func wsConnect(ginContext *gin.Context, config WebsocketConfig, login string, clientIP string) {
	ctx, ctxCancel := context.WithCancel(context.Background())
	conn := WebsocketConn{
		send:    make(chan []byte, 256),
		sendPre: make(chan *websocket.PreparedMessage, 10),
		config:  config,
		context: ctx,
		cancel:  ctxCancel,
	}

	var err error
	conn.conn, err = upgrader.Upgrade(ginContext.Writer, ginContext.Request, nil)
	if err != nil {
		conn.closeWs("websocket init error", err)
		return
	}

	conn.wsParm.WsClientMsg = login
	conn.wsParm.WsClientAddr = clientIP
	conn.conn.SetPingHandler(conn.pingHandler)
	conn.conn.SetPongHandler(conn.pongHandler)

	if conn.config.Call.WebsocketConnect != nil {
		id, err := conn.config.Call.WebsocketConnect(&conn)
		if err != nil {
			conn.closeWs("call websocket connect fail", err)
			return
		}
		conn.connBindVal = id
	}

	go conn.readMessage()
	go conn.writeMessage()

	if conn.config.Call.WebsocketConnected != nil {
		err = conn.config.Call.WebsocketConnected(&conn)
		if err != nil {
			conn.closeWs("call websocket connected fail", err)
			return
		}
	}
}

// 连接关闭回调
func (ws *WebsocketConn) closeHandler(code int, text string) error {
	ws.cancel()
	if ws.config.Call.WebsocketClosed != nil {
		return ws.config.Call.WebsocketClosed(ws, code, text)
	}
	return nil
}

// ping回调
func (ws *WebsocketConn) pingHandler(pingData string) error {
	if ws.config.Call.WebsocketPong != nil {
		data, err := ws.config.Call.WebsocketPong(ws, pingData)
		if err != nil {
			return err
		}
		err = ws.conn.WriteControl(websocket.PongMessage, []byte(data), ws.GetDeadline(ws.config.PongWaitTime))
		if err != nil {
			return err
		}
	}
	return ws.conn.SetReadDeadline(ws.GetDeadline(ws.config.PingWaitTime))
}

// pong回调
func (ws *WebsocketConn) pongHandler(pongData string) error {
	return nil
}

// 关闭连接
func (ws *WebsocketConn) closeWs(text string, err error) {
	ws.webSocketError(text, err)
	if ws.conn != nil {
		err := ws.conn.WriteControl(websocket.CloseMessage, websocket.FormatCloseMessage(WebsocketCloseByServer, err.Error()), ws.GetDeadline(ws.config.WriteWaitTime))
		if err != nil {
			ws.webSocketError("websocket close error", err)
		}
		//err = ws.conn.Close()
		//if err != nil {
		//	ws.webSocketError("websocket close error", err)
		//}
	}
	cErr := ws.closeHandler(WebsocketCloseByServer, err.Error())
	if cErr != nil {
		ws.webSocketError("call websocket closed fail", cErr)
	}
	return
}

// 主动关闭连接
func (ws *WebsocketConn) Close(msg string, immediately bool) {
	if !immediately {
		ws.sendWait()
	}
	ws.closeWs("websocket is active close ", errors.New(msg))
}

// 发送Byte数据
func (ws *WebsocketConn) SendByte(data *[]byte) {
	atomic.AddInt64(&ws.sendCnt, 1)
	ws.send <- *data
}

// 发送String数据
func (ws *WebsocketConn) SendString(data string) {
	bytes := []byte(data)
	ws.SendByte(&bytes)
}

// 发送预处理数据
func (ws *WebsocketConn) sendPreData(preData *websocket.PreparedMessage) {
	atomic.AddInt64(&ws.sendCnt, 1)
	ws.sendPre <- preData
}

// 获取绑定值
func (ws *WebsocketConn) GetBindVal() interface{} {
	return ws.connBindVal
}

// 获取参数值
func (ws *WebsocketConn) GetParm() WebsocketParm {
	return ws.wsParm
}

// 获取配置
func (ws *WebsocketConn) GetConfig() WebsocketConfig {
	return ws.config
}

// 设置读超时时间
func (ws *WebsocketConn) GetDeadline(t time.Duration) time.Time {
	if t == 0 {
		return time.Time{}
	}
	return time.Now().Add(t)
}

// 读Websocket消息
func (ws *WebsocketConn) readMessage() {
	for {
		select {
		case <-ws.context.Done():
			return
		default:
			err := ws.conn.SetReadDeadline(ws.GetDeadline(ws.config.ReadWaitTime))
			if err != nil {
				ws.closeWs("websocket read error", err)
				return
			}
			_, message, err := ws.conn.ReadMessage()
			if err != nil {
				ws.closeWs("websocket read error", err)
				return
			}
			if ws.config.Call.WebsocketReceiver != nil {
				err := ws.config.Call.WebsocketReceiver(ws, &message)
				if err != nil {
					ws.webSocketError("call websocket receiver fail", err)
				}
			}
		}
	}
}

// 写Websocket消息
func (ws *WebsocketConn) writeMessage() {
	for {
		select {
		case <-ws.context.Done():
			return
		case message, ok := <-ws.send:
			if !ok {
				ws.closeWs("webSocket send error", errors.New("send channal is closed"))
				return
			}
			err := ws.conn.SetWriteDeadline(ws.GetDeadline(ws.config.WriteWaitTime))
			if err != nil {
				ws.closeWs("webSocket send error", err)
				return
			}
			err = ws.conn.WriteMessage(websocket.BinaryMessage, message)
			atomic.AddInt64(&ws.sendCnt, -1)
			if err != nil {
				ws.closeWs("webSocket send error", err)
				return
			}
		case preMessage, ok := <-ws.sendPre:
			if !ok {
				ws.closeWs("webSocket pre-send error", errors.New("send channal is closed"))
				return
			}
			err := ws.conn.SetWriteDeadline(ws.GetDeadline(ws.config.WriteWaitTime))
			if err != nil {
				ws.closeWs("webSocket pre-send error", err)
				return
			}
			err = ws.conn.WritePreparedMessage(preMessage)
			atomic.AddInt64(&ws.sendCnt, -1)
			if err != nil {
				ws.closeWs("webSocket pre-send error", err)
				return
			}
		}
	}
}

// 广播Byte数据
func BroadCastByte(data *[]byte, connList []*WebsocketConn) error {
	preData, err := websocket.NewPreparedMessage(websocket.BinaryMessage, *data)
	if err == nil {
		for _, conn := range connList {
			conn.sendPreData(preData)
		}
	}
	return err
}

// 广播String数据
func BroadCastString(data string, connList []*WebsocketConn) error {
	preData, err := websocket.NewPreparedMessage(websocket.BinaryMessage, []byte(data))
	if err == nil {
		for _, conn := range connList {
			conn.sendPreData(preData)
		}
	}
	return err
}

func (ws *WebsocketConn) webSocketError(text string, err error) {
	if websocket.IsCloseError(err, websocket.CloseNormalClosure) {
		return
	}
	if websocket.IsCloseError(err, websocket.CloseGoingAway) {
		return
	}
	if errors.Is(err, websocket.ErrCloseSent) {
		return
	}
	if ws.config.Call.WebsocketError != nil {
		ws.config.Call.WebsocketError(text, err)
	} else {
		fmt.Println(text+": ", err)
	}
}

func (ws *WebsocketConn) sendWait() {
	for {
		select {
		case <-ws.context.Done():
			return
		default:
			if atomic.LoadInt64(&ws.sendCnt) > 0 {
				time.Sleep(5 * time.Millisecond)
			} else {
				return
			}
		}
	}
}
