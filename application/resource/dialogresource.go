package resource

import (
	"github.com/gin-gonic/gin"
	"web/application/service"
)

func DialogResource(e *gin.Engine, service *service.DialogService) {
	e.GET("api/v1/dialogs/unreaded", service.Unread)
	e.GET("api/v1/dialogs", service.GetDialogs)
	e.GET("api/v1/dialogs/messages", service.GetMessages)
	//e.GET("api/v1/dialogs/", service.FindById)

}
