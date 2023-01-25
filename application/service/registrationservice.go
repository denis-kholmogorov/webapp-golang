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
	"reflect"
	"strconv"
	"time"
	"web/application/domain"
	"web/application/dto"
	"web/application/mapper"
	"web/application/repository"
)

type AuthService struct {
	repository   *repository.Repository
	mapper       *mapper.PersonMapper
	captchaStore *captcha.Store
}

func NewAuthService() *AuthService {
	store := captcha.NewMemoryStore(100, time.Minute*2)
	return &AuthService{
		repository:   repository.GetRepository(),
		mapper:       mapper.GetPersonMapper(),
		captchaStore: &store,
	}
}

func (s *AuthService) Registration(c *gin.Context) {
	reg := dto.RegistrationDto{}
	c.BindJSON(&reg)

	cap, err := s.repository.FindById(domain.Captcha{}, reg.CaptchaSecret)
	if err != nil || cap == nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, fmt.Sprintf("Captcha failed"))
	}

	capa := reflect.ValueOf(cap).Interface()
	capd := capa.(*domain.Captcha)
	if capd.CaptchaCode != reg.CaptchaCode || time.Now().After(capd.ExpiredTime.Add(time.Minute*10)) {
		c.AbortWithStatusJSON(http.StatusBadRequest, fmt.Sprintf("Captcha failed"))
	}

	hashPass, err := bcrypt.GenerateFromPassword([]byte(reg.Password), bcrypt.DefaultCost)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, fmt.Sprintf("Password incorrect"))
		return
	}
	account := s.mapper.RegistrationToAccount(reg, hashPass)
	//domain, err := s.repository.Create(account)
	switch {
	case err != nil:
		c.JSON(http.StatusBadRequest, fmt.Sprintf("Row with not found"))
	default:
		//c.JSON(http.StatusCreated, domain)
		c.JSON(http.StatusCreated, account)
	}
}

func (s *AuthService) Login(c *gin.Context) {
	login := domain.LoginDto{}
	person := domain.Account{}
	c.BindJSON(&login)
	domainPerson, err := s.repository.FindByFields(person, login.Email)
	switch {
	case err != nil:
		c.JSON(http.StatusBadRequest, fmt.Sprintf("Row with %s not found", login.Email))
	case domainPerson.(*domain.Account).Password != login.Password:
		c.JSON(http.StatusBadRequest, fmt.Sprintf("Password not matcher"))
	default:
		loginResponse, err := createJwtToken(domainPerson)
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

func createJwtToken(p interface{}) (*domain.LoginResponse, error) {
	person := p.(*domain.Account)

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
