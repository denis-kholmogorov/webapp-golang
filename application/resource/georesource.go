package resource

import (
	"github.com/gin-gonic/gin"
	"web/application/service"
)

func GeoResource(e *gin.Engine, service *service.GeoService) {
	e.GET("/api/v1/geo/country", service.GetCountry)
	e.GET("/api/v1/geo/country/:countryId/city", service.GetCity)
}
