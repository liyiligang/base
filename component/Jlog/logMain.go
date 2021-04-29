// Copyright 2017 The Authors. All rights reserved.
// Author: liyiligang
// Date: 2019/4/1 17:10
// Description: 日志主服务

package Jlog


// 日志服务初始化配置
type LogInitConfig struct {
	Debug         bool
	LocalPath     string
	MaxSize       int
	MaxBackups    int
	MaxAge        int
	InitialFields map[string]interface{}
}

// 初始化日志服务
func LogInit(config LogInitConfig) error {
	//流程日志
	return initProcessLog(config)
}

//Debug 输出Debug级别日志
func Debug(str string, keysAndValues ...interface{}) {
	processSugared.Debugw(str, keysAndValues...)
}

//Info 输出Info级别日志
func Info(str string, keysAndValues ...interface{}) {
	processSugared.Infow(str, keysAndValues...)
}

//Warn 输出Warn级别日志
func Warn(str string, keysAndValues ...interface{}) {
	processSugared.Warnw(str, keysAndValues...)
}

//Error 输出Error级别日志
func Error(str string, keysAndValues ...interface{}) {
	processSugared.Errorw(str, keysAndValues...)
}

//DPanic 输出DPanic级别日志
func DPanic(str string, keysAndValues ...interface{}) {
	processSugared.DPanicw(str, keysAndValues...)
}

//Panic 输出Panic级别日志
func Panic(str string, keysAndValues ...interface{}) {
	processSugared.Panicw(str, keysAndValues...)
}

//Fatal 输出Fatal级别日志
func Fatal(str string, keysAndValues ...interface{}) {
	processSugared.Fatalw(str, keysAndValues...)
}

