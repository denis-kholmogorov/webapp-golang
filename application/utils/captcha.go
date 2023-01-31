package utils

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"github.com/dchest/captcha"
	"log"
	"strconv"
	"web/application/domain"
)

func CreateCaptcha() *domain.Captcha {
	captchaCode := captcha.RandomDigits(4)
	bufferCode := bytes.Buffer{}
	for _, v := range captchaCode {
		bufferCode.WriteString(strconv.Itoa(int(v)))
	}
	return domain.NewCaptchaWithCode(bufferCode.String(), captchaCode)
}

func CreateImage(captchaId string, newCaptcha *domain.Captcha) (string, error) {
	image := captcha.NewImage(captchaId, newCaptcha.CaptchaCodeByte, 150, 75)
	buffer := bytes.Buffer{}
	_, err := image.WriteTo(&buffer)
	if err != nil {
		log.Printf("utils.captcha: CreateImage() Error write to buffer captcha %s", err)
		return "", fmt.Errorf("utils.captcha: CreateImage() Error write to buffer captcha %s", err)
	}

	return "data:image/png;base64, " + base64.StdEncoding.EncodeToString(buffer.Bytes()), nil

}
