package handler

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"onij/handler/resq"
	"onij/logic"
	"onij/util"
)

// GetRelayHandler 根据relayType获取relay列表
func GetRelayHandler(c *gin.Context) {
	typ := util.GetInt(c, "type")
	res, err := logic.NewRelayLogic().GetRelayByType(typ)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get relay"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": res})
}

// PinRelayHandler 固定某个relay
func PinRelayHandler(c *gin.Context) {
	id := util.GetInt(c, "id")
	res, err := logic.NewRelayLogic().PinRelay(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to pin relay"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": res})
}

// DelRelayByTypeHandler 根据relayType删除relay
func DelRelayByTypeHandler(c *gin.Context) {
	typ := util.GetInt(c, "type")
	err := logic.NewRelayLogic().DelByType(typ)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get relay"})
		return
	}
	c.JSON(http.StatusOK, gin.H{})
}

// DelRelayByIdHandler 根据id删除relay
func DelRelayByIdHandler(c *gin.Context) {
	id := util.GetInt(c, "id")
	res, err := logic.NewRelayLogic().DelById(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get relay"})
	}
	c.JSON(http.StatusOK, gin.H{"data": res})
}

// UpsertRelayHandler 插入或更新relay
func UpsertRelayHandler(c *gin.Context) {
	req := new(resq.UpsertRelayReq)
	err := c.BindJSON(req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}
	res, err := logic.NewRelayLogic().Save(req.ToModel())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to upsert relay"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": res})
}
