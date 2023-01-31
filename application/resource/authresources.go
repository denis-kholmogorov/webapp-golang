package resource

import (
	"github.com/gin-gonic/gin"
	"web/application/service"
)

func AuthResource(e *gin.Engine, service *service.AuthService) {
	e.GET("api/v1/auth/captcha", service.GetCaptcha)
	e.POST("api/v1/auth/register", service.Registration)
	//e.POST("api/v1/auth/login", service.Login)
	//e.POST("api/v1/auth/logout", service.Logout)

}
