package main

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jackc/pgx/v5"
	"github.com/joho/godotenv"
	"log"
	"os"
	kafkaweb "web/application/kafka"
	"web/application/repository"
	"web/application/resource"
	"web/application/security"
	"web/application/service"
)

var dbConnection *pgx.Conn

func GetDBConn() *pgx.Conn {
	return dbConnection
}
func init() {
	// loads values from .env into the system
	if err := godotenv.Load("application/.env"); err != nil {
		log.Println("No .env file found")
	}
	m, err := migrate.New(
		"file://application/db/migrations",
		"postgres://postgres:postgres@localhost:5432/postgres?sslmode=disable&search_path=public")
	if err != nil {
		log.Fatal(err)
	}
	if err := m.Up(); err != nil {
		log.Fatal(err)
	}
}

func main() {
	dbConnection := createDbConnection()
	defer dbConnection.Close(context.Background())
	repository.NewRepository(dbConnection)

	server := gin.Default()
	server.Use(security.NewSecurity().AuthMiddleware)

	resource.AccountResource(server, service.NewAccountService())
	resource.AuthResource(server, service.NewAuthService())

	go kafkaweb.Consume1(context.Background())
	go kafkaweb.Consume2(context.Background())
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
