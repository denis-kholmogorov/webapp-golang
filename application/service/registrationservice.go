package service

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"net/http"
	"os"
	"strconv"
	"time"
	"web/application/domain"
	"web/application/repository"
)

type AuthService struct {
	repository *repository.DBConnection
}

func (s *AuthService) SetRepo(r *repository.DBConnection) {
	s.repository = r
}

func (s *AuthService) Login(c *gin.Context) {
	login := domain.LoginDto{}
	person := domain.Person{}
	c.BindJSON(&login)
	domainPerson, err := s.repository.FindByFields(person, login.Email)
	switch {
	case err != nil:
		c.JSON(http.StatusBadRequest, fmt.Sprintf("Row with %s not found", login.Email))
	case domainPerson.(*domain.Person).Password != login.Password:
		c.JSON(http.StatusBadRequest, fmt.Sprintf("Password not matcher"))
	default:
		loginResponse, err := createJwtToken(domainPerson)
		if err != nil {
			c.JSON(http.StatusBadRequest, err)
		}
		c.JSON(http.StatusOK, loginResponse)
	}

}

func createJwtToken(p interface{}) (*domain.LoginResponse, error) {
	person := p.(*domain.Person)

	token := jwt.New(jwt.SigningMethodHS256)

	claims := token.Claims.(jwt.MapClaims)
	claims["exp"] = json.Number(strconv.FormatInt(time.Now().Add(100*time.Minute).Unix(), 10))
	claims["authorized"] = true
	claims["id"] = person.Id
	claims["age"] = person.Age
	claims["firstName"] = person.FirstName
	claims["lastName"] = person.LastName
	claims["email"] = person.Email
	claims["birthday"] = person.Birthday
	env, ok := os.LookupEnv("SECRET_KEY")
	if ok {
		signedString, err := token.SignedString([]byte(env))
		if err != nil {
			return nil, err
		}
		return &domain.LoginResponse{JwtToken: signedString}, err
	} else {
		return nil, fmt.Errorf("not value for signed key")
	}

}
