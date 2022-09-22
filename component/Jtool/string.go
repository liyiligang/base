package Jtool

import (
	"errors"
	"github.com/mattn/go-runewidth"
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

// GetStrWidth 获取字符串占用的宽度
// 例如: GetStrWidth("abcde"); 返回 5
// 例如: GetStrWidth("啊啊啊啊啊"); 返回 10
// 例如: GetStrWidth("123啊啊"); 返回 7
// 例如: GetStrWidth("abc哈哈123"); 返回 10
func GetStrWidth(str string) int {
	return runewidth.StringWidth(str)
}