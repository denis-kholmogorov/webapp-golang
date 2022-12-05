package main

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx"
	"os"
	"web/repository"
	"web/resource"
	"web/service"
)

func main() {

	modelService := service.PersonService{}
	r := gin.Default()
	resource.CreatePaths(r, &modelService)
	//repo, err := repository.CreateConnect()
	//conn, err := pgx.Connect(context.Background(), "postgres://postgres:postgres@localhost:5432/postgres")
	//if err != nil {
	//	log.Fatal(err)
	//}

	conn, err := pgx.Connect(context.Background(), "postgres://postgres:postgres@localhost:5432/postgres")
	defer conn.Close(context.Background())
	db := repository.DBConnection{}
	db.SetConnection(conn)
	modelService.SetRepo(&db)
	err = r.Run()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}
}
