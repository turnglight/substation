package routers

import (
	"github.com/gin-gonic/gin"
	"substation.com/api"
	"substation.com/setting"
)

func InitRouter() *gin.Engine {
	r := gin.New()

	r.Use(gin.Logger())

	r.Use(gin.Recovery())

	gin.SetMode(setting.RunMode)

	routerGroup := r.Group("/api")
	{
		routerGroup.GET("/monitor/device", api.GetDeviceList)
		routerGroup.GET("/monitor/data/line", api.GetDataForLine)
	}

	return r
}

