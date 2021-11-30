// Copyright 2019 The Authors. All rights reserved.
// Author: liyiligang
// Date: 2019/4/1 17:41
// Description: web主服务

package Jweb

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/unrolled/secure"
	"io"
	"net/http"
	"strings"
)


// web服务初始化配置
type WebInitConfig struct {
	Debug               bool
	Addr            	string
	IsHttps             bool
	PublicKeyPath  		string
	PrivateKeyPath 		string
	RedirectAddr		string					//将此地址以http重定向至https
	LogWrite 			io.Writer
	RouteCall       	func(r *gin.Engine)
	ErrorCall      		func(str string, keysAndValues ...interface{})
}

// web服务初始化
func WebInit(config WebInitConfig) *http.Server {
	if !config.Debug {
		gin.SetMode(gin.ReleaseMode)
	}

	if config.LogWrite != nil {
		gin.DefaultWriter = config.LogWrite
		gin.DefaultErrorWriter = config.LogWrite
	}

	r := gin.Default()
	r.Use(accessOrigin)

	//重定向需要放在RouteFunc前
	if config.IsHttps && config.RedirectAddr != "" {
		r.Use(redirectHttp(getPortFromAddr(config.Addr)))
	}

	if config.RouteCall != nil {
		config.RouteCall(r)
	}

	srv := &http.Server{
		Addr:    config.Addr,
		Handler: r,
	}

	go func() {
		if config.IsHttps {
			if config.IsHttps && config.RedirectAddr != "" {
				go func() {
					config.httpError(http.ListenAndServe(config.RedirectAddr, r))
				}()
			}
			config.httpError(srv.ListenAndServeTLS(config.PublicKeyPath, config.PrivateKeyPath))
		} else {
			config.httpError(srv.ListenAndServe())
		}
	}()
	return srv
}

//错误处理
func (config *WebInitConfig) httpError(err error) {
	if err != nil && err != http.ErrServerClosed {
		config.ErrorCall("Web服务初始化失败", "err", err)
	}
}

//允许跨域
func accessOrigin(c *gin.Context) {
	c.Header("Access-Control-Allow-Origin", "*")
	c.Header("Access-Control-Allow-Headers", "Content-Type,AccessToken,X-CSRF-Token, Authorization, Token")
	c.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS")
	c.Header("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, X-Requested-With, Access-Control-Allow-Headers, Content-Type")
	c.Header("Access-Control-Allow-Credentials", "true")

	if c.Request.Method == "OPTIONS"{
		c.Writer.WriteHeader(http.StatusOK)
		return
	}
	c.Next()
}

//http重定向至https
func redirectHttp(port string) gin.HandlerFunc {
	return func(c *gin.Context) {
		middleware := secure.New(secure.Options{
			SSLRedirect: true,
			SSLHost:     getIPFromAddr(c.Request.Host) + ":" + port,
		})
		err := middleware.Process(c.Writer, c.Request)
		if err != nil {
			//如果出现错误，请不要继续。
			fmt.Println(err)
			return
		}
		// 继续往下处理
		c.Next()
	}
}

//从地址获取IP
func getIPFromAddr(addr string) string {
	ip := addr
	comma := strings.LastIndex(addr, ":")
	if comma >= 0 {
		ip = string([]rune(addr)[:comma])
	}
	return ip
}

//从地址获取端口
func getPortFromAddr(addr string) string {
	port := ""
	comma := strings.LastIndex(addr, ":")
	if comma >= 0 {
		port = string([]rune(addr)[(comma + 1):])
	}
	return port
}
