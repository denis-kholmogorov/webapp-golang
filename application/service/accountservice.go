package service

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"web/application/domain"
	"web/application/dto"
	"web/application/errors"
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
	account, err := s.accountRepository.FindById(utils.GetCurrentUserId(c))
	if err != nil {
		c.JSON(http.StatusBadRequest, fmt.Sprintf("Row with %s not found", utils.GetCurrentUserId(c)))
	} else {
		c.JSON(http.StatusOK, account)
	}
}

func (s *AccountService) UpdateMe(c *gin.Context) {
	id := utils.GetCurrentUserId(c)
	account := domain.Account{Id: id, Uid: id}
	utils.BindJson(c, &account)
	email, err := s.accountRepository.Update(&account)

	if err != nil {
		c.JSON(http.StatusBadRequest, fmt.Sprintf("Row with %s not found", email))
	} else {
		c.JSON(http.StatusOK, account)
	}
}

func (s *AccountService) FindById(c *gin.Context) {
	id := c.Param("id")
	account, err := s.accountRepository.FindById(id)

	if err != nil {
		c.JSON(http.StatusBadRequest, fmt.Sprintf("Row with %s not found", id))
	} else {
		c.JSON(http.StatusOK, account)
	}
}

func (s *AccountService) FindAll(c *gin.Context) {
	searchDto := dto.AccountSearchDto{}
	utils.BindQuery(c, &searchDto)

	account, err := s.accountRepository.FindAll(searchDto, utils.GetCurrentUserId(c))
	if err != nil {
		c.JSON(http.StatusBadRequest, errors.ErrorDescription{ErrorDescription: err.Error()})
	} else {
		c.JSON(http.StatusOK, account)
	}

}
