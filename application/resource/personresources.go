package resource

import (
	"github.com/gin-gonic/gin"
	"web/application/service"
)

func PersonResource(e *gin.Engine, service *service.PersonService) {
	e.GET("/person", service.GetAll)
	e.GET("/person/fields", service.GetAllFields)
	e.GET("/person/:id", service.GetById)
	e.POST("/person", service.Create)
	e.PUT("/person", service.Update)
	e.DELETE("/person/:id", service.DeleteById)
}
