package service

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/robfig/cron"
	"github.com/segmentio/kafka-go"
	"log"
	"net/http"
	"sync"
	"web/application/domain"
	"web/application/dto"
	"web/application/errorhandler"
	kafkaservice "web/application/kafka"
	"web/application/repository"
	"web/application/utils"
)

var postService PostService
var isInitPostService bool

type PostService struct {
	postRepository *repository.PostRepository
	tagRepository  *repository.TagRepository
	likeRepository *repository.LikeRepository
	kafkaWriter    *kafka.Writer
}

func NewPostService() *PostService {
	mt := sync.Mutex{}
	mt.Lock()
	postService = PostService{
		postRepository: repository.GetPostRepository(),
		tagRepository:  repository.GetTagRepository(),
		likeRepository: repository.GetLikeRepository(),
		kafkaWriter:    kafkaservice.NewWriterMessage("createNotifications"),
	}
	isInitPostService = true
	go createCronPostPublish()
	mt.Unlock()
	return &postService
}

func (s *PostService) GetAll(c *gin.Context) {
	searchDto := dto.PostSearchDto{}
	utils.BindQuery(c, &searchDto)
	currentUserId := utils.GetCurrentUserId(c)

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
	tagIds, err := s.tagRepository.Update(&post)
	if err != nil {
		log.Println(err)
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}
	newPost, err := s.postRepository.Create(&post, authorId, tagIds)
	s.SendNotification(&domain.EventNotification{InitiatorId: post.AuthorId, Content: utils.TrimStr(newPost.Title), NotificationType: "POST"})
	if err != nil {
		log.Println(err)
		c.AbortWithError(http.StatusBadRequest, err)
	} else {
		c.Status(http.StatusCreated)
	}
}

func (s *PostService) Update(c *gin.Context) {
	post := domain.Post{}
	utils.BindJson(c, &post)
	log.Printf("Create new post %v", post)
	tagIds, err := s.tagRepository.Update(&post)
	if err != nil {
		log.Println(err)
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}
	err = s.postRepository.Update(&post, tagIds)
	if err != nil {
		log.Println(err)
		c.AbortWithError(http.StatusBadRequest, err)
	} else {
		c.JSON(http.StatusCreated, "")
	}
}

func (s *PostService) Delete(c *gin.Context) {
	postId := c.Param("postId")
	log.Printf("Delete post %s", postId)
	err := s.postRepository.Delete(postId)
	if err != nil {
		log.Println(err)
		c.AbortWithError(http.StatusBadRequest, err)
	} else {
		c.JSON(http.StatusOK, "")
	}
}

func (s *PostService) GetAllComment(c *gin.Context) {
	request := dto.PageRequest{}
	postId := c.Param("postId")
	utils.BindQuery(c, &request)
	log.Printf("Get all commets %v by post %s", request, postId)
	resp, err := s.postRepository.GetAllComments(request, postId, utils.GetCurrentUserId(c))
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

func (s *PostService) UpdateComment(c *gin.Context) {
	comment := domain.Comment{}
	utils.BindJson(c, &comment)
	log.Printf("UpdateSettings comment %v", comment.Id)
	err := s.postRepository.UpdateComment(comment)
	if err != nil {
		log.Println(err)
		c.AbortWithError(http.StatusBadRequest, err)
	} else {
		c.JSON(http.StatusOK, "")
	}
}

func (s *PostService) DeleteComment(c *gin.Context) {
	commentId := c.Param("commentId")
	log.Printf("Delete post %s", commentId)
	err := s.postRepository.Delete(commentId)
	if err != nil {
		log.Println(err)
		c.AbortWithError(http.StatusBadRequest, err)
	} else {
		c.JSON(http.StatusOK, "")
	}
}

func (s *PostService) GetAllSubComment(c *gin.Context) {
	request := dto.PageRequest{}
	commentId := c.Param("commentId")
	utils.BindQuery(c, &request)
	log.Printf("Get all commets %v by post %s", request, commentId)
	resp, err := s.postRepository.GetAllComments(request, commentId, utils.GetCurrentUserId(c))
	if err != nil {
		log.Println(err)
		c.AbortWithError(http.StatusBadRequest, err)
	} else {
		c.JSON(http.StatusOK, resp)
	}
}

func (s *PostService) CreateLike(c *gin.Context) {
	parentId := c.Param("postId")
	log.Printf("Create like on parentId %s", parentId)
	resp, err := s.likeRepository.CreateLike(parentId, utils.GetCurrentUserId(c))
	if err != nil {
		log.Println(err)
		c.AbortWithError(http.StatusBadRequest, err)
	} else {
		c.JSON(http.StatusOK, resp)
	}
}

func (s *PostService) DeleteLike(c *gin.Context) {
	parentId := c.Param("postId")
	log.Printf("Create like on post %s", parentId)
	resp, err := s.likeRepository.DeleteLike(parentId, utils.GetCurrentUserId(c))
	if err != nil {
		log.Println(err)
		c.AbortWithError(http.StatusBadRequest, err)
	} else {
		c.JSON(http.StatusOK, resp)
	}
}

func (s *PostService) CreateCommentLike(c *gin.Context) {
	parentId := c.Param("commentId")
	log.Printf("Create like on parentId %s", parentId)
	resp, err := s.likeRepository.CreateLike(parentId, utils.GetCurrentUserId(c))
	if err != nil {
		log.Println(err)
		c.AbortWithError(http.StatusBadRequest, err)
	} else {
		c.JSON(http.StatusOK, resp)
	}
}

func (s *PostService) DeleteCommentLike(c *gin.Context) {
	parentId := c.Param("commentId")
	log.Printf("Create like on parentId %s", parentId)
	resp, err := s.likeRepository.DeleteLike(parentId, utils.GetCurrentUserId(c))
	if err != nil {
		log.Println(err)
		c.AbortWithError(http.StatusBadRequest, err)
	} else {
		c.JSON(http.StatusOK, resp)
	}
}

func (s *PostService) SendNotification(message *domain.EventNotification) {
	marshal, err := json.Marshal(message)
	if err != nil {
		panic(errorhandler.ErrorResponse{Message: fmt.Sprintf("PostService:SendNotification()  Kafka could not send message  %s", err)})
	}
	err = s.kafkaWriter.WriteMessages(context.Background(), kafka.Message{
		Key:   []byte(uuid.NewString()),
		Value: marshal,
	})
	if err != nil {
		panic(errorhandler.ErrorResponse{Message: fmt.Sprintf("PostService:SendNotification()  Kafka could not send message  %s", err)})
	}
}

func (s *PostService) checkPublishPost() {
	s.postRepository.UpdatePostQueue()
}

func createCronPostPublish() {
	c := cron.New()
	err := c.AddFunc("*/15 * * * *", postService.checkPublishPost)
	if err != nil {
		fmt.Println("3 Time CRON")
	}
	c.Run()
	fmt.Println("4 Time CRON")
}
