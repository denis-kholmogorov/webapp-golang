package resource

import (
	"github.com/gin-gonic/gin"
	"web/application/service"
)

func NotificationResource(e *gin.Engine, service *service.NotificationService) {
	e.GET("/api/v1/notifications", service.GetAll)
	e.GET("/api/v1/notifications/settings", service.GetAllSettings)
	e.PUT("/api/v1/notifications/settings", service.UpdateSettings)
}
