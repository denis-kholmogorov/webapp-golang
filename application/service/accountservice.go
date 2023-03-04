package service

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"web/application/domain"
	"web/application/dto"
	"web/application/errors"
	"web/application/kafka"
	"web/application/repository"
	"web/application/utils"
)

type AccountService struct {
	accountRepository *repository.AccountRepository
	kafkaSender       *kafkasender.KafkaSender
	person            domain.Account
}

func NewAccountService() *AccountService {
	return &AccountService{
		accountRepository: repository.GetAccountRepository(),
		kafkaSender:       kafkasender.NewKafkaSender(),
	}
}

//func (s *AccountService) FindById(c *gin.Context) {
//	id := c.Param("id")
//	domainPerson, err := s.accountRepository.FindById(id)
//	if err != nil {
//		log.Printf(err.Error())
//		c.JSON(http.StatusBadRequest, fmt.Sprintf("Row with %s not found", id))
//	} else {
//		c.JSON(http.StatusOK, domainPerson)
//	}
//}

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

	account, err := s.accountRepository.FindAll(searchDto)
	if err != nil {
		c.JSON(http.StatusBadRequest, errors.ErrorDescription{ErrorDescription: err.Error()})
	} else {
		c.JSON(http.StatusOK, account)
	}

}

//
//func (s *AccountService) GetAllFields(c *gin.Context) {
//	person := domain.Account{}
//	accountSearch := domain.AccountSearchDto{}
//	bindQuery(c, &accountSearch)
//	spec := postRepository.SpecBuilder().
//		Like(accountSearch.LastName, "LastName", true).
//		Or().
//		Equals(accountSearch.FirstName, "FirstName")
//
//	domainPerson, err := s.postRepository.FindAllBySpec(person, spec)
//	if err != nil {
//		log.Printf(err.Error())
//		c.JSON(http.StatusBadRequest, fmt.Sprintf("Row with %q not found", accountSearch))
//
//	} else {
//		c.JSON(http.StatusOK, domainPerson)
//	}
//}
//
//func (s *AccountService) FindAll(c *gin.Context) {
//	domainPerson, err := s.accountRepository.FindAll()
//	if err != nil {
//		log.Printf(err.Error())
//		c.JSON(http.StatusBadRequest, fmt.Sprintf("Rows with not found"))
//	} else {
//		c.JSON(http.StatusOK, domainPerson)
//	}
//}
//
//func (s *AccountService) Create(c *gin.Context) {
//	account := domain.Account{}
//	bindJson(c, &account)
//	log.Printf("Create new account %v", account)
//	id, err := s.postRepository.Create(account)
//	if err != nil {
//		log.Println(err)
//		c.AbortWithError(http.StatusBadRequest, err)
//	} else {
//		c.JSON(http.StatusCreated, &id)
//	}
//}
//
//func (s *AccountService) Update(c *gin.Context) {
//	var person domain.Account
//	bindJson(c, &person)
//	log.Printf("Update person %v", person)
//	id, err := s.postRepository.Update(person)
//	if err != nil {
//		log.Panic(c.AbortWithError(http.StatusBadRequest, err))
//	} else {
//		c.JSON(http.StatusCreated, &id)
//	}
//}
//
//func (s *AccountService) DeleteById(c *gin.Context) {
//	var person domain.Account
//	id := c.Param("id")
//	_, err := s.postRepository.DeleteById(person, id)
//	if err != nil {
//		log.Panic(c.AbortWithError(http.StatusBadRequest, err))
//	} else {
//		c.JSON(http.StatusCreated, &id)
//	}
//}
