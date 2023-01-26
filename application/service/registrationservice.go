package service

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/dchest/captcha"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"os"
	"strconv"
	"time"
	"web/application/domain"
	"web/application/dto"
	"web/application/mapper"
	"web/application/repository"
)

type AuthService struct {
	captchaRepo  *repository.CaptchaRepository[*domain.Captcha, string]
	accountRepo  *repository.AccountRepository[*domain.Account, string]
	repository   *repository.Repository
	mapper       *mapper.PersonMapper
	captchaStore *captcha.Store
}

func NewAuthService() *AuthService {
	store := captcha.NewMemoryStore(100, time.Minute*2)
	return &AuthService{
		captchaRepo:  repository.GetCaptchaRepository(),
		accountRepo:  repository.GetAccountRepository(),
		repository:   repository.GetRepository(),
		mapper:       mapper.GetPersonMapper(),
		captchaStore: &store,
	}
}

func (s *AuthService) Registration(c *gin.Context) {
	reg := dto.RegistrationDto{}
	c.BindJSON(&reg)

	capd, err := s.captchaRepo.FindById(&domain.Captcha{}, reg.CaptchaSecret)
	if err != nil || capd.CaptchaCode != reg.CaptchaCode || time.Now().After(capd.ExpiredTime.Add(time.Minute*10)) {
		c.AbortWithStatusJSON(http.StatusBadRequest, fmt.Sprintf("Captcha failed"))
		return
	}
	hashPass, err := bcrypt.GenerateFromPassword([]byte(reg.Password), bcrypt.DefaultCost)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, fmt.Sprintf("Password incorrect"))
		return
	}
	account := s.mapper.RegistrationToAccount(reg, hashPass)
	entity, err := s.repository.Create(account)
	switch {
	case err != nil:
		c.JSON(http.StatusBadRequest, fmt.Sprintf("Row with not found"))
	default:
		c.JSON(http.StatusCreated, entity)
	}
}

func (s *AuthService) Login(c *gin.Context) {
	login := domain.LoginDto{}
	c.BindJSON(&login)
	spec := repository.SpecBuilder().Equals(login.Email, "email")
	accounts, err := s.accountRepo.FindAllBySpec(spec)
	if len(accounts) != 1 {
		c.AbortWithStatusJSON(http.StatusBadRequest, "Account not found")
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

func (s *AuthService) Captcha(c *gin.Context) {

	captchaCode := captcha.RandomDigits(4)
	bufferCode := bytes.Buffer{}
	for _, v := range captchaCode {
		bufferCode.WriteString(strconv.Itoa(int(v)))
	}

	domain := domain.Captcha{
		CaptchaCode: bufferCode.String(),
		ExpiredTime: time.Now(),
	}
	id, err := s.repository.Create(domain)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, fmt.Sprintf("Captcha failed"))
	}
	ids := fmt.Sprint(id)
	image := captcha.NewImage(ids, captchaCode, 150, 75)
	buffer := bytes.Buffer{}
	image.WriteTo(&buffer)
	s3 := "data:image/png;base64, " + base64.StdEncoding.EncodeToString(buffer.Bytes())
	c.JSON(http.StatusOK, dto.CaptchaDto{Secret: ids, Image: s3})
}

func (s *AuthService) Logout(c *gin.Context) {
	c.JSON(http.StatusOK, "")
}

func createJwtToken(account *domain.Account) (*domain.LoginResponse, error) {
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
