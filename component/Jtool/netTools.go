// Copyright 2019 The Authors. All rights reserved.
// Author: liyiligang
// Date: 2019/4/20 20:47
// Description: 网络工具包

package Jtool

import (
	"io/ioutil"
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
	resp, err := http.Get("http://myexternalip.com/raw")
	if err != nil {
		return ""
	}
	defer resp.Body.Close()
	content, _ := ioutil.ReadAll(resp.Body)
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
