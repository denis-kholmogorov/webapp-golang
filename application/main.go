package main

import (
	"fmt"
	"github.com/dgraph-io/dgo/v210"
	"github.com/dgraph-io/dgo/v210/protos/api"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"google.golang.org/grpc"
	"log"
	"os"
	"strconv"
	"web/application/errorhandler"
	"web/application/repository"
	"web/application/resource"
	"web/application/security"
	"web/application/service"
)

func init() {
	// loads values from .env into the system
	if err := godotenv.Load("application/.env"); err != nil {
		log.Println("No .env file found")
	}

}

func main() {
	connection := createDbConnection()
	StartDbMigrate(connection)
	repository.NewDGraphConn(connection)

	server := gin.Default()
	server.Use(security.NewSecurity().AuthMiddleware)
	server.Use(errorhandler.ErrorHandler())

	resource.TagResource(server, service.NewTagService())
	resource.GeoResource(server, service.NewGeoService())
	resource.AuthResource(server, service.NewAuthService())
	resource.PostResource(server, service.NewPostService())
	resource.FriendResource(server, service.NewFriendService())
	resource.AccountResource(server, service.NewAccountService())
	resource.StorageResource(server, service.NewStorageService())
	resource.DialogResource(server, service.NewDialogService())
	resource.WebsocketResource(server, service.NewWebsocketService())
	resource.NotificationResource(server, service.NewSettingsService())

	err := server.Run()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}
}

func createDbConnection() *dgo.Dgraph {
	d, err := grpc.Dial(getDBUrl(), grpc.WithInsecure())
	if err != nil {
		log.Fatal(err)
	}

	return dgo.NewDgraphClient(
		api.NewDgraphClient(d),
	)
}

func getDBUrl() string {
	url, exists := os.LookupEnv("DB_URL")
	if !exists {
		url = "localhost:9080"
	}
	return url
}

func isDropFirst() bool {
	ok, exists := os.LookupEnv("DROP_FIRST")
	isDrop := false
	if exists {
		ok, err := strconv.ParseBool(ok)
		if err != nil {
			log.Fatal(err)
		}
		isDrop = ok
	}

	return isDrop
}
