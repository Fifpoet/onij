package handler

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"onij/handler/resq"
	"onij/logic"
	"onij/util"
)

// UpsertMusicHandler 插入或更新music
// @Summary 插入或更新music
// @Description 获取用户详细信息
// @Tags 用户
// @Accept  json
// @Produce  json
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Router /music/upsert [post]
func UpsertMusicHandler(c *gin.Context) {
	var req resq.UpsertMusicReq
	var err error

	// 绑定普通表单字段
	if err = c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	res, err := logic.NewMusicLogic().Save(req.ToModel())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to upsert music"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": res})
}

// GetMusicHandler 获取music
// @Summary 获取music
// @Description 获取用户详细信息
// @Tags 用户
// @Accept  json
// @Produce  json
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Router /music/get [post]
func GetMusicHandler(c *gin.Context) {
	id := util.GetInt(c, "id")

	res, err := logic.NewMusicLogic().GetMusic(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to upsert music"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": res})
}

// ListMusicHandler 列出music
// @Summary 列出music
// @Description 获取用户详细信息
// @Tags 用户
// @Accept  json
// @Produce  json
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Router /music/list [post]
func ListMusicHandler(c *gin.Context) {
	name := c.Param("title")
	at := util.GetInt(c, "artist")
	tp := util.GetInt(c, "type")

	res, err := logic.NewMusicLogic().ListByCond(name, at, tp)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to upsert music"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": res})
}
