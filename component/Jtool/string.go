package Jtool

import (
	"errors"
	"strings"
)

// SubStrWithCount 按计数逆序获取子字符串
// 例如: LastSubStrWithCount("/AA/BB/CC/DD/EE", "/", 2)
// 返回: "/DD/EE"
func SubStrWithCount(str string, sep string, cnt int) (string, error) {
	strList := strings.Split(str, sep)
	if cnt > len(strList) {
		return "", errors.New("计数超出字符串数组长度")
	}
	strList = strList[len(strList)-cnt:]
	res := strings.Join(strList, "/")
	return res, nil
}
