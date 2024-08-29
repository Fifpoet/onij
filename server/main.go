package main

import (
	"github.com/gin-gonic/gin"
	"onij/logic"
)

func main() {
	logic.Init()

	router := SetupRouter()

	_ = router.Run(":8080")
}

func SetupRouter() *gin.Engine {

	router := gin.Default()

	// tod
	todGroup := router.Group("/tod")
	{
		todGroup.GET("/", nil)
		todGroup.GET("/:id", nil)
	}
	// relay
	relayGroup := router.Group("/relay")
	{
		relayGroup.GET("/", nil)
	}

	return router
}
