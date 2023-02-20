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

	searchDto.AuthorId = utils.GetCurrentUserId(c)
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
	value, _ := c.Get("id")
	authorId := value.(string)
	_, err := s.repository.Create(&post, authorId)
	if err != nil {
		log.Println(err)
		c.AbortWithError(http.StatusBadRequest, err)
	} else {
		c.JSON(http.StatusCreated, "&id")
	}
}
