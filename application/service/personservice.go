package service

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"web/application/domain"
	"web/application/repository"
)

type PersonService struct {
	repository *repository.DBConnection
}

func (s *PersonService) SetRepo(r *repository.DBConnection) {
	s.repository = r
}

func (s *PersonService) GetById(c *gin.Context) {
	person := domain.Person{}
	id := c.Param("id")
	domainPerson, err := s.repository.FindById(person, id)
	if err != nil {
		log.Printf(err.Error())
		c.JSON(http.StatusBadRequest, fmt.Sprintf("Row with %s not found", id))
	} else {
		c.JSON(http.StatusOK, domainPerson)
	}
}

func (s *PersonService) GetAllFields(c *gin.Context) {
	person := domain.Person{}
	personSearch := domain.PersonSearchDto{
		FirstName: c.Query("first_name"),
		LastName:  c.Query("last_name"),
	}
	spec := repository.CreateSpec().
		Like(personSearch.LastName, "LastName", personSearch).
		Or().
		Like(personSearch.FirstName, "FirstName", personSearch)

	domainPerson, err := s.repository.FindAllFields(person, spec)
	if err != nil {
		log.Printf(err.Error())
		c.JSON(http.StatusBadRequest, fmt.Sprintf("Row with %q not found", personSearch))

	} else {
		c.JSON(http.StatusOK, domainPerson)
	}
}

func (s *PersonService) GetAll(c *gin.Context) {
	person := domain.Person{}
	id := c.Param("id")
	domainPerson, err := s.repository.FindAll(person)
	if err != nil {
		log.Printf(err.Error())
		c.JSON(http.StatusBadRequest, fmt.Sprintf("Row with %s not found", id))
	} else {
		c.JSON(http.StatusOK, domainPerson)
	}
}

func (s *PersonService) Create(c *gin.Context) {
	person := domain.Person{}
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

func (s *PersonService) Update(c *gin.Context) {
	var person domain.Person
	bindJson(c, &person)
	log.Printf("Update person %v", person)
	id, err := s.repository.Update(person)
	if err != nil {
		log.Panic(c.AbortWithError(http.StatusBadRequest, err))
	} else {
		c.JSON(http.StatusCreated, &id)
	}
}

func (s *PersonService) DeleteById(c *gin.Context) {
	var person domain.Person
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
