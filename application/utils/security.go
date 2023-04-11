package utils

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"web/application/errorhandler"
)

func GetCurrentUserId(c *gin.Context) string {
	value, ok := c.Get("id")
	if !ok {
		panic(errorhandler.ErrorResponse{Message: fmt.Sprintf("Гtils:GetCurrentUserId() Error get id from context")})
	}
	return value.(string)
}
