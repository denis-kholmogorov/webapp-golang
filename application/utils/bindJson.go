package utils

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

func BindJson(c *gin.Context, value any) {
	err := c.BindJSON(value)
	if err != nil {
		fmt.Printf("Can't parse json from request %s", err.Error())
		c.AbortWithStatus(http.StatusBadRequest)
	}
}

func BindQuery(c *gin.Context, value any) {
	err := c.BindQuery(value)
	if err != nil {
		fmt.Printf(err.Error())
		c.AbortWithStatus(http.StatusBadRequest)
	}
}
