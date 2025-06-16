package resource

import (
	"github.com/gin-gonic/gin"
	"web/application/service"
)

func TagResource(e *gin.Engine, service *service.TagService) {

	e.GET("api/v1/tag", service.GetByName)
}
