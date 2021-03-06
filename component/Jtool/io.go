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
