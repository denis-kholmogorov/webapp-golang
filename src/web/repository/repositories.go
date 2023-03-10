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

const insert = "INSERT INTO"
const SELECT_ = "SELECT"
const VALUES = "VALUES"
const update = "UPDATE"
const SET = "SET"
const WHERE_ID = "WHERE id ="
const DELETE = "DELETE"
const FROM = "FROM"
const PG = "pg"
const ID = "id"
const SELECT_EXISTS_ID = "select exists(select 1 from %s where id=%s)"
const TABLE_NAME = "tableName"

type DBConnection struct {
	connection *pgx.Conn
}

func (db *DBConnection) SetConnection(connection *pgx.Conn) {
	db.connection = connection
}

// FindById Find domain by id
func (db *DBConnection) FindById(domainType interface{}, id string) (interface{}, error) {
	tx, err := db.startTransaction()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(context.Background())
	rows, err := db.connection.Query(context.Background(), createSelectQuery(domainType, id))
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	domains := *fetchRows(rows, domainType)
	if len(domains) == 0 {
		return nil, nil
	}
	return domains[0], nil
}

// TODO доделать
// FindByFields Find domain by fields
func (db *DBConnection) FindByFields(domainType interface{}, fields ...interface{}) (interface{}, error) {
	tx, err := db.startTransaction()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(context.Background())
	fmt.Println(fmt.Sprintf("%s * %s public.person WHERE 'email'=%s;", SELECT_, FROM, fields[0]))
	rows, err := db.connection.Query(context.Background(), fmt.Sprintf("%s * %s public.person WHERE email = '%s';", SELECT_, FROM, fields[0]))
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	domains := *fetchRows(rows, domainType)
	if len(domains) == 0 {
		return nil, nil
	}
	return domains[0], nil
}

// FindAll Find all domain
func (db *DBConnection) FindAll(domainType interface{}) (interface{}, error) {

	tx, err := db.startTransaction()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(context.Background())
	rows, err := db.connection.Query(context.Background(), createSelectQuery(domainType))
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	domains := *fetchRows(rows, domainType)
	return domains, nil
}

// Create new domain
func (db *DBConnection) Create(domain interface{}) (interface{}, error) {
	tx, err := db.startTransaction()
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

// Update domain
func (db *DBConnection) Update(domain interface{}) (interface{}, error) {
	tx, err := db.startTransaction()
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

// DeleteById delete domain by id
func (db *DBConnection) DeleteById(domain interface{}, id string) (bool, error) {
	tx, err := db.startTransaction()
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

func (db *DBConnection) existRowById(domain interface{}) (bool, error) {
	id, table := "", ""
	values, tags := reflect.ValueOf(domain), reflect.TypeOf(domain)
	countFields := tags.NumField()
	for i := 0; i < countFields; i++ {

		value := values.Field(i)
		field := tags.Field(i)
		_, ok := field.Tag.Lookup(PG)
		if ok {
			if field.Name == TABLE_NAME {
				table = field.Tag.Get(PG)
			} else if field.Tag.Get(PG) == ID {
				id = fieldToString(reflect.ValueOf(value.Interface()), field)
			}
		}
	}
	query := fmt.Sprintf(SELECT_EXISTS_ID, table, id)
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
	queryField, queryValues := bytes.Buffer{}, bytes.Buffer{}
	values, tags := reflect.ValueOf(s), reflect.TypeOf(s)
	countFields := tags.NumField()
	for i := 0; i < countFields; i++ {
		value := values.Field(i)
		field := tags.Field(i)
		if _, ok := field.Tag.Lookup(PG); ok {
			switch {
			case field.Name == TABLE_NAME:
				table = field.Tag.Get(PG)
			case field.Tag.Get(PG) == ID:
				log.Println(ID + field.Tag.Get(PG))
			case queryField.Len() == 0:
				queryField.WriteString(field.Tag.Get(PG))
				queryValues.WriteString(fieldToString(value, field))
			default:
				queryField.WriteString("," + field.Tag.Get(PG))
				queryValues.WriteString("," + fieldToString(value, field))
			}
		}
	}

	query = fmt.Sprintf("%s %s (%s) %s (%s) returning id;", insert, table, queryField.String(), VALUES, queryValues.String())
	fmt.Println(query)
	return query
}

func createUpdateQuery(s interface{}) string {
	query, table, id := "", "", ""
	queryFieldEqualVal := bytes.Buffer{}
	values, tags := reflect.ValueOf(s), reflect.TypeOf(s)
	countFields := tags.NumField()
	for i := 0; i < countFields; i++ {

		value, field := values.Field(i), tags.Field(i)
		if _, ok := field.Tag.Lookup(PG); ok {
			switch {
			case field.Name == TABLE_NAME:
				table = field.Tag.Get(PG)
			case field.Tag.Get(PG) == ID:
				id = fieldToString(reflect.ValueOf(value.Interface()), field)
			case queryFieldEqualVal.Len() == 0:
				queryFieldEqualVal.WriteString(field.Tag.Get(PG) + "=" + fieldToString(value, field))
			default:
				queryFieldEqualVal.WriteString("," + field.Tag.Get(PG) + "=" + fieldToString(value, field))

			}
		}
	}
	query = fmt.Sprintf("%s %s %s %s %s %s returning id;", update, table, SET, queryFieldEqualVal.String(), WHERE_ID, id)
	fmt.Println(query)
	return query
}

func createDeleteQuery(s interface{}, id string) string {
	tableName := ""
	tags := reflect.TypeOf(s)

	for i := 0; i < tags.NumField(); i++ {
		tag := tags.Field(i)
		if _, ok := tag.Tag.Lookup(PG); ok && tag.Name == TABLE_NAME {
			tableName = tag.Tag.Get(PG)
		}
	}
	query := fmt.Sprintf("%s %s %s %s%s;", DELETE, FROM, tableName, WHERE_ID, id)
	fmt.Println(query)
	return query
}

func getTypeId(domain interface{}) interface{} {
	val, tags := reflect.ValueOf(domain), reflect.TypeOf(domain)
	for i := 0; i < val.NumField(); i++ {
		if tags.Field(i).Tag.Get(PG) == ID {
			switch val.Field(i).Kind() {
			case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
				return val.Field(i).Int()
			case reflect.String:
				return fmt.Sprintf("'%s'", val.Field(i).String())
			default:
				return fmt.Sprintf("'%s'", val.Field(i).String())
			}
		}
	}
	return ""
}

func createSelectQuery(s interface{}, id ...string) string {
	query, table := "", ""
	queryField := bytes.Buffer{}
	tags := reflect.TypeOf(s)
	countFields := tags.NumField()
	for i := 0; i < countFields; i++ {
		field := tags.Field(i)
		if _, ok := field.Tag.Lookup(PG); ok {
			switch {
			case field.Name == TABLE_NAME:
				table = field.Tag.Get(PG)
			case queryField.Len() == 0:
				queryField.WriteString(field.Tag.Get(PG))
			default:
				queryField.WriteString("," + field.Tag.Get(PG))
			}
		}
	}

	if len(id) == 1 {
		query = fmt.Sprintf("%s %s %s %s %s %s;", SELECT_, queryField.String(), FROM, table, WHERE_ID, id[0])
	} else {
		query = fmt.Sprintf("%s %s %s %s;", SELECT_, queryField.String(), FROM, table)
	}
	fmt.Println(query)
	return query
}

func fieldToString(value reflect.Value, tag reflect.StructField) string {
	switch value.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return strconv.Itoa(int(value.Int()))
	case reflect.String:
		return fmt.Sprintf("'%s'", value.String())
	case reflect.Struct:
		switch value.Type() {
		case reflect.TypeOf(time.Time{}):
			return fmt.Sprintf("'%s'", value.Interface().(time.Time).Format(tag.Tag.Get("time_format")))
		}
	default:
		return ""
	}
	return ""
}

func fetchRows(rows pgx.Rows, domainType interface{}) *[]interface{} {
	var domains []interface{}
	for rows.Next() {
		data, _ := rows.Values()
		fields := rows.FieldDescriptions()
		domain := reflect.New(reflect.TypeOf(domainType))
		for i, field := range fields {
			typeDomain := domain.Type().Elem()
			for l := 0; l < typeDomain.NumField(); l++ {
				if typeDomain.Field(l).Tag.Get(PG) == field.Name && data[i] != nil {
					field := domain.Elem().Field(l)
					setValue(field, data[i])
					break
				}
			}
		}
		domains = append(domains, domain.Interface())
	}
	return &domains
}

func setValue(field reflect.Value, value any) {
	switch reflect.ValueOf(value).Kind() {
	case reflect.Int64:
		field.SetInt(value.(int64))
	case reflect.Int32:
		field.SetInt(int64(value.(int32)))
	case reflect.Int16:
		field.SetInt(int64(value.(int16)))
	case reflect.Int8:
		field.SetInt(int64(value.(int8)))
	case reflect.Float32:
		field.SetFloat(float64(value.(float32)))
	case reflect.Float64:
		field.SetFloat(value.(float64))
	case reflect.Uintptr:
		field.SetUint(value.(uint64))
	case reflect.Int:
		field.SetInt(int64(value.(int)))
	default:
		field.Set(reflect.ValueOf(value))
	}
}

func (db *DBConnection) startTransaction() (pgx.Tx, error) {
	tx, err := db.connection.BeginTx(context.Background(), pgx.TxOptions{IsoLevel: pgx.ReadCommitted})
	if err != nil {
		return nil, err
	}
	return tx, nil
}
