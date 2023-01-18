package configuration

import (
	"log"
	"os"
	"path"
	"strings"
)

type SecurityConfig struct {
	whiteList []string
}

func NewSecurityConfig() *SecurityConfig {
	whiteListStr, ok := os.LookupEnv("WHITE_LIST")
	if !ok {
		log.Printf("Security: WhiteList in .env not found")
	}
	log.Printf("Security: Has been added to whiteList %s", whiteListStr)
	whiteList := strings.Split(whiteListStr, ",")
	return &SecurityConfig{whiteList: whiteList}
}

func (conf *SecurityConfig) checkIsWhitePath(requestPath string) bool {
	for _, s := range conf.whiteList {
		match, err := path.Match(s, requestPath)
		if err != nil {
			log.Printf("SecurityConfig: Match whitelist throw exception %s", err)
			return false
		}
		return match
	}
	return false
}
