package post05

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"

	_ "github.com/lib/pq"
)

type Userdata struct {
	ID          int
	Username    string
	Name        string
	Surname     string
	Description string
}

var (
	Hostname = ""
	Port     = 2345
	Username = ""
	Password = ""
	Database = ""
)

func openConnection() (*sql.DB, error) {
	conn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		Hostname, Port, Username, Password, Database)

	db, err := sql.Open("postgres", conn)
	if err != nil {
		return nil, err
	}
	return db, nil
}

func exists(username string) int {
	username = strings.ToLower(username)
	db, err := openConnection()
	if err != nil {
		return -1
	}
	defer db.Close()

	userid := -1
	statement := fmt.Sprintf(`SELECT "id" FROM "users" where username = '%s'`, username)
	rows, err := db.Query(statement)

	for rows.Next() {
		var id int
		err = rows.Scan(&id)
		if err != nil {
			return -1
		}
		userid = id
	}
	defer rows.Close()
	return userid
}

func Adduser(data Userdata) int {
	data.Username = strings.ToLower(data.Username)
	db, err := openConnection()
	if err != nil {
		return -1
	}
	defer db.Close()

	userId := exists(data.Username)
	if userId != -1 {
		fmt.Println("User already exists:", Username)
		return userId
	}

	insertStatement := `insert into "users" ("username") values ($1)`

	_, err = db.Exec(insertStatement, data.Username)

	if err != nil {
		fmt.Println(err)
		return -1
	}

	userId = exists(data.Username)
	if userId == -1 {
		fmt.Println("Failed to add user:", data.Username)
		return userId
	}

	insertStatement = `insert into "userdata" ("userid", "name", "surname", "description") values ($1, $2, $3, $4)`
	_, err = db.Exec(insertStatement, userId, data.Name, data.Surname, data.Description)
	if err != nil {
		fmt.Println("db.Exec()", err)
		return -1
	}
	return userId
}

func deleteUser(id int) error {
	db, err := openConnection()
	if err != nil {
		return err
	}

	defer db.Close()
}
