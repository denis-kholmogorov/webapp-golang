package service

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"sync"
	"web/application/domain"
	"web/application/dto"
	"web/application/repository"
	"web/application/utils"
)

var dialogService DialogService
var isInitDialogService bool

type DialogService struct {
	dialogRepository *repository.DialogRepository
}

func NewDialogService() *DialogService {
	mt := sync.Mutex{}
	mt.Lock()
	if !isInitDialogService {
		dialogService = DialogService{
			dialogRepository: repository.GetDialogRepository(),
		}
		isInitDialogService = true
	}
	mt.Unlock()
	return &dialogService
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
	currentUserId := utils.GetCurrentUserId(c)
	messages, err := s.dialogRepository.GetMessages(currentUserId, page)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, fmt.Sprintf("GetMessages() Row with not found"))
	} else {
		c.JSON(http.StatusOK, messages)
	}
}

func (s DialogService) SaveMessage(message *domain.Message) (*domain.Message, error) {
	return s.dialogRepository.SaveMessage(message)

}

func (s DialogService) UpdateDialogs(c *gin.Context) {
	companionId := c.Param("companionId")
	currentUserId := utils.GetCurrentUserId(c)
	err := s.dialogRepository.UpdateMessages(currentUserId, companionId)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, fmt.Sprintf("UpdateDialogs() Row with not found"))
	} else {
		c.JSON(http.StatusOK, "")
	}
}
