package resource

import (
	"github.com/gin-gonic/gin"
	"web/application/service"
)

func SettingsResource(e *gin.Engine, service *service.SettingsService) {
	e.GET("/api/v1/notifications/settings", service.GetAll)
	e.PUT("/api/v1/notifications/settings", service.Update)
}
