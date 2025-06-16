package utils

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"web/application/errorhandler"
)

func BindJson(c *gin.Context, value any) {
	err := c.BindJSON(value)
	if err != nil {
		panic(errorhandler.ErrorResponse{Message: fmt.Sprintf("Utils:BindJson() Can't parse json from request %s", err)})
	}
}

func BindQuery(c *gin.Context, value any) {
	err := c.BindQuery(value)
	if err != nil {
		panic(errorhandler.ErrorResponse{Message: fmt.Sprintf("Utils:BindQuery() Can't parse query from request %s", err)})
	}
}
