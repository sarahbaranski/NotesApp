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

// func main() {
// 	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
// 		"dbname=%s sslmode=disable",
// 		host, port, user, dbname)

// 	db, err := sql.Open("postgres", psqlInfo)
// 	if err != nil {
// 		panic(err)
// 	}
// 	defer db.Close()

// 	err = db.Ping()
// 	if err != nil {
// 		panic(err)
// 	}

// 	fmt.Println("Successfully connected!")
// }

// func dsn(dbName string) string {
// 	return fmt.Sprintf("postgres://%s:password@%s:%d/%s?sslmode=disable", user, host, port, dbName)
// }

func dbConnection() (*sql.DB, error) {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"dbname=%s sslmode=disable",
		host, port, user, dbname)

	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		fmt.Printf("Error %s when opening DB\n", err)
		return nil, err
	}

	// db, err := sql.Open("postgres", dns(dbName))
	// if err != nil {
	// 	fmt.Printf("Error %s when opening DB\n", err)
	// 	return nil, err
	// }

	ctx, cancelfunc := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelfunc()
	res, err := db.ExecContext(ctx, "CREATE DATABASE IF NOT EXISTS "+dbname)
	if err != nil {
		fmt.Printf("Error %s when creating DB\n", err)
		return nil, err
	}
	no, err := res.RowsAffected()
	if err != nil {
		fmt.Printf("Error %s when fetching rows", err)
		return nil, err
	}
	fmt.Printf("rows affected: %d\n", no)

	db.Close()
	db, err = sql.Open("postgres", psqlInfo)
	if err != nil {
		fmt.Printf("Error %s when opening DB", err)
		return nil, err
	}

	// db, err = sql.Open("postgres", dns(dbName))
	// if err != nil {
	// 	fmt.Printf("Error %s when opening DB", err)
	// 	return nil, err
	// }

	db.Close()
	db, err = sql.Open("postgres", fmt.Sprintf("host=%s port=%d user=%s "+
		"dbname=%s sslmode=disable", host, port, user, dbname))
	if err != nil {
		fmt.Printf("Error %s when opening DB", err)
		return nil, err
	}

	db.SetMaxOpenConns(20)
	db.SetMaxIdleConns(20)
	db.SetConnMaxLifetime(time.Minute * 5)

	ctx, cancelfunc = context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelfunc()
	err = db.PingContext(ctx)
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
