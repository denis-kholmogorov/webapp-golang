package resource

import (
	"github.com/gin-gonic/gin"
	"web/application/service"
)

func AuthResource(e *gin.Engine, service *service.AuthService) {
	e.POST("/login", service.Login)

}
