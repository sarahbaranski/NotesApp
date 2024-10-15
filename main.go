package main

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	_ "github.com/lib/pq"
)

const (
	host   = "localhost"
	port   = 5432
	user   = "sbaranski"
	dbname = "notes_app"
)

func dbConnection() (*sql.DB, error) {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"dbname=%s sslmode=disable",
		host, port, user, dbname)

	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		fmt.Printf("Error %s when opening DB\n", err)
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		fmt.Printf("Errors %s pinging DB", err)
		return nil, err
	}
	fmt.Printf("Connected to DB %s successfully\n", dbname)
	return db, nil
}

func createToDoTable(db *sql.DB) error {
	query := `CREATE TABLE IF NOT EXISTS todo_items(id varchar(36), name varchar(255), completed boolean)`
	ctx, cancelfunc := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelfunc()
	res, err := db.ExecContext(ctx, query)
	if err != nil {
		fmt.Printf("Error %s when creating todo_items table", err)
		return err
	}
	rows, err := res.RowsAffected()
	if err != nil {
		fmt.Printf("Error %s when getting rows affected", err)
		return err
	}
	fmt.Printf("Rows affected when creating table: %d", rows)
	return nil
}

func main() {
	db, err := dbConnection()
	if err != nil {
		fmt.Printf("Error %s when getting db connection", err)
		return
	}
	defer db.Close()
	fmt.Printf("Successfully connected to database")
	err = createToDoTable(db)
	if err != nil {
		fmt.Printf("Create todo_items table failed with error %s", err)
		return
	}
}
