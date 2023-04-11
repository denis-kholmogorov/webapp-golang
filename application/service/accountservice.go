package service

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"web/application/domain"
	"web/application/dto"
	"web/application/repository"
	"web/application/utils"
)

type AccountService struct {
	accountRepository *repository.AccountRepository
	person            domain.Account
}

func NewAccountService() *AccountService {
	return &AccountService{
		accountRepository: repository.GetAccountRepository(),
	}
}

func (s *AccountService) GetMe(c *gin.Context) {
	account := s.accountRepository.FindById(utils.GetCurrentUserId(c))
	c.JSON(http.StatusOK, account)

}

func (s *AccountService) UpdateMe(c *gin.Context) {
	id := utils.GetCurrentUserId(c)
	account := domain.Account{Id: id, Uid: id}
	utils.BindJson(c, &account)
	s.accountRepository.Update(&account)
	c.JSON(http.StatusOK, account)

}

func (s *AccountService) FindById(c *gin.Context) {
	id := c.Param("id")
	account := s.accountRepository.FindById(id)
	c.JSON(http.StatusOK, account)
}

func (s *AccountService) FindAll(c *gin.Context) {
	searchDto := dto.AccountSearchDto{}
	utils.BindQuery(c, &searchDto)

	account := s.accountRepository.FindAll(searchDto, utils.GetCurrentUserId(c))
	c.JSON(http.StatusOK, account)

}
