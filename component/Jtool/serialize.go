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
	"reflect"
	"strconv"
	"strings"
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

//int to string
func IntToString(i int) string {
	return strconv.FormatInt(int64(i), 10)
}

//int32 to string
func Int32ToString(i int32) string {
	return strconv.FormatInt(int64(i), 10)
}

//int64 to string
func Int64ToString(i int64) string {
	return strconv.FormatInt(i, 10)
}

//uint64 to string
func Uint64ToString(i uint64) string {
	return strconv.FormatInt(int64(i), 10)
}

//float32 to string
func Float32ToString(f float32, prec int) string {
	return strconv.FormatFloat(float64(f), 'f', prec,32)
}

//float64 to string
func Float64ToString(f float64, prec int) string {
	return strconv.FormatFloat(f, 'f', prec,64)
}

//string to int64
func StringToInt64(s string) int64 {
	i, _ := strconv.ParseInt(s, 10, 64)
	return i
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
//%u7248%u6743%u6240%u6709 %u4E8C%u96F6%u4E8C%u4E8C %u674E%u6613%u529B%u521A
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
	for i := 0; i < valLen; i++ {
		itemList[i] = val.Index(i).Interface()
	}
	return itemList
}
