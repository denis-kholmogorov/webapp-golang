package service

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"web/application/domain"
	"web/application/dto"
	"web/application/repository"
	"web/application/utils"
)

type FriendService struct {
	friendRepository *repository.FriendRepository
	friendship       domain.Friendship
}

func NewFriendService() *FriendService {
	return &FriendService{
		friendRepository: repository.GetFriendRepository(),
	}
}

func (s *FriendService) RequestFriend(c *gin.Context) {
	friendId := c.Param("id")
	err := s.friendRepository.RequestFriend(utils.GetCurrentUserId(c), friendId)
	if err != nil {
		c.JSON(http.StatusBadRequest, fmt.Sprintf("Row with %s not found", utils.GetCurrentUserId(c)))
	} else {
		c.JSON(http.StatusOK, "")
	}
}

func (s *FriendService) ApproveFriend(c *gin.Context) {
	friendId := c.Param("id")
	err := s.friendRepository.ApproveFriend(utils.GetCurrentUserId(c), friendId)
	if err != nil {
		c.JSON(http.StatusBadRequest, fmt.Sprintf("Row with %s not found", utils.GetCurrentUserId(c)))
	} else {
		c.JSON(http.StatusOK, "")
	}
}

func (s *FriendService) FindAll(c *gin.Context) {
	page := dto.PageRequest{}
	statusCode := dto.StatusCode{}
	utils.BindQuery(c, &page)
	utils.BindQuery(c, &statusCode)
	response, err := s.friendRepository.FindAll(utils.GetCurrentUserId(c), statusCode, page)
	if err != nil {
		c.JSON(http.StatusBadRequest, fmt.Sprintf("Row with %s not found", utils.GetCurrentUserId(c)))
	} else {
		c.JSON(http.StatusOK, response)
	}
}