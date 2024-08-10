package main

import (
	"github.com/gin-gonic/gin"
	"onij/inject"
)

func main() {
	app := inject.InitializeApp()

	router := SetupRouter(app)

	_ = router.Run(":8080")
}

func SetupRouter(app *inject.App) *gin.Engine {

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
