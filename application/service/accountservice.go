package service

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"web/application/domain"
	"web/application/kafka"
	"web/application/repository"
)

type AccountService struct {
	repository  *repository.Repository
	kafkaSender *kafkasender.KafkaSender
	person      domain.Account
}

func NewPersonService() *AccountService {
	return &AccountService{
		repository:  repository.GetRepository(),
		kafkaSender: kafkasender.NewKafkaSender(),
	}
}

func (s *AccountService) GetById(c *gin.Context) {
	id := c.Param("id")
	domainPerson, err := s.repository.FindById(s.person, id)
	s.kafkaSender.SendMessage(domainPerson)
	if err != nil {
		log.Printf(err.Error())
		c.JSON(http.StatusBadRequest, fmt.Sprintf("Row with %s not found", id))
	} else {
		c.JSON(http.StatusOK, domainPerson)
	}
}

func (s *AccountService) GetAllFields(c *gin.Context) {
	person := domain.Account{}
	personSearch := domain.PersonSearchDto{}
	bindQuery(c, &personSearch)
	spec := repository.SpecBuilder().
		Like(personSearch.LastName, "LastName", personSearch, true).
		Or().
		Equals(personSearch.FirstName, "FirstName", personSearch)

	domainPerson, err := s.repository.FindAllFields(person, spec)
	if err != nil {
		log.Printf(err.Error())
		c.JSON(http.StatusBadRequest, fmt.Sprintf("Row with %q not found", personSearch))

	} else {
		c.JSON(http.StatusOK, domainPerson)
	}
}

func (s *AccountService) GetAll(c *gin.Context) {
	person := domain.Account{}
	id := c.Param("id")
	domainPerson, err := s.repository.FindAll(person)
	if err != nil {
		log.Printf(err.Error())
		c.JSON(http.StatusBadRequest, fmt.Sprintf("Row with %s not found", id))
	} else {
		c.JSON(http.StatusOK, domainPerson)
	}
}

func (s *AccountService) Create(c *gin.Context) {
	person := domain.Account{}
	bindJson(c, &person)
	log.Printf("Create new person %v", person)
	id, err := s.repository.Create(person)
	if err != nil {
		log.Println(err)
		c.AbortWithError(http.StatusBadRequest, err)
	} else {
		c.JSON(http.StatusCreated, &id)
	}
}

func (s *AccountService) Update(c *gin.Context) {
	var person domain.Account
	bindJson(c, &person)
	log.Printf("Update person %v", person)
	id, err := s.repository.Update(person)
	if err != nil {
		log.Panic(c.AbortWithError(http.StatusBadRequest, err))
	} else {
		c.JSON(http.StatusCreated, &id)
	}
}

func (s *AccountService) DeleteById(c *gin.Context) {
	var person domain.Account
	id := c.Param("id")
	_, err := s.repository.DeleteById(person, id)
	if err != nil {
		log.Panic(c.AbortWithError(http.StatusBadRequest, err))
	} else {
		c.JSON(http.StatusCreated, &id)
	}
}

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
