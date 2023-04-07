package main

import (
	"context"
	"fmt"
	"github.com/dgraph-io/dgo/v210"
	"github.com/dgraph-io/dgo/v210/protos/api"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"google.golang.org/grpc"
	"log"
	"os"
	"strconv"
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
	startDbMigrate(connection)
	repository.NewDGraphConn(connection)

	server := gin.Default()
	server.Use(security.NewSecurity().AuthMiddleware)

	resource.GeoResource(server, service.NewGeoService())
	resource.AuthResource(server, service.NewAuthService())
	resource.PostResource(server, service.NewPostService())
	resource.FriendResource(server, service.NewFriendService())
	resource.AccountResource(server, service.NewAccountService())
	resource.StorageResource(server, service.NewStorageService())
	resource.DialogResource(server, service.NewDialogService())
	resource.WebsocketResource(server, service.NewWebsocketService())

	err := server.Run()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}
}

func createDbConnection() *dgo.Dgraph {
	// Dial a gRPC connection. The address to dial to can be configured when
	// setting up the dgraph cluster.
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

func startDbMigrate(conn *dgo.Dgraph) {

	err := conn.Alter(context.Background(), &api.Operation{
		DropAll: isDropFirst(),
	})

	err = conn.Alter(context.Background(), &api.Operation{
		Schema: CreateAccountType,
	})
	if err != nil {
		log.Fatal("Alter CreateAccountType schemas has been closed with error")
	}

	err = conn.Alter(context.Background(), &api.Operation{
		Schema: CreateCaptchaType,
	})
	if err != nil {
		log.Fatal("Alter CreateCaptchaType schemas has been closed with error")
	}

	err = conn.Alter(context.Background(), &api.Operation{
		Schema: CreateCountryType,
	})
	if err != nil {
		log.Fatal("Alter CreateCountryType schemas has been closed with error")
	}

	err = conn.Alter(context.Background(), &api.Operation{
		Schema: CreatePostType,
	})
	if err != nil {
		log.Fatal("Alter CreatePostType schemas has been closed with error")
	}

	err = conn.Alter(context.Background(), &api.Operation{
		Schema: CreateCityType,
	})
	if err != nil {
		log.Fatal("Alter CreateCityType schemas has been closed with error")
	}

	err = conn.Alter(context.Background(), &api.Operation{
		Schema: CreateCommentType,
	})
	if err != nil {
		log.Fatal("Alter CreateCommentType schemas has been closed with error")
	}

	err = conn.Alter(context.Background(), &api.Operation{
		Schema: CreateTagType,
	})
	if err != nil {
		log.Fatal("Alter CreateTagType schemas has been closed with error")
	}

	err = conn.Alter(context.Background(), &api.Operation{
		Schema: CreateFriendshipType,
	})
	if err != nil {
		log.Fatal("Alter CreateFriendshipType schemas has been closed with error")
	}

	err = conn.Alter(context.Background(), &api.Operation{
		Schema: CreateLikeType,
	})
	if err != nil {
		log.Fatal("Alter CreateLikeType schemas has been closed with error")
	}

	err = conn.Alter(context.Background(), &api.Operation{
		Schema: CreateDialogType,
	})
	if err != nil {
		log.Fatal("Alter CreateDialogType schemas has been closed with error")
	}

	err = conn.Alter(context.Background(), &api.Operation{
		Schema: CreateMessageType,
	})
	if err != nil {
		log.Fatal("Alter CreateMessageType schemas has been closed with error")
	}

	if isDropFirst() {
		txn := conn.NewTxn()
		marshalR := []byte(InsertCountryRu)
		marshalRB := []byte(InsertRB)
		_, err = txn.Mutate(context.Background(), &api.Mutation{SetJson: marshalR})
		_, err = txn.Mutate(context.Background(), &api.Mutation{SetJson: marshalRB, CommitNow: true})
		if err != nil {
			log.Fatal("Import new data has been closed with error")
		}
	}

}
