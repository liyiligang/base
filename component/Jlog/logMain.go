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

package Jlog

import (
	"go.uber.org/zap"
	"gopkg.in/natefinch/lumberjack.v2"
	"io"
)

var globalSugared *zap.SugaredLogger

type LogInitConfig struct {
	Debug         bool
	LocalPath     string
	MaxSize       int
	MaxBackups    int
	MaxAge        int
	InitialFields map[string]interface{}
}

type logIOWrite struct {
	msg 	string
	key 	string
	logger 	*zap.SugaredLogger
}

func (w *logIOWrite) Write(p []byte) (n int, err error){
	w.logger.Infow(w.msg, w.key, string(p))
	return
}

func InitGlobalLog(config LogInitConfig) {
	globalSugared = InitLog(config)
}

func InitLog(config LogInitConfig) *zap.SugaredLogger {

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
	return initLogCore(coreConfig{
		output:    ioWriter,
		encoder:   encoder,
		outConfig: outConfig,
	})
}

func Debug(str string, keysAndValues ...interface{}) {
	globalSugared.Debugw(str, keysAndValues...)
}

func Info(str string, keysAndValues ...interface{}) {
	globalSugared.Infow(str, keysAndValues...)
}

func Warn(str string, keysAndValues ...interface{}) {
	globalSugared.Warnw(str, keysAndValues...)
}

func Error(str string, keysAndValues ...interface{}) {
	globalSugared.Errorw(str, keysAndValues...)
}

func DPanic(str string, keysAndValues ...interface{}) {
	globalSugared.DPanicw(str, keysAndValues...)
}

func Panic(str string, keysAndValues ...interface{}) {
	globalSugared.Panicw(str, keysAndValues...)
}

func Fatal(str string, keysAndValues ...interface{}) {
	globalSugared.Fatalw(str, keysAndValues...)
}

func IOWrite (key string, logger *zap.SugaredLogger) io.Writer {
	if logger == nil {
		logger = globalSugared
	}
	ioWrite := &logIOWrite{msg: "", key: key, logger: logger}
	return ioWrite
}
