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
	postRepository *repository.PostRepository
	tagRepository  *repository.TagRepository
	likeRepository *repository.LikeRepository
	post           domain.Post
}

func NewPostService() *PostService {
	return &PostService{
		postRepository: repository.GetPostRepository(),
		tagRepository:  repository.GetTagRepository(),
		likeRepository: repository.GetLikeRepository(),
	}
}

func (s *PostService) GetAll(c *gin.Context) {
	searchDto := dto.PostSearchDto{}
	utils.BindQuery(c, &searchDto)
	currentUserId := utils.GetCurrentUserId(c)
	if len(searchDto.AccountIds) == 0 {
		searchDto.AccountIds = append(searchDto.AccountIds, currentUserId)
	}
	posts, err := s.postRepository.GetAll(searchDto, currentUserId)
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
	err := s.tagRepository.Create(&post)
	if err != nil {
		log.Println(err)
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}
	_, err = s.postRepository.Create(&post, authorId)
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
	err := s.tagRepository.Update(&post)
	if err != nil {
		log.Println(err)
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}
	_, err = s.postRepository.Update(&post)
	if err != nil {
		log.Println(err)
		c.AbortWithError(http.StatusBadRequest, err)
	} else {
		c.JSON(http.StatusCreated, "&id")
	}
}

func (s *PostService) GetAllComment(c *gin.Context) {
	request := dto.PageRequest{}
	postId := c.Param("postId")
	utils.BindQuery(c, &request)
	log.Printf("Get all commets %v by post %s", request, postId)
	resp, err := s.postRepository.GetAllComments(request, postId)
	if err != nil {
		log.Println(err)
		c.AbortWithError(http.StatusBadRequest, err)
	} else {
		c.JSON(http.StatusOK, resp)
	}
}

func (s *PostService) CreateComment(c *gin.Context) {
	postId := c.Param("postId")
	comment := domain.Comment{}
	utils.BindJson(c, &comment)
	log.Printf("Create comment %v by post %s", comment, postId)
	resp, err := s.postRepository.CreateComment(comment, postId, utils.GetCurrentUserId(c))
	if err != nil {
		log.Println(err)
		c.AbortWithError(http.StatusBadRequest, err)
	} else {
		c.JSON(http.StatusOK, resp)
	}
}

func (s *PostService) GetAllSubComment(c *gin.Context) {
	request := dto.PageRequest{}
	commentId := c.Param("commentId")
	utils.BindQuery(c, &request)
	log.Printf("Get all commets %v by post %s", request, commentId)
	resp, err := s.postRepository.GetAllComments(request, commentId)
	if err != nil {
		log.Println(err)
		c.AbortWithError(http.StatusBadRequest, err)
	} else {
		c.JSON(http.StatusOK, resp)
	}
}

func (s *PostService) CreateLike(c *gin.Context) {
	postId := c.Param("postId")
	log.Printf("Create like on post %s", postId)
	resp, err := s.likeRepository.CreateLike(postId, utils.GetCurrentUserId(c))
	if err != nil {
		log.Println(err)
		c.AbortWithError(http.StatusBadRequest, err)
	} else {
		c.JSON(http.StatusOK, resp)
	}
}

func (s *PostService) DeleteLike(c *gin.Context) {
	postId := c.Param("postId")
	log.Printf("Create like on post %s", postId)
	resp, err := s.likeRepository.DeleteLike(postId, utils.GetCurrentUserId(c))
	if err != nil {
		log.Println(err)
		c.AbortWithError(http.StatusBadRequest, err)
	} else {
		c.JSON(http.StatusOK, resp)
	}
}
