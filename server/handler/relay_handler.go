package handler

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"onij/handler/resq"
	"onij/logic"
	"onij/util"
)

// GetRelayHandler 根据relayType获取relay列表
// @Summary 根据relayType获取relay列表
// @Description 获取用户详细信息
// @Tags 用户
// @Accept  json
// @Produce  json
// @Param   type     path    int     true        "用户ID"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Router /type/{type} [get]
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
// @Summary 固定某个relay
// @Description 获取用户详细信息
// @Tags 用户
// @Accept  json
// @Produce  json
// @Param   id     path    int     true        "用户ID"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Router /pin/{id} [get]
func PinRelayHandler(c *gin.Context) {
	id := util.GetInt(c, "id")
	res, err := logic.NewRelayLogic().PinRelay(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to pin relay"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": res})
}

// DelRelayByIdHandler 根据id删除relay
// @Summary 根据id删除relay
// @Description 获取用户详细信息
// @Tags 用户
// @Accept  json
// @Produce  json
// @Param   id     path    int     true        "用户ID"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Router /{id} [delete]
func DelRelayByIdHandler(c *gin.Context) {
	id := util.GetInt(c, "id")
	res, err := logic.NewRelayLogic().DelById(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get relay"})
	}
	c.JSON(http.StatusOK, gin.H{"data": res})
}

// UpsertRelayHandler 插入或更新relay
// @Summary 插入或更新relay
// @Description 获取用户详细信息
// @Tags 用户
// @Accept  json
// @Produce  json
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Router /upsert [post]
func UpsertRelayHandler(c *gin.Context) {
	var req resq.UpsertRelayReq
	var err error

	// 绑定普通表单字段
	if err = c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 处理文件上传
	req.File, err = c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "File upload error: " + err.Error()})
		return
	}

	res, err := logic.NewRelayLogic().Save(req.ToModel())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to upsert relay"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": res})
}

// GetRelayByPasswordHandler 根据password获取relay
// @Summary 根据password获取relay
// @Description 获取用户详细信息
// @Tags 用户
// @Accept  json
// @Produce  json
// @Param   pwd     path    int     true        "用户ID"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Router /password/{pwd} [get]
func GetRelayByPasswordHandler(c *gin.Context) {
	pwd := util.GetInt(c, "pwd")
	res, err := logic.NewRelayLogic().GetRelayByPwd(pwd)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get relay"})
	}
	c.JSON(http.StatusOK, gin.H{"data": res})
}

// GetRelayByPasswordAndIdHandler 根据password和id获取relay
// @Summary 根据password和id获取relay
// @Description 获取用户详细信息
// @Tags 用户
// @Accept  json
// @Produce  json
// @Param   id     path    int     true        "用户ID"
// @Param   pwd     path    int     true        "用户ID"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Router /password/{pwd}/id/{id} [get]
func GetRelayByPasswordAndIdHandler(c *gin.Context) {
	pwd := util.GetInt(c, "pwd")
	id := util.GetInt(c, "id")
	res, err := logic.NewRelayLogic().GetRelayByPwdAndId(pwd, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get relay"})
	}
	c.JSON(http.StatusOK, gin.H{"data": res})
}
