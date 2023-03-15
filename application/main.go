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

	resource.AccountResource(server, service.NewAccountService())
	resource.AuthResource(server, service.NewAuthService())
	resource.GeoResource(server, service.NewGeoService())
	resource.StorageResource(server, service.NewStorageService())
	resource.PostResource(server, service.NewPostService())
	resource.FriendResource(server, service.NewFriendService())

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
	err = conn.Alter(context.Background(), &api.Operation{
		Schema: CreateCaptchaType,
	})

	err = conn.Alter(context.Background(), &api.Operation{
		Schema: CreateCountryType,
	})
	err = conn.Alter(context.Background(), &api.Operation{
		Schema: CreatePostType,
	})
	err = conn.Alter(context.Background(), &api.Operation{
		Schema: CreateCityType,
	})
	err = conn.Alter(context.Background(), &api.Operation{
		Schema: CreateCommentType,
	})
	err = conn.Alter(context.Background(), &api.Operation{
		Schema: CreateTagType,
	})
	err = conn.Alter(context.Background(), &api.Operation{
		Schema: CreateFriendshipType,
	})
	err = conn.Alter(context.Background(), &api.Operation{
		Schema: CreateLikeType,
	})

	if err != nil {
		log.Fatal("Alter schemas has been closed with error")
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
