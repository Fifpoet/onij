package handler

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"onij/infra"
	"onij/logic"
	"strconv"
)

type RelayHandler interface {
}

type relayHandler struct {
	infra *infra.AllInfra
	logic *logic.AllLogic
}

func NewRelayHandler(i *infra.AllInfra, l *logic.AllLogic) RelayHandler {
	return &relayHandler{
		infra: i,
		logic: l,
	}
}

func (r *relayHandler) GetRelayHandler(c *gin.Context) {
	id := c.Param("type")
	typ, err := strconv.Atoi(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid id"})
		return
	}
	res, err := r.logic.RelayLogic.GetRelayByType(typ)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get relay"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": res})
}
