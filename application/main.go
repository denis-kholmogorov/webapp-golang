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

func startDbMigrate(conn *dgo.Dgraph) {

	err := conn.Alter(context.Background(), &api.Operation{
		DropAll: false,
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
		Schema: CreateCityType,
	})

	txn := conn.NewTxn()

	//country := domain.Country{
	//	DType: []string{"Country"},
	//	Title: "Russia",
	//	Cities: []domain.City{
	//		domain.City{
	//			DType: []string{"City"},
	//			Title: "Moscow",
	//		},
	//	},
	//}
	//marshal, err := json.Marshal(country)
	marshalR := []byte(InsertCountryRu)
	marshalRB := []byte(InsertRB)
	if err != nil {
		return
	}
	_, err = txn.Mutate(context.Background(), &api.Mutation{SetJson: marshalR})
	_, err = txn.Mutate(context.Background(), &api.Mutation{SetJson: marshalRB, CommitNow: true})
	if err != nil {
		return
	}

	if err != nil {
		log.Fatal(err)
	}

	if err != nil {
		log.Fatal("The drop all operation should have succeeded")
	}
}
