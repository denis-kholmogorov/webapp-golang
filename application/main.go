package main

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
	"github.com/joho/godotenv"
	"log"
	"os"
	"web/application/repository"
	"web/application/resource"
	"web/application/security"
	"web/application/service"
)

// TODO использовать в репозиториях дженерики
type ApplicationContext struct {
}

func init() {
	// loads values from .env into the system
	if err := godotenv.Load("application/.env"); err != nil {
		log.Println("No .env file found")
	}
}

func main() {
	conn := createDbConnection()
	defer conn.Close(context.Background())
	repo := repository.NewRepository(conn)

	server := gin.Default()
	server.Use(security.NewSecurity().AuthMiddleware)

	resource.PersonResource(server, service.NewPersonService(repo))
	resource.AuthResource(server, service.NewAuthService(repo))

	err := server.Run()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}
}

func createDbConnection() *pgx.Conn {
	conn, err := pgx.Connect(context.Background(), getDBUrl())
	if err != nil {
		log.Fatal(err)
	}
	return conn
}

func getDBUrl() string {
	url, exists := os.LookupEnv("DB_URL")
	if !exists {
		url = "postgres://postgres:postgres@localhost:5432/postgres"
	}
	return url
}
