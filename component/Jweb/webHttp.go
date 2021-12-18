package Jweb

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type HttpFunc interface {
	HttpReceiver([]byte) ([]byte, error, int)
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


