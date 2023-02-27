package service

import (
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"web/application/domain"
	"web/application/dto"
	"web/application/repository"
	"web/application/utils"
)

type PostService struct {
	repository *repository.PostRepository
	post       domain.Post
}

func NewPostService() *PostService {
	return &PostService{
		repository: repository.GetPostRepository(),
	}
}

func (s *PostService) GetAll(c *gin.Context) {
	searchDto := dto.PostSearchDto{}
	utils.BindQuery(c, &searchDto)
	if len(searchDto.AccountIds) == 0 {
		searchDto.AccountIds = append(searchDto.AccountIds, utils.GetCurrentUserId(c))
	}
	posts, err := s.repository.GetAll(searchDto)
	if err != nil {
		log.Println(err)
		c.AbortWithError(http.StatusBadRequest, err)
	} else {
		c.JSON(http.StatusOK, posts)
	}
}

func (s *PostService) Create(c *gin.Context) {
	post := domain.Post{}
	utils.BindJson(c, &post)
	log.Printf("Create new post %v", post)
	authorId := utils.GetCurrentUserId(c)
	_, err := s.repository.Create(&post, authorId)
	if err != nil {
		log.Println(err)
		c.AbortWithError(http.StatusBadRequest, err)
	} else {
		c.JSON(http.StatusCreated, "&id")
	}
}

func (s *PostService) Update(c *gin.Context) {
	post := domain.Post{}
	utils.BindJson(c, &post)
	log.Printf("Create new post %v", post)
	_, err := s.repository.Update(&post)
	if err != nil {
		log.Println(err)
		c.AbortWithError(http.StatusBadRequest, err)
	} else {
		c.JSON(http.StatusCreated, "&id")
	}
}

func (s *PostService) GetAllComments(c *gin.Context) {
	request := dto.PageRequest{}
	postId := c.Param("postId")
	utils.BindQuery(c, &request)
	log.Printf("Get all commets %v by post %s", request, postId)
	resp, err := s.repository.GetAllComments(request, postId)
	if err != nil {
		log.Println(err)
		c.AbortWithError(http.StatusBadRequest, err)
	} else {
		c.JSON(http.StatusOK, resp)
	}
}
