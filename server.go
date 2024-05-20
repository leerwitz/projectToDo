package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

func main() {

	driverName := "postgres"
	databaseName := "user=postgres password=1980 dbname=mydb host=10.0.2.15 port=5432 sslmode=disable"
	database, err := sql.Open(driverName, databaseName)
	if err != nil {
		log.Fatal(err)
	}

	defer database.Close()

	router := mux.NewRouter()
	router.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
		fmt.Fprintf(writer, "hello world from %s\n", request.URL.Path)
	})
	http.ListenAndServe(":8080", router)
}
