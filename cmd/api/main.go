package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	_ "github.com/lib/pq"
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
		fmt.Println(err.Error())
		log.Panic("Connect to DB failed!")
	}
	fmt.Println("Connect to Postgress successfully!")
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
	db, err := sql.Open("postgres", dns)
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
	dns := os.Getenv("POSTGRES_DB")
	fmt.Println("POSTGRES_DB: ", dns)
	for {
		count++
		fmt.Printf("Connecting to DB attempting %d ...\n", count)
		db, err := GetDB(dns)
		if err != nil {
			fmt.Println("Data base is not ready yet...")
		} else {
			return db, nil
		}

		if count > 2 {
			return nil, err
		}

		fmt.Println("Try again after 2 seconds...")
		time.Sleep(2 * time.Second)
	}
}
