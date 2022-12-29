package security

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"net/http"
	"os"
	"path"
	"strings"
)

func AuthMiddleware(ctx *gin.Context) {
	if ctx.Request.Header["Authorization"] != nil {
		rowToken := ctx.Request.Header["Authorization"][0]
		token, err := parseToken(rowToken)
		if err != nil {
			ctx.AbortWithError(http.StatusUnauthorized, fmt.Errorf("token expired"))
		}
		addValuesToContext(ctx, token)

	} else {
		ok, err := checkWhiteList(ctx)
		if err != nil || !ok {
			ctx.AbortWithError(http.StatusForbidden, fmt.Errorf("path %s not found in whiteList", ctx.Request.URL.Path))
		}
	}

}

func parseToken(rowToken string) (*jwt.Token, error) {
	env, _ := os.LookupEnv("SECRET_KEY")
	token, err := jwt.Parse(rowToken, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodRSA); ok {
			return nil, fmt.Errorf("there's an error with the signing method")
		}
		return []byte(env), nil
	})
	return token, err
}

func addValuesToContext(ctx *gin.Context, token *jwt.Token) {
	claims := token.Claims.(jwt.MapClaims)
	fmt.Println(claims)
	ctx.Set("id", claims["id"])
	ctx.Set("firstname", claims["firstName"])
	ctx.Set("lastname", claims["lastName"])
}

func checkWhiteList(ctx *gin.Context) (bool, error) {
	whiteListStr, _ := os.LookupEnv("WHITE_LIST")
	whiteList := strings.Split(whiteListStr, ",")
	requestPath := ctx.Request.URL.Path
	for _, s := range whiteList {
		match, err := path.Match(s, requestPath)
		if err != nil {
			return false, err
		}
		if match {
			return match, nil
		}
	}
	return false, nil
}
