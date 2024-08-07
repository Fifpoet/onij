package main

import (
	"github.com/gin-gonic/gin"
	"onij/inject"
)

func main() {
	router := SetupRouter()

	_ = router.Run(":8080")
}

func SetupRouter() *gin.Engine {
	inject.InitializeApp()

	router := gin.Default()

	// tod路由
	todGroup := router.Group("/tod")
	{
		todGroup.GET("/", nil)
		todGroup.GET("/:id", nil)
	}

	return router
}
