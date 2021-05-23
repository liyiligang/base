// Copyright 2019 The Authors. All rights reserved.
// Author: liyiligang
// Date: 2019/4/1 17:41
// Description: webSocket服务

package Jweb

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"net/http"
	"time"
)

// Websocket接口配置
type WebsocketFunc interface {
	WebsocketConnect(conn *WebsocketConn) (interface{}, error)
	WebsocketConnected(conn *WebsocketConn) error
	WebsocketClose(conn *WebsocketConn, code int, text string)
	WebsocketReceiver(conn *WebsocketConn, data *[]byte)
	WebsocketPong(conn *WebsocketConn, pingData string) string
	WebsocketError(text string, err error)
}

// Websocket配置参数
type WebsocketConfig struct {
	WriteWaitTime time.Duration
	ReadWaitTime  time.Duration
	PingWaitTime  time.Duration
	PongWaitTime  time.Duration
	Call          WebsocketFunc
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
	sendCtx     context.Context
	sendClose   context.CancelFunc
	config      WebsocketConfig
	wsParm      WebsocketParm
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
	c, err := upgrader.Upgrade(ginContext.Writer, ginContext.Request, nil)
	if err != nil {
		config.Call.WebsocketError("WebSocket握手失败", err)
		return
	}

	ctx, cancal := context.WithCancel(context.Background())

	conn := WebsocketConn{
		conn:      c,
		send:      make(chan []byte, 256),
		sendPre:   make(chan *websocket.PreparedMessage, 10),
		config:    config,
		sendCtx:   ctx,
		sendClose: cancal,
	}

	conn.wsParm.WsClientMsg = login
	conn.wsParm.WsClientAddr = clientIP
	conn.conn.SetCloseHandler(conn.closeHandler)
	conn.conn.SetPingHandler(conn.pingHandler)
	conn.conn.SetPongHandler(conn.pongHandler)

	closeErrConn := func(err error) {
		wErr := c.WriteControl(websocket.CloseMessage, websocket.FormatCloseMessage(4000, err.Error()), conn.GetDeadline(config.WriteWaitTime))
		if wErr != nil {
			config.Call.WebsocketError("WebSocket关闭失败", wErr)
		}
		c.Close()
	}

	id, err := conn.config.Call.WebsocketConnect(&conn)
	if err != nil {
		closeErrConn(err)
		return
	}

	go conn.readMessage()
	go conn.writeMessage()
	conn.connBindVal = id
	err = conn.config.Call.WebsocketConnected(&conn)
	if err != nil {
		closeErrConn(err)
		return
	}
}

// 连接关闭回调
func (ws *WebsocketConn) closeHandler(code int, text string) error {
	ws.sendClose()
	ws.config.Call.WebsocketClose(ws, code, text)
	return nil
}

// ping回调
func (ws *WebsocketConn) pingHandler(pingData string) error {
	data := ws.config.Call.WebsocketPong(ws, pingData)
	ws.conn.WriteControl(websocket.PongMessage, []byte(data), ws.GetDeadline(ws.config.PongWaitTime))
	ws.conn.SetReadDeadline(ws.GetDeadline(ws.config.PingWaitTime))
	return nil
}

// pong回调
func (ws *WebsocketConn) pongHandler(pongData string) error {
	return nil
}

// 主动关闭连接
func (ws *WebsocketConn) Close() error {
	ws.closeHandler(0, "服务端主动断开连接")
	return ws.conn.Close()
}

// 发送Byte数据
func (ws *WebsocketConn) SendByte(data *[]byte) {
	ws.send <- *data
}

// 发送String数据
func (ws *WebsocketConn) SendString(data string) {
	byte := []byte(data)
	ws.SendByte(&byte)
}

// 发送预处理数据
func (ws *WebsocketConn) sendPreData(preData *websocket.PreparedMessage) {
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

//设置读超时时间
func (ws *WebsocketConn) GetDeadline(t time.Duration) time.Time {
	if t == 0 {
		return time.Time{}
	}
	return time.Now().Add(t)
}

// 读Websocket消息
func (ws *WebsocketConn) readMessage() {
	//ws.conn.SetReadLimit(maxMessageSize)		//最大读取限制
	for {
		ws.conn.SetReadDeadline(ws.GetDeadline(ws.config.ReadWaitTime))
		_, message, err := ws.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway){
				ws.config.Call.WebsocketError("WebSocket读取错误", err)
				ws.Close()
			}
			return
		}
		ws.config.Call.WebsocketReceiver(ws, &message)
	}
}

// 写Websocket消息
func (ws *WebsocketConn) writeMessage() {
	for {
		select {
		case message, ok := <-ws.send:
			if !ok {
				ws.config.Call.WebsocketError("WebSocket读取send chan失败", nil)
				continue
			}
			ws.conn.SetWriteDeadline(ws.GetDeadline(ws.config.WriteWaitTime))
			err := ws.conn.WriteMessage(websocket.BinaryMessage, message)
			if err != nil {
				ws.config.Call.WebsocketError("WebSocket发送失败", err)
			}
		case preMessage, ok := <-ws.sendPre:
			if !ok {
				ws.config.Call.WebsocketError("WebSocket广播读取send chan失败", nil)
				continue
			}
			ws.conn.SetWriteDeadline(ws.GetDeadline(ws.config.WriteWaitTime))
			err := ws.conn.WritePreparedMessage(preMessage)
			if err != nil {
				ws.config.Call.WebsocketError("WebSocket_Prepared发送失败", err)
			}
		case <-ws.sendCtx.Done():
			return
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
