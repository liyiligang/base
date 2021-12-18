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
	Key           	string
	ID        	  	int64
	StartDuration 	time.Duration
	StopDuration  	time.Duration
	Custom      	map[string]interface{}
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
		"jti":           config.ID,
		"nbf":           time.Now().Add(config.StartDuration).Unix(),
		"exp":           time.Now().Add(config.StartDuration+config.StopDuration).Unix(),
		"iat":           time.Now().Unix(),
		"startDuration": config.StartDuration,
		"stopDuration":  config.StopDuration,
	}

	if config.Custom != nil {
		for k, v := range config.Custom {
			Claims[k] = v
		}
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, Claims)
	tokenString, _ := token.SignedString([]byte(config.Key))
	return tokenString
}

func ParseToken(tokenString string, key string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("%v", token.Header["alg"])
		}
		return []byte(key), nil
	})
	if err != nil {
		return nil, err
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, errors.New("token.Claims assert fail with jwt.MapClaims")
	}
	if !token.Valid {
		return nil, errors.New("token "+ tokenString + " is invalid")
	}
	return claims, nil
}

