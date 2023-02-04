package service

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"web/application/domain"
	"web/application/kafka"
	"web/application/repository"
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
	id, _ := c.Get("id")
	account, err := s.accountRepository.FindById(id.(string))

	if err != nil {
		c.JSON(http.StatusBadRequest, fmt.Sprintf("Row with %s not found", id))
	} else {
		c.JSON(http.StatusOK, account)
	}
}

//
//func (s *AccountService) GetAllFields(c *gin.Context) {
//	person := domain.Account{}
//	accountSearch := domain.AccountSearchDto{}
//	bindQuery(c, &accountSearch)
//	spec := repository.SpecBuilder().
//		Like(accountSearch.LastName, "LastName", true).
//		Or().
//		Equals(accountSearch.FirstName, "FirstName")
//
//	domainPerson, err := s.repository.FindAllBySpec(person, spec)
//	if err != nil {
//		log.Printf(err.Error())
//		c.JSON(http.StatusBadRequest, fmt.Sprintf("Row with %q not found", accountSearch))
//
//	} else {
//		c.JSON(http.StatusOK, domainPerson)
//	}
//}
//
//func (s *AccountService) GetAll(c *gin.Context) {
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
//	id, err := s.repository.Create(account)
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
//	id, err := s.repository.Update(person)
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
//	_, err := s.repository.DeleteById(person, id)
//	if err != nil {
//		log.Panic(c.AbortWithError(http.StatusBadRequest, err))
//	} else {
//		c.JSON(http.StatusCreated, &id)
//	}
//}

func bindJson(c *gin.Context, value any) {
	err := c.BindJSON(value)
	if err != nil {
		fmt.Printf("Can't parse json from request %s", err.Error())
		c.AbortWithStatus(http.StatusBadRequest)
	}
}

func bindQuery(c *gin.Context, value any) {
	err := c.BindQuery(value)
	if err != nil {
		fmt.Printf(err.Error())
		c.AbortWithStatus(http.StatusBadRequest)
	}
}
