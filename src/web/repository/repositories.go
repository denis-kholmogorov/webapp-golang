package repository

import (
	"bytes"
	"context"
	"fmt"
	"github.com/jackc/pgx"
	"log"
	"reflect"
	"strconv"
	"time"
)

const INSERT = "INSERT INTO"
const VALUES = "VALUES"
const UPDATE = "UPDATE"
const SET = "SET"
const WHERE_ID = "WHERE id ="
const DELETE = "DELETE"
const FROM = "FROM"

type DBConnection struct {
	connection *pgx.Conn
}

func (db *DBConnection) SetConnection(connection *pgx.Conn) {
	db.connection = connection
}

func (db *DBConnection) Create(domain interface{}) (interface{}, error) {
	tx, err := db.connection.BeginTx(context.Background(), pgx.TxOptions{IsoLevel: pgx.ReadCommitted})
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(context.Background())
	id := getTypeId(domain)
	hasId, err := db.existRowById(domain)
	if err != nil {
		return nil, err
	}
	if hasId {
		err = db.connection.QueryRow(context.Background(), createUpdateQuery(domain)).Scan(&id)
	} else {
		err = db.connection.QueryRow(context.Background(), createInsertQuery(domain)).Scan(&id)
	}

	if err != nil {
		log.Printf("Error execute query insert %s", err)
		return nil, err
	}
	err = tx.Commit(context.Background())
	if err != nil {
		log.Printf("Error commit transaction insert %s", err)
		return nil, err
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
		return nil, err
	}
	err = tx.Commit(context.Background())
	if err != nil {
		log.Printf("Error commit transaction update %s", err)
		return nil, err
	}
	return &id, err
}

func (db *DBConnection) DeleteById(domain interface{}, id string) (bool, error) {
	tx, err := db.connection.BeginTx(context.Background(), pgx.TxOptions{IsoLevel: pgx.ReadCommitted})
	if err != nil {
		return false, err
	}
	defer tx.Rollback(context.Background())

	tag, err := db.connection.Exec(context.Background(), createDeleteQuery(domain, id))
	log.Printf("Row affected = %d", tag.RowsAffected())
	if err != nil {
		log.Printf("Error execute query delete by id=%s, %s", id, err)
		return false, err
	}
	err = tx.Commit(context.Background())
	if err != nil {
		log.Printf("Error commit transaction update %s", err)
		return false, err
	}
	return true, nil
}

func (db *DBConnection) existRowById(s interface{}) (bool, error) {
	id, table := "", ""
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
				id = fieldToString(reflect.ValueOf(value.Interface()), field)
			}
		}
	}
	query := fmt.Sprintf("select exists(select 1 from %s where id=%s)", table, id)
	rows, err := db.connection.Query(context.Background(), query)
	defer rows.Close()

	if err != nil {
		return false, err
	}
	var isExist bool
	for rows.Next() {
		err := rows.Scan(&isExist)
		if err != nil {
			return false, err
		}
	}
	return isExist, nil
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
				queryValues.WriteString(fieldToString(value, field))
			} else {
				queryField.WriteString("," + field.Tag.Get("pg"))
				queryValues.WriteString("," + fieldToString(value, field))
			}
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
			if field.Name == "tableName" {
				table = field.Tag.Get("pg")
			} else if field.Tag.Get("pg") == "id" {
				id = fieldToString(reflect.ValueOf(value.Interface()), field)
			} else if queryFieldEqualVal.Len() == 0 {
				queryFieldEqualVal.WriteString(field.Tag.Get("pg") + "=" + fieldToString(value, field))
			} else {
				queryFieldEqualVal.WriteString("," + field.Tag.Get("pg") + "=" + fieldToString(value, field))
			}
		}
	}

	query = fmt.Sprintf("%s %s %s %s %s %s returning id;", UPDATE, table, SET, queryFieldEqualVal.String(), WHERE_ID, id)
	fmt.Println(query)

	return query
}

func createDeleteQuery(s interface{}, id string) string {
	tableName := ""
	tags := reflect.TypeOf(s)

	for i := 0; i < tags.NumField(); i++ {
		tag := tags.Field(i)
		_, ok := tag.Tag.Lookup("pg")
		if ok {
			if tag.Name == "tableName" {
				tableName = tag.Tag.Get("pg")
			}
		}
	}
	query := fmt.Sprintf("%s %s %s %s%s;", DELETE, FROM, tableName, WHERE_ID, id)
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
			default:
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

func fieldToString(value reflect.Value, tag reflect.StructField) string {
	switch value.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return strconv.Itoa(int(value.Int()))
	case reflect.String:
		return "'" + value.String() + "'"
	case reflect.Struct:
		if value.Type() == (reflect.TypeOf(time.Time{})) {
			t := value.Interface().(time.Time)
			return fmt.Sprintf("'%s'", t.Format(tag.Tag.Get("time_format")))
		}
	default:
		return ""
	}
	return ""
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
