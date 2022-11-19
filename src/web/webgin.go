package main

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx"
	"log"
	"os"
	"web/resource"
	"web/service"
)

func main() {

	modelService := service.PersonService{}
	r := gin.Default()
	resource.CreatePaths(r, &modelService)
	//repo := repository.CreateConnect()
	conn, err := pgx.Connect(context.Background(), "postgres://postgres:postgres@localhost:5432/postgres")
	if err != nil {
		log.Fatal(err)
	}
	modelService.SetRepo(conn)
	err = r.Run()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}
}
