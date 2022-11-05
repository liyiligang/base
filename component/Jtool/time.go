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
	"time"
)

const TimeFormat = "2006-01-02 15:04:05"
const TimeFormatSeries = "20060102150405"
const TimeFormatUTC = "2006-01-02T15:04:05.000Z"

// GetCurTimeFormatStandard 获取标准时间格式字符串
// 2006-01-02 15:04:05
func GetCurTimeFormatStandard() string {
	currentTime := time.Now()
	return currentTime.Format(TimeFormat)
}

// GetCurTimeFormatSeries 获取标准时间连续格式字符串
// 20060102150405
func GetCurTimeFormatSeries() string {
	currentTime := time.Now()
	return currentTime.Format(TimeFormatSeries)
}

// TimeUnixToFormat 时间戳转时间字符串
// 2006-01-02 15:04:05
func TimeUnixToFormat(timeUnix int64) string {
	t := time.Unix(timeUnix, 0)
	return t.Format(TimeFormat)
}

// TimeUnixToTime 时间戳转time.Time (秒)
func TimeUnixToTime(timeUnix int64) time.Time {
	return time.Unix(timeUnix, 0)
}

// TimeMsUnixToTime 时间戳转time.Time (豪秒)
func TimeMsUnixToTime(timeUnix int64) time.Time {
	return time.UnixMilli(timeUnix)
}

// GetMillisecondTimeInt64 获取当前时间戳 (毫秒)
func GetMillisecondTimeInt64() int64 {
	return time.Now().UnixMilli()
}

// TimeStringToMsTimeInt64 时间转时间戳 (毫秒)
func TimeStringToMsTimeInt64(time time.Time) int64 {
	return time.UnixMilli()
}

// MsTimeUnixTo 时间戳转倒计时 (毫秒)
// 秒, 分钟, 小时, 天
func MsTimeUnixTo(timeUnix int64) string {
	if timeUnix < 1000 {
		return "1 秒"
	}

	timeUnix = timeUnix/1000
	if timeUnix < 60 {
		return Int64ToString(timeUnix) + " 秒"
	}

	timeUnix = timeUnix / 60
	if timeUnix < 60 {
		return Int64ToString(timeUnix) + " 分钟"
	}

	timeUnix = timeUnix / 24
	if timeUnix < 24 {
		return Int64ToString(timeUnix) + " 小时"
	}
	return Int64ToString(timeUnix) + " 天"
}
