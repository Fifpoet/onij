package handler

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"onij/handler/resq"
	"onij/logic"
	"strconv"
)

// UpsertPerformerHandler 插入或更新performer
// @Summary 插入或更新performer
// @Description 获取用户详细信息
// @Tags 用户
// @Accept  json
// @Produce  json
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Router /performer/upsert [post]
func UpsertPerformerHandler(c *gin.Context) {
	var req resq.UpsertPerformerReq
	var err error

	// 绑定普通表单字段
	if err = c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	res, err := logic.NewPerformerLogic().Save(req.ToModel())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to upsert performer"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": res})
}

// GetPerformerHandler 获取performer
// @Summary 获取performer
// @Description 获取用户详细信息
// @Tags 用户
// @Accept  json
// @Produce  json
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Router /performer/get [get]
func GetPerformerHandler(c *gin.Context) {
	var err error
	pTypeStr := c.Query("type")
	pType, _ := strconv.Atoi(pTypeStr)
	name := c.Query("name")

	res, err := logic.NewPerformerLogic().GetByNameAndType(name, pType)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get performer"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": res})
}
