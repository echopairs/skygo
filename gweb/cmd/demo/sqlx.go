package main

import (
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"log"
)

var schema = `
CREATE TABLE person (
	first_name 	text,
	last_name 	text,
	email 		text
)
`

type Person struct {
	FirstName string `db:"first_name"`
	LastName  string `db:"last_name"`
	Email     string
}

func main() {
	source := "root:admin@tcp(127.0.0.1:3306)/pairs"
	db, err := sqlx.Connect("mysql", source)
	if err != nil {
		log.Fatal(err)
	}
	//db.MustExec(schema)
	// 1. Named queries
	//db.NamedExec("INSERT INTO person (first_name, last_name, email) VALUES (:first_name, :last_name, :email)", &Person{"Jane", "Citizen", "jane.citzen@example.com"})

	// 2. Select get slice
	people := []Person{}
	err = db.Select(&people, "select *from person where first_name = ?", "Jane")
	if err != nil {
		log.Fatal(err)
	}
	for _, v := range people {
		fmt.Println(v)
	}

	// 3. Get one
	z := Person{}
	err = db.Get(&z, "select first_name from person where first_name =?", "zheng")
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(z)

	// 4. 数据量大的时候，为了减少内存分配使用rows.Next
	p := Person{}
	rows, err := db.Queryx("select * from person")
	if err != nil {
		fmt.Println(err)
	}
	for rows.Next() {
		err := rows.StructScan(&p)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(p)
	}

	// 5. name 綁定插入 (map映射/結構體映射)
	_, err = db.NamedExec(`Insert INTO person (first_name, last_name, email) VALUES (:first, :last, :email)`,
		map[string]interface{}{
			"first": "Bin",
			"last":  "Smuth",
			"email": "sss@qq.com",
		})

	// 6. Select 條件在 結構體中/map中
	_, err = db.NamedQuery(`select *from person Where first_name=:fn`, map[string]interface{}{
		"fn": "Bin"})

	row, err := db.NamedQuery("select *from person where first_name=:first_name", p)
	if err != nil {
		fmt.Println(err)
	} else {
		for row.Next() {
			pp := Person{}
			row.StructScan(&pp)
			fmt.Println("pp: ", pp)
		}
	}
}
