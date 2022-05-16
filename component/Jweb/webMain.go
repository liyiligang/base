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
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/liyiligang/base/component/Jtool"
	"github.com/unrolled/secure"
	"io"
	"log"
	"net/http"
)

type WebConfig struct {
	Debug          bool
	Origin         bool
	Addr           string
	PublicKeyPath  string
	PrivateKeyPath string
	AccessWrite    io.Writer
	ErrorWrite     io.Writer
	RouteCall      func(r *gin.Engine)
}

func WebInit(config WebConfig) *http.Server {
	if !config.Debug {
		gin.SetMode(gin.ReleaseMode)
	}
	gin.DisableConsoleColor()
	if config.AccessWrite != nil {
		gin.DefaultWriter = config.AccessWrite
	}
	if config.ErrorWrite != nil {
		gin.DefaultErrorWriter = config.ErrorWrite
	}

	r := gin.Default()
	if config.Origin {
		r.Use(accessOrigin)
	}
	if config.RouteCall != nil {
		config.RouteCall(r)
	}

	srv := &http.Server{
		Addr:    config.Addr,
		Handler: r,
	}
	go func() {
		if config.PublicKeyPath != "" && config.PrivateKeyPath != "" {
			config.initError(srv.ListenAndServeTLS(config.PublicKeyPath, config.PrivateKeyPath))
		}else {
			config.initError(srv.ListenAndServe())
		}
	}()
	return srv
}

//错误处理
func (config *WebConfig) initError(err error) {
	if err != nil && err != http.ErrServerClosed {
		log.Panic("web server start up failed: ", err)
	}
}

//允许跨域
func accessOrigin(c *gin.Context) {
	c.Header("Access-Control-Allow-Origin", "*")
	c.Header("Access-Control-Allow-Headers", "Content-Type,AccessToken,X-CSRF-Token, Authorization, Token")
	c.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS")
	c.Header("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, X-Requested-With, Access-Control-Allow-Headers, Content-Type")
	c.Header("Access-Control-Allow-Credentials", "true")

	if c.Request.Method == "OPTIONS" {
		c.Abort()
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
			SSLHost:     Jtool.GetIPFromAddr(c.Request.Host) + ":" + port,
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
