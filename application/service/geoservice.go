package service

import (
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
	countries := s.geoRepository.FindAll()
	c.JSON(http.StatusOK, countries)
}

func (s GeoService) GetCity(c *gin.Context) {
	id := c.Param("countryId")
	cities := s.geoRepository.FindCitiesByCountryId(id)
	c.JSON(http.StatusOK, cities)

}
