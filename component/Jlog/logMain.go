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
	"strings"
)

var logHandle *zap.SugaredLogger

type LogConfig struct {
	Debug         bool
	Path          string
	Level         string
	MaxSize       int
	MaxBackups    int
	MaxAge        int
	Json          bool
	InitialFields map[string]interface{}
}

type logIOWrite struct{
	level       string
}
func (w *logIOWrite) Write(p []byte) (n int, err error) {
	switch strings.ToLower(w.level) {
	case "error":
		Error(string(p))
	case "warn":
		Warn(string(p))
	case "info":
		Info(string(p))
	default:
		Debug(string(p))
	}
	return
}

func InitLog(config LogConfig) {

	//配置编码格式
	encoder := fileEncoderConfig()

	//配置输出文件
	ioWriter := &lumberjack.Logger{
		Filename:   config.Path,       // 日志文件路径
		MaxSize:    config.MaxSize,    // 每个日志文件保存的最大尺寸 单位：M
		MaxBackups: config.MaxBackups, // 最多保存多少个日志文件
		MaxAge:     config.MaxAge,     // 日志文件最多保存多少天
		Compress:   false,             // 是否压缩
		LocalTime:  true,              // 是否使用本地时间
	}

	//其他日志配置
	outConfig := fileOutputConfig()
	outConfig.Level = getLevel(config.Level)
	if config.Debug {
		outConfig.Development = true
	}
	outConfig.InitialFields = config.InitialFields

	//初始化服务核心
	logHandle = initLogCore(coreConfig{
		output:    ioWriter,
		encoder:   encoder,
		outConfig: outConfig,
		jsonEncoder: config.Json,
	})
}

func Debug(keysAndValues ...interface{}) {
	logHandle.Debug(keysAndValues...)
}

func Info(keysAndValues ...interface{}) {
	logHandle.Info(keysAndValues...)
}

func Warn(keysAndValues ...interface{}) {
	logHandle.Warn(keysAndValues...)
}

func Error(keysAndValues ...interface{}) {
	logHandle.Error(keysAndValues...)
}

func DPanic(keysAndValues ...interface{}) {
	logHandle.DPanic(keysAndValues...)
}

func Panic(keysAndValues ...interface{}) {
	logHandle.Panic(keysAndValues...)
}

func Fatal(keysAndValues ...interface{}) {
	logHandle.Fatal(keysAndValues...)
}

func DebugW(str string, keysAndValues ...interface{}) {
	logHandle.Debugw(str, keysAndValues...)
}

func InfoW(str string, keysAndValues ...interface{}) {
	logHandle.Infow(str, keysAndValues...)
}

func WarnW(str string, keysAndValues ...interface{}) {
	logHandle.Warnw(str, keysAndValues...)
}

func ErrorW(str string, keysAndValues ...interface{}) {
	logHandle.Errorw(str, keysAndValues...)
}

func DPanicW(str string, keysAndValues ...interface{}) {
	logHandle.DPanicw(str, keysAndValues...)
}

func PanicW(str string, keysAndValues ...interface{}) {
	logHandle.Panicw(str, keysAndValues...)
}

func FatalW(str string, keysAndValues ...interface{}) {
	logHandle.Fatalw(str, keysAndValues...)
}

func GetLogger(level string) io.Writer {
	return &logIOWrite{level: level}
}

func getLevel(level string) zap.AtomicLevel {
	switch strings.ToLower(level) {
	case "fatal":
		return zap.NewAtomicLevelAt(zap.FatalLevel)
	case "panic":
		return zap.NewAtomicLevelAt(zap.PanicLevel)
	case "dpanic":
		return zap.NewAtomicLevelAt(zap.DPanicLevel)
	case "error":
		return zap.NewAtomicLevelAt(zap.ErrorLevel)
	case "warn":
		return zap.NewAtomicLevelAt(zap.WarnLevel)
	case "info":
		return zap.NewAtomicLevelAt(zap.InfoLevel)
	default:
		return zap.NewAtomicLevelAt(zap.DebugLevel)
	}
}
