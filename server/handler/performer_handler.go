package handler

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"onij/handler/resq"
	"onij/logic"
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
