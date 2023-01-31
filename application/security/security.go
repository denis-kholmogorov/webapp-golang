package security

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"log"
	"os"
	"path"
	"strconv"
	"strings"
)

type Security struct {
	whiteList []string
}

func NewSecurity() *Security {
	whiteListStr, ok := os.LookupEnv("WHITE_LIST")
	if !ok {
		log.Printf("Security: WhiteList in .env not found")
	}
	log.Printf("Security: Has been added to whiteList %s", whiteListStr)
	whiteList := strings.Split(whiteListStr, ",")
	return &Security{whiteList: whiteList}
}

func (conf *Security) AuthMiddleware(ctx *gin.Context) {
	//if !conf.hasPathInWhiteList(ctx) {

	//if ctx.Request.Header["Authorization"] != nil {
	//	rowToken := ctx.Request.Header["Authorization"][0]
	//	token, err := parseToken(rowToken)
	//	if err != nil {
	//		ctx.AbortWithError(http.StatusUnauthorized, fmt.Errorf("token expired"))
	//	}
	//	addValuesToContext(ctx, token)
	//} else {
	//	ctx.AbortWithError(http.StatusForbidden, fmt.Errorf("path %s not found in whiteList", ctx.Request.URL.Path))
	//}
	//}

}

func (conf *Security) hasPathInWhiteList(ctx *gin.Context) bool {
	requestPath := ctx.Request.URL.Path
	for _, s := range conf.whiteList {
		match, err := path.Match(s, requestPath)
		if err != nil {
			log.Printf("SecurityConfig: Match whitelist throw exception %s", err)
			return false
		}
		if match {
			return match
		}
	}
	return false
}

func parseToken(rowToken string) (*jwt.Token, error) {
	tokenClear := strings.TrimPrefix(rowToken, "Bearer ")
	env, _ := os.LookupEnv("SECRET_KEY")
	token, err := jwt.Parse(tokenClear, func(token *jwt.Token) (interface{}, error) {
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
	ctx.Set("id", strconv.FormatFloat(claims["id"].(float64), 'G', -1, 64))
	ctx.Set("firstname", claims["firstName"])
	ctx.Set("lastname", claims["lastName"])
}
