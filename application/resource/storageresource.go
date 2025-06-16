package resource

import (
	"github.com/gin-gonic/gin"
	"web/application/service"
)

func StorageResource(e *gin.Engine, service *service.StorageService) {
	e.POST("/api/v1/storage", service.Upload)
}
