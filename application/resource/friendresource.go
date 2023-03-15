package resource

import (
	"github.com/gin-gonic/gin"
	"web/application/service"
)

func FriendResource(e *gin.Engine, service *service.FriendService) {
	e.GET("/api/v1/friends", service.FindAll)
	e.POST("/api/v1/friends/:id/request", service.RequestFriend)
}
