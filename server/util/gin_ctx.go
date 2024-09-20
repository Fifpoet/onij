package util

import (
	"github.com/gin-gonic/gin"
	"strconv"
)

func GetInt(c *gin.Context, key string) int {
	idStr := c.Param(key)
	id, _ := strconv.Atoi(idStr)
	return id
}
