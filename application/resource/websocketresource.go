package resource

import (
	"github.com/gin-gonic/gin"
	"web/application/service"
)

func WebsocketResource(e *gin.Engine, service *service.WebsocketService) {
	e.GET("api/v1/streaming/ws", service.Connect)
}
