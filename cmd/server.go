package main

import (
	"database/sql"
	"log"
	"net/http"

	. "github.com/leerwitz/projectToDo/internal/handlers"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

func main() {

	driverName := "postgres"
	databaseName := "user=postgres password=1980 dbname=postgres host=10.0.2.15 port=5432 sslmode=disable"
	database, err := sql.Open(driverName, databaseName)

	if err != nil {
		log.Fatal(err)
	}

	err = database.Ping()
	if err != nil {
		log.Fatal(err)
	}

	defer database.Close()

	router := mux.NewRouter()
	router.Use(enableCors)

	router.HandleFunc("/task", GetAllTask(database)).Methods("GET")
	router.HandleFunc("/task/{id}", GetTaskByID(database)).Methods("GET")
	router.HandleFunc("/task", PostTask(database)).Methods("POST", "OPTIONS")
	router.HandleFunc("/task/{id}", PutTaskById(database)).Methods("PUT")
	router.HandleFunc("/task/{id}", PatchTaskById(database)).Methods("PATCH")
	router.HandleFunc("/task/{id}", DeleteTaskById(database)).Methods("DELETE", "OPTIONS")

	router.NotFoundHandler = http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		if request.Method == `OPTIONS` {
			writer.WriteHeader(http.StatusNoContent)
		} else {
			http.Error(writer, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	http.ListenAndServe(":8080", router)
}
