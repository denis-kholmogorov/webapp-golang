package service

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"sync"
	"web/application/dto"
	"web/application/repository"
	"web/application/utils"
)

var tagService TagService
var isInitTagService bool

type TagService struct {
	tagRepository *repository.TagRepository
}

func NewTagService() *TagService {
	mt := sync.Mutex{}
	mt.Lock()
	tagService = TagService{
		tagRepository: repository.GetTagRepository(),
	}
	isInitTagService = true
	mt.Unlock()
	return &tagService
}

func (s *TagService) GetByName(c *gin.Context) {
	searchDto := dto.TagSearchDto{}
	utils.BindQuery(c, &searchDto)
	tags := s.tagRepository.FindAll(searchDto)
	c.JSON(http.StatusOK, tags)

}
