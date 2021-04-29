// Copyright 2017 The Authors. All rights reserved.
// Author: liyiligang
// Date: 2019/4/1 17:10
// Description: 主流程日志

package Jlog

import (
	"go.uber.org/zap"
	"gopkg.in/natefinch/lumberjack.v2"
)

var processSugared *zap.SugaredLogger

//初始化流程日志
func initProcessLog(config LogInitConfig) error {

	//配置编码格式
	encoder := fileEncoderConfig()

	//配置输出文件
	ioWriter := &lumberjack.Logger{
		Filename:   config.LocalPath,  // 日志文件路径
		MaxSize:    config.MaxSize,    // 每个日志文件保存的最大尺寸 单位：M
		MaxBackups: config.MaxBackups, // 最多保存多少个日志文件
		MaxAge:     config.MaxAge,     // 日志文件最多保存多少天
		Compress:   false,             // 是否压缩
		LocalTime:  true,              // 是否使用本地时间
	}

	//其他日志配置
	outConfig := fileOutputConfig()
	if config.Debug {
		outConfig.Level = zap.NewAtomicLevelAt(zap.DebugLevel)
		outConfig.Development = true
	}
	outConfig.InitialFields = config.InitialFields

	//初始化服务核心
	var err error
	processSugared, err = initLogCore(coreConfig{
		output:    ioWriter,
		encoder:   encoder,
		outConfig: outConfig,
	})
	return err
}
