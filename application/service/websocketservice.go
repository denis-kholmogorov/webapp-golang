package service

import (
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"sync"
	"web/application/domain"
	"web/application/dto"
	"web/application/utils"
)

var websocketService WebsocketService
var isInitWebsocketService bool

type WebsocketService struct {
	dialogService *DialogService
	connections   map[string]*websocket.Conn
}

func NewWebsocketService() *WebsocketService {
	mt := sync.Mutex{}
	mt.Lock()
	if !isInitWebsocketService {
		websocketService = WebsocketService{
			NewDialogService(),
			map[string]*websocket.Conn{},
		}
		isInitWebsocketService = true
	}
	mt.Unlock()
	return &websocketService
}

var upgrader = websocket.Upgrader{
	// Solve cross-domain problems
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
} // use default options

func (s *WebsocketService) Connect(c *gin.Context) {
	//Upgrade get request to webSocket protocol
	ws, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Print("upgrade:", err)
		return
	}
	defer ws.Close()
	s.connections[utils.GetCurrentUserId(c)] = ws

	for {
		//read data from ws
		socketDto := dto.SocketDto[domain.Message]{}
		err := ws.ReadJSON(&socketDto)
		if err != nil {
			log.Println("read:", err)
			break
		}

		saveMessage, err := s.dialogService.SaveMessage(socketDto.Data)

		if s.connections[saveMessage.RecipientId] != nil {
			err = s.connections[saveMessage.RecipientId].WriteJSON(dto.NewMessageSocketDto(saveMessage))
			if err != nil {
				log.Println("write:", err)
				break
			}
		} //write ws data
	}
}
