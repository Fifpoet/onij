package handler

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"onij/handler/resq"
	"onij/logic"
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

	// 处理文件上传, 没有则忽略
	req.Cover, _ = c.FormFile("cover")
	req.Mp, _ = c.FormFile("mp")
	req.Lyric, _ = c.FormFile("lyric")
	req.Sheet, _ = c.FormFile("sheet")

	res, err := logic.NewMusicLogic().Save(req.ToModel())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to upsert music"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": res})
}
