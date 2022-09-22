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
	"crypto/md5"
	"encoding/hex"
	"io"
	"io/ioutil"
	"mime/multipart"
	"os"
	"strings"
)

// CreateSysTmpFile 创建系统临时文件
func CreateSysTmpFile(fileName string, data []byte) (string, error) {
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

// IsFileExist 判断文件是否存在
func IsFileExist(filePath string) bool {
	s, err := os.Lstat(filePath)
	return !os.IsNotExist(err) && !s.IsDir()
}

// IsDirExist 判断文件夹是否存在
func IsDirExist(dirPath string) (bool, error) {
	s, err := os.Stat(dirPath)
	if err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}
		return false, err
	}
	return s.IsDir(), nil
}

// MakeDir 创建文件夹
func MakeDir(dirPath string) error {
	err := os.MkdirAll(dirPath, 755)
	if err != nil {
		return err
	}
	return nil
}

// MakeDirIfNoExist 判断文件夹是否存在, 不存在则创建
func MakeDirIfNoExist(dirPath string) (bool, error) {
	isExist, err := IsDirExist(dirPath)
	if err != nil {
		return false, err
	}
	if isExist {
		return true, nil
	} else {
		err := MakeDir(dirPath)
		if err != nil {
			return false, err
		}
		return false, nil
	}
}

// 按文件名查找指定目录下的一个文件
func FindDirFileWithFileName(dirPath string, fileName string) (string, error) {
	fileList, err := ioutil.ReadDir(dirPath)
	if err != nil {
		return "", err
	}
	for _, file := range fileList {
		if !file.IsDir() {
			name := file.Name()
			n := strings.LastIndex(file.Name(), ".")
			if n >= 0 {
				name = name[:n]
			}
			if name ==  fileName{
				return file.Name(), nil
			}
		}
	}
	return "", nil
}

// ReadFileWithSize 按字节读取文件
func ReadFileWithSize(fileName string, size uint64, call func([]byte) error) error {
	file, err := os.Open(fileName)
	if err != nil {
		return err
	}
	defer file.Close()
	return ReadIOWithSize(file, size, call)
}

// GetFileMd5 获取文件Md5值(file)
func GetFileMd5(file multipart.File) (string, error) {
	hash := md5.New()
	_, err := io.Copy(hash, file)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(hash.Sum(nil)), nil
}

// GetFileMd5WithPath 获取文件Md5值(path)
func GetFileMd5WithPath(filePath string) (string, error) {
	body, err := ioutil.ReadFile(filePath)
	if err != nil {
		return "", err
	}
	m := md5.New()
	m.Write(body)
	return hex.EncodeToString(m.Sum(nil)), nil
}
