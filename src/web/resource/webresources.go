package resource

import (
	"github.com/gin-gonic/gin"
	"web/service"
)

func CreatePaths(e *gin.Engine, service *service.PersonService) {
	//e.GET("/model", service.GetAll)
	//e.GET("/model/:id", service.GetById)
	e.POST("/model", service.Create)
	//e.PUT("/model", service.Update)
	//e.DELETE("/model/:id", service.DeleteById)
}
