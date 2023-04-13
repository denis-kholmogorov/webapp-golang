package service

import (
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

func (s *FriendService) Request(c *gin.Context) {
	friendId := c.Param("id")
	s.friendRepository.RequestFriend(utils.GetCurrentUserId(c), friendId)
	c.Status(http.StatusOK)
}

func (s *FriendService) Approve(c *gin.Context) {
	friendId := c.Param("id")
	s.friendRepository.ApproveFriend(utils.GetCurrentUserId(c), friendId)
	c.JSON(http.StatusOK, "")
}

func (s *FriendService) FindAll(c *gin.Context) {
	page := dto.PageRequest{}
	statusCode := dto.StatusCode{}
	utils.BindQuery(c, &page)
	utils.BindQuery(c, &statusCode)
	response := s.friendRepository.FindAll(utils.GetCurrentUserId(c), statusCode, page)
	c.JSON(http.StatusOK, response)
}

func (s *FriendService) Delete(c *gin.Context) {
	friendId := c.Param("id")
	s.friendRepository.Delete(utils.GetCurrentUserId(c), friendId)
	c.JSON(http.StatusOK, "")

}

func (s *FriendService) Count(c *gin.Context) {
	count := s.friendRepository.Count(utils.GetCurrentUserId(c))
	c.JSON(http.StatusOK, count)

}

func (s *FriendService) Recommendations(c *gin.Context) {
	response := s.friendRepository.Recommendations(utils.GetCurrentUserId(c))
	c.JSON(http.StatusOK, response)
}

func (s *FriendService) Block(c *gin.Context) {
	friendId := c.Param("id")
	s.friendRepository.Block(utils.GetCurrentUserId(c), friendId)
	c.Status(http.StatusOK)
}
