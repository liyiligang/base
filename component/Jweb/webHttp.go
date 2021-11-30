package Jweb

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type HttpFunc interface {
	HttpReceiver([]byte) ([]byte, error, int)
	HttpUploadFile(c *gin.Context) ([]byte, error)
	HttpDownloadFile(c *gin.Context) (string, error)
}

type HttpConfig struct {
	Call          HttpFunc
}

func (config *HttpConfig) HttpHandle(c *gin.Context) {
	raw, err := c.GetRawData()
	if err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return
	}
	ans, err, errCode:= config.Call.HttpReceiver(raw)
	if errCode != 0 {
		if err != nil{
			c.String(errCode, err.Error())
		}else {
			c.String(errCode, "")
		}
		return
	}
	c.String(http.StatusOK, string(ans))
}

func (config *HttpConfig) HttpUploadFile(c *gin.Context) {
	ans, err := config.Call.HttpUploadFile(c)
	if err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return
	}
	c.String(http.StatusOK, string(ans))
}

func (config *HttpConfig) HttpDownloadFile(c *gin.Context) {
	ans, err := config.Call.HttpDownloadFile(c)
	if err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return
	}
	c.File(ans)
}
