package resource

import (
	"github.com/gin-gonic/gin"
	"web/application/service"
)

func PostResource(e *gin.Engine, service *service.PostService) {
	e.GET("api/v1/post", service.GetAll)
	e.POST("api/v1/post", service.Create)
	e.PUT("api/v1/post", service.Update)
	e.GET("api/v1/post/:postId/comment", service.GetAllComments)

}
