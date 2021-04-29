// Copyright 2019 The Authors. All rights reserved.
// Author: liyiligang
// Date: 2019/4/15 18:19
// Description: 令牌主服务

package Jtoken

import (
	"errors"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"io/ioutil"
	"time"
)

var key string

type TokenConfig struct {
	Key           string
	UserID        int64
	StartDuration time.Duration
	StopDuration  time.Duration
	UserPara      *map[string]interface{}
}

type TokenAns struct {
	Token    string
	UserID   int64
	UserPara map[string]interface{}
}

func GetSecretByPath(keyPath string) (string, error) {
	if key == "" {
		bytes, err := ioutil.ReadFile(keyPath)
		if err != nil {
			return "", err
		} else {
			key = string(bytes)
		}
	}
	return key, nil
}

func GetToken(config TokenConfig) string {

	//iss: 签发者
	//sub: 面向的用户
	//aud: 接收方
	//exp: 过期时间
	//nbf: 生效时间
	//iat: 签发时间
	//jti: 唯一身份标识

	Claims := jwt.MapClaims{
		"jti":           config.UserID,
		"nbf":           time.Now().Add(config.StartDuration).Unix(),
		"exp":           time.Now().Add(config.StartDuration+config.StopDuration).Unix(),
		"iat":           time.Now().Unix(),
		"startDuration": config.StartDuration,
		"stopDuration":  config.StopDuration,
	}

	if config.UserPara != nil {
		for k, v := range *config.UserPara {
			Claims[k] = v
		}
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, Claims)
	tokenString, _ := token.SignedString([]byte(config.Key))
	return tokenString
}

func ParseToken(tokenString string, key string) (int64, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("%v", token.Header["alg"])
		}
		return []byte(key), nil
	})
	if err != nil {
		return 0, err
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return 0, errors.New("token解析失败")
	}
	return int64(claims["jti"].(float64)), nil
}

