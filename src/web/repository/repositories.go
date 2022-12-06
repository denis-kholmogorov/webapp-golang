package repository

import (
	"bytes"
	"context"
	"fmt"
	"github.com/jackc/pgx"
	"log"
	"reflect"
	"strconv"
)

const INSERT = "INSERT INTO"
const VALUES = "VALUES"
const UPDATE = "UPDATE"
const SET = "SET"
const WHERE_ID = "WHERE id ="

type DBConnection struct {
	connection *pgx.Conn
}

func (db *DBConnection) SetConnection(connection *pgx.Conn) {
	db.connection = connection
}

func (db *DBConnection) Create(domain interface{}) (interface{}, error) {
	tx, err := db.connection.BeginTx(context.Background(), pgx.TxOptions{IsoLevel: pgx.ReadCommitted})
	defer tx.Rollback(context.Background())
	id := getTypeId(domain)
	if err != nil {
		return nil, err
	}
	err = db.connection.QueryRow(context.Background(), createInsertQuery(domain)).Scan(&id)
	if err != nil {
		log.Printf("Error execute query insert %s", err)
	}
	err = tx.Commit(context.Background())
	if err != nil {
		log.Printf("Error commit transaction insert %s", err)
	}
	return &id, err
}

func (db *DBConnection) Update(domain interface{}) (interface{}, error) {
	tx, err := db.connection.BeginTx(context.Background(), pgx.TxOptions{IsoLevel: pgx.ReadCommitted})
	defer tx.Rollback(context.Background())
	id := getTypeId(domain)
	err = db.connection.QueryRow(context.Background(), createUpdateQuery(domain)).Scan(&id)
	if err != nil {
		log.Printf("Error execute query update %s", err)
	}
	err = tx.Commit(context.Background())
	if err != nil {
		log.Printf("Error commit transaction update %s", err)
	}
	return &id, err
}

func createInsertQuery(s interface{}) string {
	query, table := "", ""
	queryField := bytes.Buffer{}
	queryValues := bytes.Buffer{}
	values, tags := reflect.ValueOf(s), reflect.TypeOf(s)
	countFields := tags.NumField()
	for i := 0; i < countFields; i++ {

		value := values.Field(i)
		field := tags.Field(i)
		_, ok := field.Tag.Lookup("pg")

		if ok {
			if field.Name == "tableName" {
				table = field.Tag.Get("pg")
			} else if field.Tag.Get("pg") == "id" {
				log.Println("id" + field.Tag.Get("pg"))
			} else if queryField.Len() == 0 {
				queryField.WriteString(field.Tag.Get("pg"))
				queryValues.WriteString(fieldToString(reflect.ValueOf(value.Interface())))
			} else {
				queryField.WriteString("," + field.Tag.Get("pg"))
				queryValues.WriteString("," + fieldToString(reflect.ValueOf(value.Interface())))
			}
		} else {
			fmt.Println("(not specified)")
		}
	}

	query = fmt.Sprintf("%s %s (%s) %s (%s) returning id;", INSERT, table, queryField.String(), VALUES, queryValues.String())
	fmt.Println(query)
	return query
}

func createUpdateQuery(s interface{}) string {
	query, table, id := "", "", ""
	queryFieldEqualVal := bytes.Buffer{}
	values, tags := reflect.ValueOf(s), reflect.TypeOf(s)
	countFields := tags.NumField()
	for i := 0; i < countFields; i++ {

		value := values.Field(i)
		field := tags.Field(i)
		_, ok := field.Tag.Lookup("pg")

		if ok {
			if field.Tag.Get("pg") == "id" {
				id = fieldToString(reflect.ValueOf(value.Interface()))
			}
			if field.Name == "tableName" {
				table = field.Tag.Get("pg")
			} else if field.Tag.Get("pg") == "id" {

			} else if queryFieldEqualVal.Len() == 0 {
				queryFieldEqualVal.WriteString(field.Tag.Get("pg") + "=" + fieldToString(reflect.ValueOf(value.Interface())))
			} else {
				queryFieldEqualVal.WriteString("," + field.Tag.Get("pg") + "=" + fieldToString(reflect.ValueOf(value.Interface())))
			}
		} else {
			fmt.Println("(not specified)")
		}
	}

	query = fmt.Sprintf("%s %s %s %s %s %s returning id;", UPDATE, table, SET, queryFieldEqualVal.String(), WHERE_ID, id)
	fmt.Println(query)

	return query
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

//func CreateConnect() (*DBConnection, error) {
//	context := context.Background()
//	conn, err := pgx.Connect(context, "postgres://postgres:postgres@localhost:5432/postgres")
//	if err != nil {
//		return nil, err
//	}
//	defer conn.Close(context)
//	db := DBConnection{conn}
//	return &db, nil
//}