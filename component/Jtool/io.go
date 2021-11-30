// Copyright 2021 The Authors. All rights reserved.
// Author: liyiligang
// Date: 2021/11/30 10:34
// Description:

package Jtool

import (
	"errors"
	"io"
)

func ReadIOWithSize(reader io.Reader, size uint64, call func ([]byte) error) error {
	if reader == nil || size == 0 || call == nil {
		return errors.New("parameter is error")
	}
	buf := make([]byte, size)
	for {
		n, err := reader.Read(buf)
		if err != nil {
			if errors.Is(err, io.EOF){
				break
			}
			return err
		}
		err = call(buf[:n])
		if err != nil {
			return err
		}
	}
	return nil
}
