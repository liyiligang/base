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

//获取标准时间格式字符串
func GetCurTimeFormatStandard() string {
	currentTime := time.Now()
	return currentTime.Format(TimeFormat)
}

//获取标准时间连续格式字符串
func GetCurTimeFormatSeries() string {
	currentTime := time.Now()
	return currentTime.Format(TimeFormatSeries)
}

//时间戳转时间字符串
func TimeUnixToFormat(timeUnix int64) string {
	t := time.Unix(timeUnix, 0)
	return t.Format(TimeFormat)
}

//时间戳转time.Time
func TimeUnixToTime(timeUnix int64) time.Time {
	return time.Unix(timeUnix, 0)
}

//get time int64 (微秒)
func GetMicrosecondTimeInt64() int64 {
	return time.Now().UnixNano() / 1e6
}

//time.Time to int64 (微秒)
func TimeStringToMsTimeInt64(time time.Time) int64 {
	return time.UnixNano() / 1e6
}
