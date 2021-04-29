// Copyright 2019 The Authors. All rights reserved.
// Author: liyiligang
// Date: 2020/1/21 15:24
// Description:

package Jconfig

import (
	"bytes"
	"github.com/spf13/viper"
	"log"
)

func ReadConfigFromPath(config interface{}, path string) {
	viperPath := viper.New()
	viperPath.SetConfigName("config")
	viperPath.SetConfigType("toml")
	if path == "" {
		viperPath.AddConfigPath(".")
	} else {
		viperPath.AddConfigPath(path)
	}

	if err := viperPath.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			log.Fatal("找不到本地配置文件" + path + "/Jconfig.toml")
		} else {
			log.Fatal("读取本地配置文件错误", "err: ", err)
		}
	}

	if err := viperPath.Unmarshal(config); err != nil {
		log.Fatal("本地配置文件解析失败", "err: ", err)
	}
	log.Println("本地配置文件读取成功")
}

func ReadConfigFromByte(config interface{}, data []byte) error {
	viperByte := viper.New()
	viperByte.SetConfigType("toml")
	if err := viperByte.ReadConfig(bytes.NewBuffer(data)); err != nil {
		if err != nil {
			return err
		}
	}
	if err := viperByte.Unmarshal(config); err != nil {
		return err
	}
	return nil
}
