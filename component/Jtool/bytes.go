// Copyright 2021 The Authors. All rights reserved.
// Author: liyiligang
// Date: 2021/11/29 22:21
// Description: 字节切片处理

package Jtool

import (
	"bytes"
)

func ReadByteWithSize(data []byte, size uint64, call func ([]byte) error) error {
	if data == nil || len(data) == 0 {
		return nil
	}
	buffer := bytes.NewBuffer(data)
	return ReadIOWithSize(buffer, size, call)
}