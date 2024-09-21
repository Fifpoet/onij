package main

import (
	"github.com/gin-gonic/gin"
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

	// 创建 relay 分组路由
	relayGroup := router.Group("/relay")
	{
		// 定义每个路由及其对应的处理函数
		relayGroup.GET("/type/:type", handler.GetRelayHandler)              // 根据 relayType 获取 relay 列表
		relayGroup.POST("/pin", handler.PinRelayHandler)                    // 固定某个 relay
		relayGroup.DELETE("/:id", handler.DelRelayByIdHandler)              // 根据 id 删除 relay
		relayGroup.POST("/upsert", handler.UpsertRelayHandler)              // 插入或更新 relay
		relayGroup.GET("/password/:pwd", handler.GetRelayByPasswordHandler) // 根据 password 获取 relay
	}

	return router
}
