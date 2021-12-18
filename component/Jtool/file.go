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
	"io/ioutil"
	"os"
)

//创建系统临时文件
func CreateSysTmpFile(fileName string, data []byte)(string, error){
	tmpFile, err := ioutil.TempFile(os.TempDir(), fileName)
	if err != nil {
		return "", err
	}
	_, err = tmpFile.Write(data)
	if err != nil {
		return "", err
	}
	return tmpFile.Name(), err
}

//判断文件是否存在
func IsFileExist(filePath string) bool{
	s, err := os.Lstat(filePath)
	return !os.IsNotExist(err) && !s.IsDir()
}

//判断文件夹是否存在
func IsDirExist(dirPath string) bool{
	s, err:=os.Stat(dirPath)
	if err != nil{
		return false
	}
	return s.IsDir()
}

//创建文件夹
func MakeDir(dirPath string) error {
	err := os.Mkdir(dirPath,755)
	if err != nil{
		return err
	}
	return nil
}

//按字节读取文件
func ReadFileWithSize(fileName string, size uint64, call func ([]byte) error) error {
	file, err := os.Open(fileName)
	if err != nil {
		return err
	}
	defer file.Close()
	return ReadIOWithSize(file, size, call)
}
