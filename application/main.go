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

func init() {
	// loads values from .env into the system
	if err := godotenv.Load(); err != nil {
		log.Print("No .env file found")
	}
}

func main() {

	personService := service.PersonService{}
	authService := service.AuthService{}
	server := gin.Default()
	server.Use(security.AuthMiddleware)
	resource.PersonResource(server, &personService)
	resource.AuthResource(server, &authService)
	conn, err := pgx.Connect(context.Background(), getDBUrl())
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close(context.Background())
	db := repository.DBConnection{}
	db.SetConnection(conn)
	personService.SetRepo(&db)
	authService.SetRepo(&db)
	err = server.Run()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}
}

func getDBUrl() string {
	url, exists := os.LookupEnv("DB_URL")
	if !exists {
		url = "postgres://postgres:postgres@localhost:5432/postgres"
	}
	return url
}
