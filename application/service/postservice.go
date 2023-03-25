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
	ids, err := s.tagRepository.Update(&post)
	if err != nil {
		log.Println(err)
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}
	err = s.postRepository.Update(&post, ids)
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
	log.Printf("Update comment %v", comment.Id)
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
