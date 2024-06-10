/*
The package works on 2 tables on a PostgreSQL data base server.
The names of the tables are:
  - Users
  - Userdata

The definitions of the tables in the PostgreSQL server are:

	CREATE TABLE Users (
	ID SERIAL,
	Username VARCHAR(100) PRIMARY KEY
	);
	CREATE TABLE Userdata (
	UserID Int NOT NULL,
	Name VARCHAR(100),
	Surname VARCHAR(100),
	Description VARCHAR(200)
	);

This is rendered as code
This is not rendered as code
*/
package post05

import (
	"context"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v5"
	"strings"
)

type Userdata struct {
	ID          int
	Username    string
	Name        string
	Surname     string
	Description string
}

/*
This block of global variables holds the connection details
to the Postgres server
Hostname: is the IP or the hostname of the server
Port: is the TCP port the DB server listens to
Username: is the username of the database user
Password: is the password of the database user
Database: is the name of the Database in PostgreSQL
*/
var (
	Hostname = ""
	Port     = 2345
	Username = ""
	Password = ""
	Database = ""
)

func openConnection() (*pgx.Conn, error) {
	config, err := pgx.ParseConfig(fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		Hostname, Port, Username, Password, Database))
	if err != nil {
		return nil, err
	}

	config.RuntimeParams = map[string]string{
		"statement_timeout": "30000",
	}

	conn, err := pgx.ConnectConfig(context.Background(), config)
	if err != nil {
		return nil, err
	}
	return conn, nil
}

func exists(username string) int {
	username = strings.ToLower(username)
	conn, err := openConnection()
	if err != nil {
		return -1
	}
	defer conn.Close(context.Background())

	var userid int
	err = conn.QueryRow(
		context.Background(), `SELECT "id" FROM "users" where username = $1`, username).Scan(&userid)
	if err != nil {
		return -1
	}
	return userid
}

/*
Add new Userdata struct
*/
func Adduser(data Userdata) int {
	data.Username = strings.ToLower(data.Username)
	db, err := openConnection()
	if err != nil {
		fmt.Println("Adding error")
		return -1
	}
	defer db.Close(context.Background())

	userId := exists(data.Username)
	if userId != -1 {
		fmt.Println("User already exists:", Username)
		return userId
	}

	_, err = db.Exec(context.Background(),
		`insert into "users" ("username") values ($1)`, data.Username)

	if err != nil {
		fmt.Println(err)
		return -1
	}

	userId = exists(data.Username)
	if userId == -1 {
		fmt.Println("Failed to add user:", data.Username)
		return userId
	}

	_, err = db.Exec(context.Background(),
		`insert into "userdata" ("userid", "name", "surname", "description") values ($1, $2, $3, $4)`,
		userId, data.Name, data.Surname, data.Description)

	if err != nil {
		fmt.Println("db.Exec()", err)
		return -1
	}
	return userId
}

func DeleteUser(id int) error {
	db, err := openConnection()
	if err != nil {
		return err
	}

	defer db.Close(context.Background())

	var username string
	err = db.QueryRow(
		context.Background(), `SELECT "username" FROM "users" where id = $1`, id).Scan(&username)
	if err != nil {
		return err
	}

	if exists(username) != id {
		return fmt.Errorf("User with ID %d does not exist", id)
	}

	if exists(username) != id {
		return fmt.Errorf("User with ID %d does not exist", id)
	}

	deleteStatement := `delete from "userdata" where userid=$1`
	_, err = db.Exec(context.Background(), deleteStatement, id)
	if err != nil {
		return err
	}

	deleteStatement = `delete from "users" where id=$1`
	_, err = db.Exec(context.Background(), deleteStatement, id)
	if err != nil {
		return err
	}
	return nil
}

// BUG(1): Function ListUsers() not working as expected
func ListUsers() ([]Userdata, error) {
	Data := []Userdata{}
	db, err := openConnection()
	if err != nil {
		return Data, err
	}
	defer db.Close(context.Background())

	rows, err := db.Query(context.Background(), `SELECT "id", "username","name","surname","description" 
	FROM "users", "userdata" WHERE users.id = userdata.userid`)
	if err != nil {
		return Data, err
	}

	for rows.Next() {
		var id int
		var username string
		var surname string
		var name string
		var description string
		err = rows.Scan(&id, &username, &name, &surname, &description)
		temp := Userdata{
			ID:          id,
			Username:    username,
			Name:        name,
			Surname:     surname,
			Description: description,
		}

		Data = append(Data, temp)
		if err != nil {
			return Data, err
		}
	}
	defer rows.Close()
	return Data, nil
}

func UpdateUser(d Userdata) error {
	db, err := openConnection()
	if err != nil {
		return err
	}
	defer db.Close(context.Background())

	userId := exists(d.Username)
	if userId != -1 {
		return errors.New("User does not exist")
	}

	d.ID = userId
	updateStatement := `UPDATE "userdata" SET "name"=$1, "surname"=$2, "description"=$3 WHERE "userid"=$4`
	_, err = db.Exec(context.Background(), updateStatement, d.Name, d.Surname, d.Description, d.ID)
	if err != nil {
		return err
	}

	return nil
}
