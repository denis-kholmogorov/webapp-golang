package service

import (
	"context"
	"fmt"
	"github.com/georgysavva/scany/pgxscan"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx"
	"log"
	"net/http"
	"reflect"
	"strconv"
)

type Person struct {
	tableName struct{} `pg:"public.person"`
	Id        int      `json:"id" pg:"id"`
	Age       int      `json:"age" pg:"age"`
	FirstName string   `json:"firstName" pg:"first_name"`
}

type PersonService struct {
	repository *pgx.Conn
}

func (s *PersonService) SetRepo(r *pgx.Conn) {
	s.repository = r
}

func (s *PersonService) GetById(c *gin.Context) {
	person := Person{}
	id := c.Param("id")
	err := pgxscan.Get(context.Background(), s.repository, &person, createQuery(person, id), id)

	if err != nil {
		log.Printf(err.Error())
		c.JSON(http.StatusBadRequest, fmt.Sprintf("Row with %s not found", id))
	} else {
		c.JSON(http.StatusOK, person)
	}
}

func (s *PersonService) GetAll(c *gin.Context) {
	var persons []*Person
	err := pgxscan.Select(context.Background(), s.repository, &persons, createQuery(Person{}))
	if err != nil {
		log.Println("Rows not found")
		c.JSON(http.StatusBadRequest, fmt.Sprintln("Rows not found"))
	} else {
		c.JSON(http.StatusOK, persons)
	}
}

func (s *PersonService) Create(c *gin.Context) {
	background := context.Background()
	tx, err := s.repository.BeginTx(background, pgx.TxOptions{IsoLevel: pgx.Serializable})
	if err != nil {
		log.Panic(err)
	}
	var person Person
	var id int
	err = c.BindJSON(&person)
	defer tx.Rollback(background)
	err = s.repository.QueryRow(background, createInsert(person)).Scan(&id)
	if err != nil {
		log.Println(err)
	}
	err = tx.Commit(background)
	if err != nil {
		log.Fatal(err)
	}
	c.JSON(http.StatusOK, id)
}

func createQuery(s interface{}, id ...interface{}) string {
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

func createInsert(s interface{}) string {
	values := ""
	table := ""
	fields := ""
	query := ""

	val := reflect.ValueOf(s)
	tags := reflect.TypeOf(s)
	for i := 0; i < val.NumField(); i++ {

		valueField := val.Field(i)
		tag := tags.Field(i)

		if _, ok := tag.Tag.Lookup("pg"); ok {
			if tag.Name == "tableName" {
				table = tag.Tag.Get("pg")
			} else if tag.Name == "Id" {
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

func fieldToString(val reflect.Value) string {
	switch val.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return strconv.Itoa(int(val.Int()))
	case reflect.String:
		return "'" + val.String() + "'"
	}
	return ""
}

//func (s *PersonService)GetById(c *gin.Context)  {
//	person := Person{}
//	//id := c.Param("id")
//	//err := pgxscan.Get(context.Background(), s.repository, &person, "select id, age, first_name from public.person where id = $1", id)
//	row, err := s.repository.Query(context.Background(),"select id, age, first_name from public.person where id = 2")
//	if err != nil {
//		log.Printf(err.Error())
//		//c.JSON(http.StatusBadRequest, fmt.Sprintf("Row with %s not found", id))
//	}
//	if row != nil {
//		for row.Next() {
//			err := row.Scan(&person.Id,&person.Age,&person.FirstName)
//			if err != nil {
//				fmt.Println(err)
//			}
//		}
//	}
//	fmt.Println(person)
//}

//
//type PersonService struct {
//repository *repository.DBConnection
//}
//
//func (s *PersonService) SetRepo(r *repository.DBConnection)  {
//
//	s.repository = r
//}
//
//func (s *PersonService)GetById(c *gin.Context)  {
//	person := Person{}
//	id := c.Param("id")
//	_, err := s.repository.GetBuId(person,c.Param("id"))
//
//	if err != nil {
//		log.Printf(err.Error())
//		c.JSON(http.StatusBadRequest, fmt.Sprintf("Row with %s not found", id))
//	} else {
//		c.JSON(http.StatusOK, person)
//	}
//}
//
//func (s *PersonService)GetAll(c *gin.Context)  {
//	var persons []*Person
//	err := pgxscan.Select(context.Background(), s.conn, &persons, "select id, age, first_name from public.person")
//	if err != nil {
//		log.Println("Rows not found")
//		c.JSON(http.StatusBadRequest, fmt.Sprintln("Rows not found"))
//	} else {
//		c.JSON(http.StatusOK, persons)
//	}
//}
//
//
//func (s *PersonService)Create(c *gin.Context) {
//	tx, err := s.conn.BeginTx(c, pgx.TxOptions{IsoLevel: pgx.Serializable})
//	var person Person
//	err = c.BindJSON(&person)
//	rows := [][]interface{}{
//		{person.Age, person.FirstName},
//	}
//
//	from, err := s.conn.CopyFrom(c, pgx.Identifier{"person"}, []string{"age", "first_name"}, pgx.CopyFromRows(rows))
//	if err != nil {
//		log.Println(err)
//		tx.Rollback(c)
//	} else {
//		log.Println(from)
//		tx.Commit(c)
//		c.JSON(http.StatusCreated,nil)
//	}
//}
//
//func (s *PersonService)Test(c *gin.Context) {
//	//var name string
//	//var age int64
//
//	var persons []*Person
//	pgxscan.Select(context.Background(), s.conn, &persons, "select id, age, first_name from public.person")
//	//err := s.conn.QueryRow(context.Background(), "select age, first_name from public.person").Scan()
//	//if err != nil {
//	//	fmt.Fprintf(os.Stderr, "QueryRow failed: %v\n", err)
//	//	os.Exit(1)
//	//}
//
//	for p := range persons {
//	fmt.Println(p)
//
//}
//
//
//}
//
//func (s *PersonService) Update(c *gin.Context) {
//
//}
//
//func (s *PersonService) Delete(c *gin.Context) {
//}
//
//rows := [][]interface{}{
//{person.Age, person.FirstName},
//}
//defer tx.Rollback(background)
//from, err := s.repository.CopyFrom(c, pgx.Identifier{"person"}, []string{"age", "first_name"}, pgx.CopyFromRows(rows))
//if err != nil {
//log.Println(err)
//tx.Rollback(background)
//} else {
//log.Println(from)
//tx.Commit(background)
//c.JSON(http.StatusCreated, nil)
//}
//
//batch := &pgx.Batch{}
//batch.Queue("insert into ledger(description, amount) values($1, $2)", "q1", 1)
//batch.Queue("insert into ledger(description, amount) values($1, $2)", "q2", 2)
//br := s.repository.SendBatch(context.Background(), batch)
//log.Println(br.QueryRow())
