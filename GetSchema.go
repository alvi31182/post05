package main

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"os"
)

func main() {
	arguments := os.Args
	if len(arguments) != 6 {
		fmt.Println("Please provide: hostname port username password db")
		return
	}

	host := arguments[1]
	pgport := arguments[2]
	username := arguments[3]
	password := arguments[4]
	dbname := arguments[5]

	conn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", host, pgport, username, password, dbname)

	db, err := sql.Open("postgres", conn)
	if err != nil {
		fmt.Println("Open():", err)
		return
	}

	defer db.Close()

	rows, err := db.Query(`SELECT "datname" FROM "pg_database" WHERE datistemplate = false`)
	if err != nil {
		fmt.Println("Query():", err)
		return
	}

	for rows.Next() {
		var name string
		err = rows.Scan(&name)
		if err != nil {
			fmt.Println("Scan():", err)
			return
		}
		fmt.Println("*", name)
	}
	defer rows.Close()

	query := `SELECT table_name FROM information_schema.tables WHERE table_schema = 'public' ORDER BY table_name`
	rows, err = db.Query(query)
	if err != nil {
		fmt.Println("Query():", err)
		return
	}
	for rows.Next() {
		var name string
		err = rows.Scan(&name)
		if err != nil {
			fmt.Println("Scan():", err)
			return
		}
		fmt.Println("+T", name)
	}

}
