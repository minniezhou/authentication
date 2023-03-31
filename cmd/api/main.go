package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
)

type Config struct {
	DB *sql.DB
}

const webPort = "80"

var count int64

func main() {
	fmt.Println("This is the authentication service...")
	// connect to DB
	db, err := ConnectToDB()
	if err != nil {
		log.Panic("Connect to DB failed!")
	}
	config := Config{
		DB: db,
	}
	h := config.NewHandler()
	err = http.ListenAndServe(fmt.Sprintf(":%s", webPort), h.router)
	if err != nil {
		log.Panic(err)
	}
}

func GetDB(dns string) (*sql.DB, error) {
	db, err := sql.Open("pgx", dns)
	if err != nil {
		return nil, err
	}
	err = db.Ping()
	if err != nil {
		return nil, err
	}
	return db, nil
}

func ConnectToDB() (*sql.DB, error) {
	dns := os.Getenv("PostgresDB")
	for {
		count++
		fmt.Printf("Connecting to DB attempting %d ...", count)
		db, err := GetDB(dns)
		if err != nil {
			fmt.Println("Data base is not ready yet...")
		} else {
			return db, nil
		}

		if count > 10 {
			return nil, err
		}

		fmt.Println("Try again after 2 seconds...")
		time.Sleep(2 * time.Second)
	}
}
