package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	_ "github.com/leerwitz/projectToDo/docs"
	. "github.com/leerwitz/projectToDo/internal/handlers"
	"github.com/rs/cors"
	httpSwagger "github.com/swaggo/http-swagger"
)

const clientURl string = `http://localhost`
const listenPort string = `:8080`

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
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")

	driverName := "postgres"
	databaseName := fmt.Sprintf("user=%s password=%s dbname=%s host=%s port=%s sslmode=disable",
		dbUser, dbPassword, dbName, dbHost, dbPort)

	database, err := sql.Open(driverName, databaseName)

	if err != nil {
		log.Fatal(err)
	}

	err = database.Ping()
	if err != nil {
		log.Fatal(err)
	}

	defer database.Close()
	log.Println("Successfully connected to the database!")

	router := mux.NewRouter()

	router.HandleFunc("/task", GetAllTaskByTitle(database)).Methods("GET")
	router.HandleFunc("/task/{id}", GetTaskByID(database)).Methods("GET")
	router.HandleFunc("/task", PostTask(database)).Methods("POST", "OPTIONS")
	router.HandleFunc("/task/{id}", PatchTaskById(database)).Methods("PATCH")
	router.HandleFunc("/task/{id}", DeleteTaskById(database)).Methods("DELETE", "OPTIONS")

	router.PathPrefix("/swagger/").Handler(httpSwagger.WrapHandler)

	handler := cors.New(cors.Options{
		AllowedOrigins:   []string{clientURl},
		AllowedMethods:   []string{"GET", "POST", "PATCH", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Content-Type", "Authorization"},
		AllowCredentials: true,
	}).Handler(router)
	log.Fatal(http.ListenAndServe(listenPort, handler))

}
