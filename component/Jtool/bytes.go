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
	"archive/zip"
	"bytes"
	"io"
)

func ReadByteWithSize(data []byte, size uint64, call func ([]byte) error) error {
	if data == nil || len(data) == 0 {
		return nil
	}
	buffer := bytes.NewBuffer(data)
	return ReadIOWithSize(buffer, size, call)
}

func CompressByteWithZip(dataName string, data []byte) ([]byte, error){
	var err error
	var zipBuffer = new(bytes.Buffer)
	var zipWriter = zip.NewWriter(zipBuffer)
	var zipEntry io.Writer
	zipEntry, err = zipWriter.Create(dataName)
	if err != nil {
		return nil, err
	}
	_, err = zipEntry.Write(data)
	if err != nil {
		return nil, err
	}
	err = zipWriter.Close()
	if err != nil {
		return nil, err
	}
	return zipBuffer.Bytes(), nil
}