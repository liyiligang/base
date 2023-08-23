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
	"io"
	"net"
	"net/http"
	"strconv"
	"strings"
)

// 获取内网IP地址
func GetPrivateIP() string {
	conn, _ := net.Dial("udp", "8.8.8.8:80")
	defer conn.Close()
	localAddr := conn.LocalAddr().String()
	idx := strings.LastIndex(localAddr, ":")
	return localAddr[0:idx]
}

// 获取外网IP地址
func GetPublicIP() string {
	resp, err := http.Get("http://ip.dhcp.cn/?ip")
	if err != nil {
		return ""
	}
	defer resp.Body.Close()
	content, _ := io.ReadAll(resp.Body)
	return string(content)
}

func GetIPFromAddr(addr string) string {
	ip := addr
	comma := strings.LastIndex(addr, ":")
	if comma >= 0 {
		ip = string([]rune(addr)[:comma])
	}
	return ip
}

func GetPortFromAddr(addr string) string {
	port := ""
	comma := strings.LastIndex(addr, ":")
	if comma >= 0 {
		port = string([]rune(addr)[(comma + 1):])
	}
	return port
}

func GetPortIntFromAddr(addr string) int {
	port, _ := strconv.Atoi(GetPortFromAddr(addr))
	return port
}
