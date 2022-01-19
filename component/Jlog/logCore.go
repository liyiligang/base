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
	"io"
	"os"
	"runtime"
	"sort"
	"strconv"
	"strings"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

//核心配置
type coreConfig struct {
	output    io.Writer
	outConfig zap.Config
	encoder   zapcore.EncoderConfig
}

//创建日志核心服务
func initLogCore(config coreConfig) *zap.SugaredLogger {
	//初始化日志服务
	var cores []zapcore.Core

	//定义输出级别
	priority := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl >= config.outConfig.Level.Level()
	})

	//本地日志
	if config.output != nil {
		localLog := zapcore.AddSync(config.output)
		localLogEncoder := zapcore.NewJSONEncoder(config.encoder)
		cores = append(cores, zapcore.NewCore(localLogEncoder, localLog, priority))
	}

	//控制台日志
	if config.outConfig.Development {
		console := zapcore.Lock(os.Stdout)
		consoleEncoderConfig := zap.NewDevelopmentEncoderConfig()
		consoleEncoderConfig.EncodeTime = consoleTimeEncoder
		consoleEncoder := zapcore.NewConsoleEncoder(consoleEncoderConfig)
		cores = append(cores, zapcore.NewCore(consoleEncoder, console, priority))
	}

	//创建主日志服务
	core := zapcore.NewTee(cores...)
	logger := zap.New(core, buildOptions(config.outConfig)...)

	//创建日志服务接口指针
	return logger.Sugar()
}

//设置日志服务参数
func buildOptions(logConfig zap.Config) []zap.Option {

	var opts []zap.Option

	//设置debug模式
	if logConfig.Development {
		opts = append(opts, zap.Development())
	}

	//输出定位信息
	if !logConfig.DisableCaller {
		opts = append(opts, zap.AddCaller())
		opts = append(opts, zap.AddCallerSkip(1))
	}

	//输出堆栈信息
	stackLevel := zap.PanicLevel
	if logConfig.Development {
		stackLevel = zap.ErrorLevel
	}
	if !logConfig.DisableStacktrace {
		opts = append(opts, zap.AddStacktrace(stackLevel))
	}

	//日志丢弃模式
	if logConfig.Sampling != nil {
		opts = append(opts, zap.WrapCore(func(core zapcore.Core) zapcore.Core {
			return zapcore.NewSamplerWithOptions(core, time.Second, logConfig.Sampling.Initial, logConfig.Sampling.Thereafter)
		}))
	}

	//输出追加信息
	if len(logConfig.InitialFields) > 0 {
		fs := make([]zapcore.Field, 0, len(logConfig.InitialFields))
		keys := make([]string, 0, len(logConfig.InitialFields))
		for k := range logConfig.InitialFields {
			keys = append(keys, k)
		}
		sort.Strings(keys)
		for _, k := range keys {
			fs = append(fs, zap.Any(k, logConfig.InitialFields[k]))
		}
		opts = append(opts, zap.Fields(fs...))
	}

	return opts
}

//文件日志输出配置
func fileOutputConfig() zap.Config {
	return zap.Config{
		Level:             zap.NewAtomicLevelAt(zap.InfoLevel),
		Development:       false,
		DisableCaller:     false,
		DisableStacktrace: false,
		Sampling:          nil,
		// Sampling: &zap.SamplingConfig{ //负载较大时丢弃非关键日志
		// 	Initial:    100,
		// 	Thereafter: 100,
		// },
		//InitialFields: map[string]interface{}{"@server": "ServerName"}, //默认字段
		// OutputPaths:      []string{"stderr"},
		// ErrorOutputPaths: []string{"stderr"},
	}
}

//文件日志编码配置
func fileEncoderConfig() zapcore.EncoderConfig {
	return zapcore.EncoderConfig{
		TimeKey:        "@timestamp",
		LevelKey:       "@level",
		NameKey:        "logger",
		CallerKey:      "@location",
		MessageKey:     "@message",
		StacktraceKey:  "@stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.CapitalLevelEncoder,    //原生字段名大小写
		EncodeTime:     fileTimeEncoder,                //时间参数格式
		EncodeDuration: zapcore.SecondsDurationEncoder, //时间精确度(秒, 纳秒)
		EncodeCaller:   fileCallerEncoder,              //定位格式(包名, 行数, 函数名)
	}
}

//文件日志定位信息格式(文件名:行:函数名)
func fileCallerEncoder(caller zapcore.EntryCaller, enc zapcore.PrimitiveArrayEncoder) {

	var callerStr string
	if !caller.Defined {
		callerStr = "caller is not defined"
	}

	fileN := strings.Split(caller.File, "/")
	funcN := strings.Split(runtime.FuncForPC(caller.PC).Name(), "/")
	callerStr = funcN[len(funcN)-1] + "/" + fileN[len(fileN)-1] + "(" + strconv.Itoa(caller.Line) + ")"
	enc.AppendString(callerStr)
}

//控制台时间格式
func consoleTimeEncoder(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString(t.Format("15:04:05.000"))
}

//日志文件时间格式
func fileTimeEncoder(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString(t.Format("20060102-15:04:05.000"))
}
