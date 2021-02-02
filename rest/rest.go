package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"substation.com/routers"
	"substation.com/setting"
)

func main() {
	router := routers.InitRouter()

	gin.SetMode(gin.ReleaseMode)

	s := &http.Server{
		Addr:           fmt.Sprintf(":%d", setting.HTTPPort),
		Handler:        router,
		ReadTimeout:    setting.ReadTimeout,
		WriteTimeout:   setting.WriteTimeout,
		MaxHeaderBytes: 1 << 20,
	}

	s.ListenAndServe()
}