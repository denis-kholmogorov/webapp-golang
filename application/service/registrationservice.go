package service

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"os"
	"strconv"
	"time"
	"web/application/domain"
	"web/application/dto"
	"web/application/errorhandler"
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
	utils.BindJson(c, &reg)

	capt := s.captchaRepository.FindById(reg.CaptchaSecret)
	if capt.CaptchaCode != reg.CaptchaCode || time.Now().After(capt.ExpiredTime.Add(time.Minute*10)) {
		panic(errorhandler.MarshalError{Message: fmt.Sprintf("AuthService:Registration() Error with capthca code")})
	}

	accounts := s.accountRepository.FindByEmail(reg.Email)
	if len(accounts) != 0 {
		panic(errorhandler.ErrorResponse{Message: fmt.Sprintf("AuthService.Registration() Email already exists %s", reg.Email)})
	}

	hashPass, err := bcrypt.GenerateFromPassword([]byte(reg.Password), bcrypt.DefaultCost)
	if err != nil {
		panic(errorhandler.ErrorResponse{Message: fmt.Sprintf("AuthService.Registration() Generate password is fail %s", reg.Email)})
	}

	account := s.mapper.RegistrationToAccount(reg, hashPass)
	s.accountRepository.Save(account)

	c.Status(http.StatusCreated)
}

func (s *AuthService) Login(c *gin.Context) {
	login := domain.LoginDto{}
	utils.BindJson(c, &login)

	accounts := s.accountRepository.FindByEmail(login.Email)
	if len(accounts) != 1 {
		panic(errorhandler.NotFoundError{Message: fmt.Sprintf("AuthService.Login() Account not found by email %s", login.Email)})
	}

	err := bcrypt.CompareHashAndPassword([]byte(accounts[0].Password), []byte(login.Password))
	if err != nil {
		panic(errorhandler.ErrorResponse{Message: fmt.Sprintf("AuthService.Login() Password not equals  %s", login.Email)})
	}

	loginResponse := createJwtToken(accounts[0])
	c.JSON(http.StatusOK, loginResponse)
}

func (s *AuthService) Logout(c *gin.Context) {
	c.Status(http.StatusOK)
}

func (s *AuthService) GetCaptcha(c *gin.Context) {
	newCaptcha := utils.CreateCaptcha()
	captchaId := s.captchaRepository.Save(newCaptcha)
	image := utils.CreateImage(captchaId, newCaptcha)
	c.JSON(http.StatusOK, dto.CaptchaDto{Secret: captchaId, Image: image})
}

func createJwtToken(account domain.Account) *domain.LoginResponse {
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
			panic(errorhandler.ErrorResponse{Message: fmt.Sprintf("Registration.createJwtToken() Error with create token for account %s", account.Id)})
		}
		return &domain.LoginResponse{
			AccessToken:  signedString,
			RefreshToken: signedString,
		}
	} else {
		panic(errorhandler.ErrorResponse{Message: fmt.Sprintf("Registration.createJwtToken() Not value SECRET_KEY for signed key")})
	}
}
