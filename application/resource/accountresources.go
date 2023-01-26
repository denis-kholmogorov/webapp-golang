package resource

import (
	"github.com/gin-gonic/gin"
	"web/application/service"
)

func AccountResource(e *gin.Engine, service *service.AccountService) {
	e.GET("api/v1/account", service.GetAll)
	e.GET("api/v1/account/fields", service.GetAllFields)
	e.GET("api/v1/account/me", service.GetById)
	e.GET("api/v1/account/:id", service.GetById)
	e.POST("api/v1/account", service.Create)
	e.PUT("api/v1/account", service.Update)
	e.DELETE("api/v1/account/:id", service.DeleteById)
}
