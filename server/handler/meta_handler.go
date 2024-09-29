package handler

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"onij/logic"
)

// GetMetaHandler 获取枚举
// @Summary 获取枚举
// @Description 获取用户详细信息
// @Tags 用户
// @Accept  json
// @Produce  json
// @Param   type     path    int     true        "用户ID"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Router /meta/get [get]
func GetMetaHandler(c *gin.Context) {
	res, err := logic.NewLocalLogic().GetMeta()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get relay"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": res})
}
