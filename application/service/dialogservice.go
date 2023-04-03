package service

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"web/application/dto"
	"web/application/repository"
	"web/application/utils"
)

type DialogService struct {
	dialogRepository *repository.DialogRepository
}

func NewDialogService() *DialogService {
	return &DialogService{
		dialogRepository: repository.GetDialogRepository(),
	}
}

func (s DialogService) Unread(c *gin.Context) {
	log.Println("DialogService:GetCountry()")
	c.JSON(http.StatusOK, "{\"data\":{\"count\":\"1\"}}")

}

func (s DialogService) GetDialogs(c *gin.Context) {
	page := dto.PageRequestOld{}
	utils.BindQuery(c, &page)
	authorId := utils.GetCurrentUserId(c)
	cities, err := s.dialogRepository.GetDialogs(authorId, page)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, fmt.Sprintf("GetDialogs() Row with not found"))
	} else {
		c.JSON(http.StatusOK, cities)
	}
}

func (s DialogService) GetMessages(c *gin.Context) {
	page := dto.PageRequestOld{}
	utils.BindQuery(c, &page)
	authorId := utils.GetCurrentUserId(c)
	messages, err := s.dialogRepository.GetMessages(authorId, page)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, fmt.Sprintf("GetMessages() Row with not found"))
	} else {
		c.JSON(http.StatusOK, messages)
	}
}
