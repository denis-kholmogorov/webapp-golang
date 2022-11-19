package repository

import (
	"context"
	"fmt"
	"github.com/jackc/pgx"
	"log"
	"reflect"
)

type DBConnection struct {
	connection *pgx.Conn
}

func CreateConnect() *DBConnection {
	conn, err := pgx.Connect(context.Background(), "postgres://postgres:postgres@localhost:5432/postgres")
	if err != nil {
		log.Fatal(err)
	}

	defer func(conn *pgx.Conn, ctx context.Context) {
		err := conn.Close(ctx)
		if err != nil {
			log.Fatal("Connection db hasn't been closed")
		}
	}(conn, context.Background())

	connection := &DBConnection{conn}
	query, err := connection.connection.Query(context.Background(), "SELECT id, age, first_name FROM public.person WHERE id=2")
	if err != nil {
		return nil
	}
	log.Println(query)
	return connection
}

func (D DBConnection) GetBuId(s interface{}, id string) (interface{}, error) {
	query := ""
	table := ""

	st := reflect.TypeOf(s)
	countFields := st.NumField()
	for i := 0; i < countFields; i++ {

		field := st.Field(i)
		if alias, ok := field.Tag.Lookup("pg"); ok {
			if field.Name == "tableName" {
				table = field.Tag.Get("pg")
			} else if i == countFields-1 {
				query += alias
			} else {
				query += alias + ","
			}
		} else {
			fmt.Println("(not specified)")
		}
	}
	query = "SELECT " + query + " FROM " + table + " WHERE id=$id;"
	query = "SELECT id, age, first_name FROM public.person WHERE id=2"

	fmt.Println(query)

	rows, err := D.connection.Query(context.Background(), query)
	r := rows
	for t, v := range r.FieldDescriptions() {
		println(t, v.Name)
	}
	if err != nil {
		log.Println(err)
		return nil, err
	}
	log.Println(r)
	return s, nil
}
