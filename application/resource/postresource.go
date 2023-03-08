package resource

import (
	"github.com/gin-gonic/gin"
	"web/application/service"
)

func PostResource(e *gin.Engine, service *service.PostService) {
	e.GET("api/v1/post", service.GetAll)
	e.POST("api/v1/post", service.Create)
	e.PUT("api/v1/post", service.Update)
	e.GET("api/v1/post/:postId/comment", service.GetAllComment)
	e.GET("api/v1/post/:postId/comment/:commentId/subcomment", service.GetAllSubComment)
	e.POST("api/v1/post/:postId/comment", service.CreateComment)
	e.POST("api/v1/post/:postId/like", service.CreateLike)
	e.DELETE("api/v1/post/:postId/like", service.DeleteLike)
	e.POST("api/v1/post/:postId/comment/:commentId/like", service.CreateCommentLike)
	e.DELETE("api/v1/post/:postId/comment/:commentId/like", service.DeleteCommentLike)

}
