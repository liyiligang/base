// Copyright 2021 The Authors. All rights reserved.
// Author: liyiligang
// Date: 2021/02/18 7:38
// Description: 时间处理

package Jtool

import "time"

//获取标准时间格式字符串
func GetCurTimeFormatStandard() string {
	currentTime := time.Now()
	return currentTime.Format("2006-01-02_15:04:05")
}
