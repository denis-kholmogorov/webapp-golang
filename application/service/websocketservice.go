package service

import (
	"context"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/segmentio/kafka-go"
	"log"
	"net/http"
	"sync"
	"web/application/domain"
	"web/application/dto"
	kafkaservice "web/application/kafka"
	"web/application/utils"
)

var websocketService WebsocketService
var isInitWebsocketService bool

type WebsocketService struct {
	connections map[string]*websocket.Conn
	kafkaWriter *kafka.Writer
}

func NewWebsocketService() *WebsocketService {
	mt := sync.Mutex{}
	mt.Lock()
	if !isInitWebsocketService {
		websocketService = WebsocketService{
			connections: map[string]*websocket.Conn{},
			kafkaWriter: kafkaservice.NewWriterMessage("receiveMessage"),
		}
		go websocketService.SendMessage(context.Background())
		isInitWebsocketService = true
	}
	mt.Unlock()
	return &websocketService
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

func (s *WebsocketService) Connect(c *gin.Context) {
	ws, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Print("ERROR::Websocket upgrade websocket:", err)
		return
	}
	defer s.deleteConn(utils.GetCurrentUserId(c))
	defer ws.Close()

	log.Printf("Websocket before add account %d", len(s.connections))
	s.connections[utils.GetCurrentUserId(c)] = ws
	log.Printf("Websocket after add account %d", len(s.connections))

	for {
		socketDto := dto.SocketDto[domain.Message]{}
		err = ws.ReadJSON(&socketDto)
		if err != nil {
			log.Println("ERROR::WebsocketService.Conn() read message:", err)
			break
		}
		err = s.ReceiveMessage(socketDto.Data)
		if err != nil {
			log.Println("ERROR: Websocket kafka send message:", err)
			break
		}
	} //write ws data
}

func (s *WebsocketService) ReceiveMessage(message *domain.Message) error {
	marshal, err := json.Marshal(message)
	if err != nil {
		log.Println("ERROR: WebsocketService.ReceiveMessage() cannot marshal", err)
		return err
	}
	err = s.kafkaWriter.WriteMessages(context.Background(), kafka.Message{
		Key:   []byte(uuid.NewString()),
		Value: marshal,
	})
	if err != nil {
		log.Print("ERROR: WebsocketService.ReceiveMessage() could not read message ", err)
		return err
	}
	return nil
}

func (s *WebsocketService) SendSocketMessage(savedMessage *domain.Message) error {

	if s.connections[savedMessage.RecipientId] != nil {
		err := s.connections[savedMessage.RecipientId].WriteJSON(dto.NewMessageSocketDto(savedMessage))
		if err != nil {
			log.Printf("ERROR: Websocket send message to account %s %v:", savedMessage.RecipientId, err)
			return err
		}
	}
	return nil
}

func (s *WebsocketService) SendMessage(ctx context.Context) {
	r := kafkaservice.NewReaderMessage("sendMessage", "sendMessage")
	for {
		msg, err := r.ReadMessage(ctx)
		if err != nil {
			log.Printf("ERROR: Kafka could not read message " + err.Error())
			break
		}

		message := domain.Message{}
		err = json.Unmarshal(msg.Value, &message)
		if err != nil {
			log.Println("Kafka can't read message")
			break
		}
		err = s.SendSocketMessage(&message)
		if err != nil {
			log.Print("ERROR: Kafka could not send message ", err)
			break
		}
	}
}

func (s *WebsocketService) deleteConn(id string) {
	log.Printf("Websocket before deleted account%d", len(s.connections))
	delete(s.connections, id)
	log.Printf("Websocket after deleted account%d", len(s.connections))

}
