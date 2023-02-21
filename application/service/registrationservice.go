package service

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"golang.org/x/crypto/bcrypt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"
	"web/application/domain"
	"web/application/dto"
	"web/application/mapper"
	"web/application/repository"
	"web/application/utils"
)

type AuthService struct {
	accountRepository *repository.AccountRepository
	captchaRepository *repository.CaptchaRepository
	mapper            *mapper.PersonMapper
}

func NewAuthService() *AuthService {

	return &AuthService{
		captchaRepository: repository.GetCaptchaRepository(),
		accountRepository: repository.GetAccountRepository(),
		mapper:            mapper.GetPersonMapper(),
	}
}

func (s *AuthService) Registration(c *gin.Context) {
	reg := dto.RegistrationDto{}
	err := c.BindJSON(&reg)
	if err != nil {
		log.Fatal(err)
	}

	capt, err := s.captchaRepository.FindById(reg.CaptchaSecret)
	if err != nil || capt.CaptchaCode != reg.CaptchaCode || time.Now().After(capt.ExpiredTime.Add(time.Minute*10)) {
		c.AbortWithStatusJSON(http.StatusBadRequest, fmt.Sprintf("Registration incorrect"))
		return
	}

	accounts, err := s.accountRepository.FindByEmail(reg.Email)
	if len(accounts) != 0 || err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, fmt.Sprintf("Email already exists"))
		return
	}

	hashPass, err := bcrypt.GenerateFromPassword([]byte(reg.Password), bcrypt.DefaultCost)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, fmt.Sprintf("Password incorrect"))
		return
	}
	account := s.mapper.RegistrationToAccount(reg, hashPass)
	err = s.accountRepository.Save(account)
	switch {
	case err != nil:
		c.JSON(http.StatusBadRequest, fmt.Sprintf("Row with not found"))
	default:
		c.Status(http.StatusCreated)
	}
}

func (s *AuthService) Login(c *gin.Context) {
	login := domain.LoginDto{}
	c.BindJSON(&login)
	accounts, err := s.accountRepository.FindByEmail(login.Email)
	if len(accounts) != 1 {
		c.AbortWithStatusJSON(http.StatusBadRequest, fmt.Sprintf("Login: has been found wore then one account by this email"))
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(accounts[0].Password), []byte(login.Password))
	switch {
	case err != nil:
		c.JSON(http.StatusBadRequest, fmt.Sprintf("Error %s", err))
	default:
		loginResponse, err := createJwtToken(accounts[0])
		if err != nil {
			c.JSON(http.StatusBadRequest, err)
		}
		c.JSON(http.StatusOK, loginResponse)
	}

}

func (s *AuthService) Logout(c *gin.Context) {
	c.JSON(http.StatusOK, "")
}

func (s *AuthService) GetCaptcha(c *gin.Context) {
	newCaptcha := utils.CreateCaptcha()
	captchaId, err := s.captchaRepository.Save(newCaptcha)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, fmt.Sprintf("getCaptcha failed"))
		return
	}

	image, err := utils.CreateImage(*captchaId, newCaptcha)

	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, fmt.Sprintf("getCaptcha failed"))
		return
	}
	c.JSON(http.StatusOK, dto.CaptchaDto{Secret: *captchaId, Image: image})
}

func createJwtToken(account domain.Account) (*domain.LoginResponse, error) {
	token := jwt.New(jwt.SigningMethodHS256)

	claims := token.Claims.(jwt.MapClaims)
	claims["exp"] = json.Number(strconv.FormatInt(time.Now().Add(100*time.Minute).Unix(), 10))
	claims["authorized"] = true
	claims["id"] = account.Id
	claims["age"] = account.Age
	claims["firstName"] = account.FirstName
	claims["lastName"] = account.LastName
	claims["email"] = account.Email
	claims["birthday"] = account.BirthDate
	env, ok := os.LookupEnv("SECRET_KEY")
	if ok {
		signedString, err := token.SignedString([]byte(env))
		if err != nil {
			return nil, err
		}
		return &domain.LoginResponse{
			AccessToken:  signedString,
			RefreshToken: signedString,
		}, err
	} else {
		return nil, fmt.Errorf("not value for signed key")
	}

}
