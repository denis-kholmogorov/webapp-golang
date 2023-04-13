package service

import (
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"web/application/dto"
	"web/application/repository"
	"web/application/utils"
)

type SettingsService struct {
	settingsRepository *repository.SettingsRepository
}

func NewSettingsService() *SettingsService {
	return &SettingsService{
		settingsRepository: repository.GetSettingsRepository(),
	}
}

func (s SettingsService) GetAll(c *gin.Context) {
	log.Println("SettingsService:GetSettings()")
	settings := s.settingsRepository.GetSettings(utils.GetCurrentUserId(c))
	c.JSON(http.StatusOK, settings)
}

func (s SettingsService) Update(c *gin.Context) {
	log.Println("SettingsService:UpdateSettings()")
	item := dto.SettingsItem{}
	utils.BindJson(c, &item)
	settings := s.settingsRepository.GetSettings(utils.GetCurrentUserId(c))
	s.settingsRepository.UpdateSettings(settings, getNameSettings(&item), item.Enable)
	c.JSON(http.StatusOK, settings)
}

func getNameSettings(item *dto.SettingsItem) string {
	switch item.NotificationType {
	case "POST":
		return "enablePost"
	case "POST_COMMENT":
		return "enablePostComment"
	case "COMMENT_COMMENT":
		return "enableCommentComment"
	case "FRIEND_REQUEST":
		return "enableMessage"
	case "MESSAGE":
		return "enableFriendRequest"
	case "FRIEND_BIRTHDAY":
		return "enableFriendBirthday"
	case "SEND_EMAIL_MESSAGE":
		return "enableSendEmailMessage"
	}
	return ""
}
