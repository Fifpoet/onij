package handler

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"onij/logic"
	"strconv"
)

func GetRelayHandler(c *gin.Context) {
	id := c.Param("type")
	typ, err := strconv.Atoi(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid id"})
		return
	}
	res, err := logic.NewRelayLogic().GetRelayByType(typ)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get relay"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": res})
}
