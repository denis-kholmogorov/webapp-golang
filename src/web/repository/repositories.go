package repository

import (
	"context"
	"fmt"
	"github.com/jackc/pgx"
	"log"
	"reflect"
	"strconv"
)

type DBConnection struct {
	connection *pgx.Conn
}

func (db *DBConnection) SetConnection(connection *pgx.Conn) {
	db.connection = connection
}

func CreateConnect() (*DBConnection, error) {
	context := context.Background()
	conn, err := pgx.Connect(context, "postgres://postgres:postgres@localhost:5432/postgres")
	if err != nil {
		return nil, err
	}
	defer conn.Close(context)
	db := DBConnection{conn}
	return &db, nil
}

func (db *DBConnection) Create(domain interface{}) (interface{}, error) {
	context := context.Background()
	tx, err := db.connection.BeginTx(context, pgx.TxOptions{IsoLevel: pgx.ReadCommitted})
	defer tx.Rollback(context)
	id := getTypeId(domain)
	err = db.connection.QueryRow(context, createInsertQuery(domain)).Scan(&id)
	if err != nil {
		log.Println(err)
	}
	err = tx.Commit(context)
	if err != nil {
		log.Fatal(err)
	}
	return &id, err
}

func getTypeId(domain interface{}) interface{} {
	val := reflect.ValueOf(domain)
	tags := reflect.TypeOf(domain)
	for i := 0; i < val.NumField(); i++ {
		if tags.Field(i).Tag.Get("pg") == "id" {
			switch val.Field(i).Kind() {
			case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
				return val.Field(i).Int()
			case reflect.String:
				return "'" + val.Field(i).String() + "'"
			}
		}
	}
	return ""
}

func createInsertQuery(s interface{}) string {
	query, fields, table, values := "", "", "", ""

	val := reflect.ValueOf(s)
	tags := reflect.TypeOf(s)
	for i := 0; i < val.NumField(); i++ {

		valueField := val.Field(i)
		tag := tags.Field(i)

		if _, ok := tag.Tag.Lookup("pg"); ok {
			if tag.Name == "tableName" {
				table = tag.Tag.Get("pg")
			} else if tag.Tag.Get("pg") == "id" {
				log.Println("id" + tag.Tag.Get("pg"))
			} else if i == val.NumField()-1 {
				values += fieldToString(reflect.ValueOf(valueField.Interface()))
				fields += tag.Tag.Get("pg")
			} else {
				values += fieldToString(reflect.ValueOf(valueField.Interface())) + ","
				fields += tag.Tag.Get("pg") + ","
			}
		} else {
			fmt.Println("(not specified)")
		}
	}

	query = "INSERT INTO " + table + " (" + fields + ") values (" + values + ") returning id;"
	fmt.Println(query)

	return query
}

func querySelect(s interface{}, id ...interface{}) string {
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
	if len(id) != 0 {
		query = "SELECT " + query + " FROM " + table + " WHERE id =$1;"
	} else {
		query = "SELECT " + query + " FROM " + table
	}

	fmt.Println(query)
	return query
}

func fieldToString(val reflect.Value) string {
	switch val.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return strconv.Itoa(int(val.Int()))
	case reflect.String:
		return "'" + val.String() + "'"
	}
	return "" //TODO добавить ошибку
}
