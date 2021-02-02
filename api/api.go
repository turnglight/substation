package api

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"substation.com/models"
)


type DeviceParam struct {
	MonitorId int32
	CmdType int32
}

func GetDeviceList(c *gin.Context){
	equipments := models.GetDeviceList()
	c.JSONP(http.StatusOK, gin.H{
		"code": 200,
		"data": equipments,
	})
}

func GetDataForLine(c *gin.Context){
	var deviceParam DeviceParam
	if c.ShouldBind(&deviceParam) == nil {
		data := models.GetDataForLine(deviceParam.MonitorId, deviceParam.CmdType)
		c.JSONP(http.StatusOK, gin.H{
			"code": 200,
			"data": data,
		})
	}
}

