package utils

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"github.com/dchest/captcha"
	"strconv"
	"web/application/domain"
	"web/application/errorhandler"
)

func CreateCaptcha() *domain.Captcha {
	captchaCode := captcha.RandomDigits(4)
	bufferCode := bytes.Buffer{}
	for _, v := range captchaCode {
		bufferCode.WriteString(strconv.Itoa(int(v)))
	}
	return domain.NewCaptchaWithCode(bufferCode.String(), captchaCode)
}

func CreateImage(captchaId string, newCaptcha *domain.Captcha) string {
	image := captcha.NewImage(captchaId, newCaptcha.CaptchaCodeByte, 150, 75)
	buffer := bytes.Buffer{}
	_, err := image.WriteTo(&buffer)
	if err != nil {
		panic(errorhandler.ErrorResponse{Message: fmt.Sprintf("Utils:CreateImage() Error write to buffer captcha %s", err)})
	}
	return "data:image/png;base64, " + base64.StdEncoding.EncodeToString(buffer.Bytes())

}
