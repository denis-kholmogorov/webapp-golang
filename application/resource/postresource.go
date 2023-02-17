package resource

import (
	"github.com/gin-gonic/gin"
	"web/application/service"
)

func PostResource(e *gin.Engine, service *service.PostService) {
	e.POST("api/v1/post", service.Create)
	e.GET("api/v1/post", service.GetAll)
}
