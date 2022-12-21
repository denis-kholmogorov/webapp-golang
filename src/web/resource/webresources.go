package resource

import (
	"github.com/gin-gonic/gin"
	"web/service"
)

func CreatePaths(e *gin.Engine, service *service.PersonService) {
	e.GET("/person", service.GetAll)
	e.GET("/person/:id", service.GetById)
	e.POST("/person", service.Create)
	e.PUT("/person", service.Update)
	e.DELETE("/person/:id", service.DeleteById)
}
