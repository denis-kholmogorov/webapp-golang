package utils

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"web/application/errors"
)

func GetCurrentUserId(c *gin.Context) string {
	value, ok := c.Get("id")
	if ok {
		return value.(string)
	}
	c.AbortWithStatusJSON(http.StatusUnauthorized, errors.ErrorDescription{ErrorDescription: "Not found user id"})
	return ""

}
