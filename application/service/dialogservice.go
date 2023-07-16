package service

import (
	"context"
	"encoding/json"
	"fmt"
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

var dialogService DialogService
var isInitDialogService bool

type DialogService struct {
	dialogRepository *repository.DialogRepository
	kafkaWriter      *kafka.Writer
}

func NewDialogService() *DialogService {
	mt := sync.Mutex{}
	mt.Lock()
	if !isInitDialogService {
		dialogService = DialogService{
			dialogRepository: repository.GetDialogRepository(),
			kafkaWriter:      kafkaservice.NewWriterMessage("sendMessage"),
		}
		go dialogService.ConsumeMessage(context.Background())
		isInitDialogService = true
	}
	mt.Unlock()
	return &dialogService
}

func (s DialogService) Unread(c *gin.Context) {
	log.Println("DialogService:GetCountry()")
	c.JSON(http.StatusOK, dto.Count{TotalElement: 1})

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

func (s DialogService) SaveMessage(message *domain.Message) (*domain.Message, error) {
	return s.dialogRepository.SaveMessage(message)

}

func (s DialogService) ConsumeMessage(ctx context.Context) {
	r := kafkaservice.NewReaderMessage("receiveMessage", "receiveMessage")
	for {
		msg, err := r.ReadMessage(ctx)
		if err != nil {
			log.Printf("ERROR: Kafka could not read message " + err.Error())
			break
		}

		message := domain.Message{}
		err = json.Unmarshal(msg.Value, &message)
		savedMessage, err := s.SaveMessage(&message)
		if err != nil {
			log.Println("Kafka can't read message")
			break
		}
		err = s.SendMessage(savedMessage)
		if err != nil {
			log.Print("ERROR: Kafka could not read message ", err)
			break
		}
	}
}

func (s *DialogService) SendMessage(message *domain.Message) error {
	marshal, _ := json.Marshal(message)
	err := s.kafkaWriter.WriteMessages(context.Background(), kafka.Message{
		Key:   []byte(uuid.NewString()),
		Value: marshal,
	})
	if err != nil {
		log.Print("ERROR: Kafka could not read message ", err)
		return err
	}
	return nil
}
