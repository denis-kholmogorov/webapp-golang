package resource

import (
	"github.com/gin-gonic/gin"
	"web/application/service"
)

func FriendResource(e *gin.Engine, service *service.FriendService) {
	e.GET("/api/v1/friends", service.FindAll)
	e.GET("/api/v1/friends/count", service.Count)
	e.POST("/api/v1/friends/:id/request", service.Request)
	e.PUT("/api/v1/friends/:id/approve", service.Approve)
	e.PUT("/api/v1/friends/block/:id", service.Block)
	e.PUT("/api/v1/friends/unblock/:id", service.UnBlock)
	e.DELETE("/api/v1/friends/:id", service.Delete)
	e.GET("/api/v1/friends/recommendations", service.Recommendations)
}
