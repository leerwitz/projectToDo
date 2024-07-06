package main

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	_ "github.com/leerwitz/projectToDo/docs"
	. "github.com/leerwitz/projectToDo/internal/handlers"
	"github.com/rs/cors"
	httpSwagger "github.com/swaggo/http-swagger"
)

// @title Task API
// @version 1.0
// @description This is a sample server for managing tasks.
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:8080
// @BasePath /
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

	router.HandleFunc("/task", GetAllTaskByTitle(database)).Methods("GET")
	router.HandleFunc("/task/{id}", GetTaskByID(database)).Methods("GET")
	router.HandleFunc("/task", PostTask(database)).Methods("POST", "OPTIONS")
	router.HandleFunc("/task/{id}", PutTaskById(database)).Methods("PUT")
	router.HandleFunc("/task/{id}", PatchTaskById(database)).Methods("PATCH")
	router.HandleFunc("/task/{id}", DeleteTaskById(database)).Methods("DELETE", "OPTIONS")

	router.PathPrefix("/swagger/").Handler(httpSwagger.WrapHandler)

	handler := cors.Default().Handler(router)
	log.Fatal(http.ListenAndServe(":8080", handler))

}
