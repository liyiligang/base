// Copyright 2019 The Authors. All rights reserved.
// Author: liyiligang
// Date: 2019/10/30 7:38
// Description: 序列化工具包

package Jtool

import (
	"reflect"
	"strconv"
	"strings"
	"time"
	"unsafe"
)

//结构体序列化
func StructToBytes(ptr unsafe.Pointer, len uintptr) *[]byte {
	var x reflect.SliceHeader
	x.Len = int(len)
	x.Cap = int(len)
	x.Data = uintptr(ptr)
	return (*[]byte)(unsafe.Pointer(&x))
}

//结构体反序列化
func BytesToStruct(b *[]byte) unsafe.Pointer {
	return unsafe.Pointer((*reflect.SliceHeader)(unsafe.Pointer(b)).Data)
}

//int64 to string
func Int64ToString(i int64) string {
	return strconv.FormatInt(i, 10)
}

//uint64 to string
func Uint64ToString(i uint64) string {
	return strconv.FormatInt(int64(i), 10)
}

//[]int to string
func IntSliceToString(i []int, symbol string) string {
	var s []string
	for _, v := range i {
		s = append(s, strconv.FormatInt(int64(v), 10))
	}
	return strings.Join(s, symbol)
}

//string to []int
func StringToIntSlice(s string, symbol string) []int {
	var i []int
	sp := strings.Split(s, symbol)
	for _, v := range sp {
		if pi, err := strconv.ParseInt(v, 10, 0); err == nil {
			i = append(i, int(pi))
		} else {
			return nil
		}
	}
	return i
}

//structToMap
func StructToMap(st interface{}) map[string]interface{} {
	var data = make(map[string]interface{})
	stVal := reflect.ValueOf(st).Elem()
	stType := stVal.Type()
	for i := 0; i < stVal.NumField(); i++ {
		data[strings.Title(stType.Field(i).Name)] = stVal.Field(i).Interface()
	}
	return data
}

//get time int64 (微秒)
func GetMicrosecondTimeInt64() int64 {
	return time.Now().UnixNano() / 1e6
}

//time to int64 (微秒)
func TimeStringToMsTimeInt64(time time.Time) int64 {
	return time.UnixNano() / 1e6
}

//[]rune to []byte
func RunesToBytes(data []byte) []byte {
	r := []rune(string(data))
	var b []byte
	for _, v := range r {
		b = append(b, byte(v))
	}
	return b
}

//any slice to interface slice
func SliceToInterfaceSlice(data interface{}) []interface{} {
	val := reflect.ValueOf(data)
	if val.Kind() != reflect.Slice {
		return nil
	}
	valLen := val.Len()
	if valLen == 0 {
		return nil
	}
	itemList := make([]interface{}, valLen)
	for i := 0; i < valLen; i++{
		itemList[i] = val.Index(i).Interface()
	}
	return itemList
}