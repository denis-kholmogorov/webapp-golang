package service

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"web/application/repository"
)

type GeoService struct {
	geoRepository *repository.GeoRepository
}

func NewGeoService() *GeoService {
	return &GeoService{
		geoRepository: repository.GetGeoRepository(),
	}
}

func (s GeoService) GetCountry(c *gin.Context) {
	log.Println("GeoService:GetCountry()")
	countries, err := s.geoRepository.FindAll()
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, fmt.Sprintf("GetCountry() Row with not found"))
	} else {
		c.JSON(http.StatusOK, countries)
	}

}

func (s GeoService) GetCity(c *gin.Context) {
	id := c.Param("countryId")
	cities, err := s.geoRepository.FindCitiesByCountryId(id)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, fmt.Sprintf("GetCountry() Row with not found"))
	} else {
		c.JSON(http.StatusOK, cities)
	}
}
