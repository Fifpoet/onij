package main

import (
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"onij/handler"
	"onij/logic"
)

func main() {
	logic.Init()

	router := SetupRouter()

	_ = router.Run(":8080")
}

func SetupRouter() *gin.Engine {
	router := gin.Default()

	// 注册 API 文档路由
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// relay
	relayGroup := router.Group("/relay")
	{
		relayGroup.GET("/type/:type", handler.GetRelayHandler)              // 根据 relayType 获取 relay 列表
		relayGroup.POST("/pin/:id", handler.PinRelayHandler)                // 固定某个 relay
		relayGroup.DELETE("/:id", handler.DelRelayByIdHandler)              // 根据 id 删除 relay
		relayGroup.POST("/upsert", handler.UpsertRelayHandler)              // 插入或更新 relay
		relayGroup.GET("/password/:pwd", handler.GetRelayByPasswordHandler) // 根据 password 获取 relay
	}

	// music
	musicGroup := router.Group("/music")
	{
		musicGroup.POST("/upsert", handler.UpsertMusicHandler) // 插入或更新 music
		musicGroup.GET("/get", handler.GetMusicHandler)
		musicGroup.GET("/list", handler.ListMusicHandler)
	}

	// meta
	metaGroup := router.Group("/meta")
	{
		metaGroup.GET("/get", handler.GetMetaHandler)
	}

	return router
}
