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
