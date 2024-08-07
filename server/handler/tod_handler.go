package handler

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"onij/logic"
)

func GetWeeklyTodHandler(c *gin.Context) {

	res, err := logic.NewTodLogic().GetWeeklyTodList()
	if err != nil {
		c.JSON(http.StatusInternalServerError, "")
		c.Abort()
	}

	c.JSON(http.StatusOK, res)
}
