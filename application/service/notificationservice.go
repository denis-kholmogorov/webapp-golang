package service

import (
	"context"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/segmentio/kafka-go"
	"log"
	"net/http"
	"sync"
	"web/application/domain"
	"web/application/dto"
	kafkaservice "web/application/kafka"
	"web/application/repository"
	"web/application/utils"
)

var notificationService NotificationService
var isInitNotificationService bool

type NotificationService struct {
	notificationRepository *repository.NotificationRepository
	friendRepository       *repository.FriendRepository
	kafkaWriter            *kafka.Writer
}

func NewSettingsService() *NotificationService {
	mt := sync.Mutex{}
	mt.Lock()
	if !isInitNotificationService {
		notificationService = NotificationService{
			notificationRepository: repository.GetSettingsRepository(),
			friendRepository:       repository.GetFriendRepository(),
			kafkaWriter:            kafkaservice.NewWriterMessage("sendNotification"),
		}
		go notificationService.consumeNotification(context.Background())
		isInitNotificationService = true
	}
	mt.Unlock()
	return &notificationService
}

func (s NotificationService) GetAll(c *gin.Context) {
	log.Println("NotificationService:GetAll()")
	settings := s.notificationRepository.GetAll(utils.GetCurrentUserId(c))
	c.JSON(http.StatusOK, settings)
}

func (s NotificationService) GetAllSettings(c *gin.Context) {
	log.Println("NotificationService:GetSettings()")
	settings := s.notificationRepository.GetSettings(utils.GetCurrentUserId(c))
	c.JSON(http.StatusOK, settings)
}

func (s NotificationService) UpdateSettings(c *gin.Context) {
	log.Println("NotificationService:UpdateSettings()")
	item := dto.SettingsItem{}
	utils.BindJson(c, &item)
	settings := s.notificationRepository.GetSettings(utils.GetCurrentUserId(c))
	s.notificationRepository.UpdateSettings(settings, getNameSettings(&item), item.Enable)
	c.JSON(http.StatusOK, settings)
}

func (s NotificationService) consumeNotification(ctx context.Context) {
	r := kafkaservice.NewReaderMessage("createNotifications", "createNotifications")
	for {
		msg, err := r.ReadMessage(ctx)
		if err != nil {
			log.Printf("ERROR: Kafka could not read message %s", err)
			continue
		}
		event := domain.EventNotification{}
		err = json.Unmarshal(msg.Value, &event)
		if err != nil {
			log.Printf("ERROR: Unmarshal kafka messgae %s", err)
			continue
		}
		log.Printf("READ MESSAGE %s", event)
		err = s.createAndSend(&event)
		if err != nil {
			log.Printf("ERROR: Unmarshal kafka messgae %s", err)
			continue
		}
	}
}

func (s NotificationService) createAndSend(event *domain.EventNotification) error {
	log.Println("NotificationService:create()")
	friends, err := s.friendRepository.GetMyFriends(event.InitiatorId)
	if err != nil {
		return err
	}
	notifications := make([]*domain.Notification, len(friends))
	for i, friend := range friends {
		notifications[i] = domain.NewNotification(event.InitiatorId, friend, event.Content, event.NotificationType)
	}

	list, err := s.notificationRepository.SaveAll(notifications)
	if err != nil {
		return err
	}
	for _, notify := range list.List {
		s.sendNotification(notify)
	}
	return nil
}

func (s NotificationService) sendNotification(notify domain.Notification) {
	log.Println("NotificationService:sendNotification()")
	marshal, _ := json.Marshal(notify)
	err := s.kafkaWriter.WriteMessages(context.Background(), kafka.Message{
		Key:   []byte(uuid.NewString()),
		Value: marshal,
	})
	if err != nil {
		log.Printf("ERROR: NotificationService:sendNotification Kafka could not send message %s", marshal)
		log.Printf("ERROR: NotificationService:sendNotification Kafka could not send message %s", err)
	}
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
