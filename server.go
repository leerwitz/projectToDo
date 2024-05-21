package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

type Task struct {
	Id     int64
	Title  string
	Text   string
	Author string
	Urgent bool
}

func GetAllTask(database *sql.DB) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		query := "SELECT id , title, text, author, urgent FROM task"
		rows, err := database.Query(query)
		if err != nil {
			http.Error(writer, err.Error(), http.StatusInternalServerError)
			return
		}

		defer rows.Close()
		var tasks []Task

		for rows.Next() {
			var curTask Task
			if err := rows.Scan(&curTask.Id, &curTask.Title, &curTask.Text, &curTask.Author, &curTask.Urgent); err != nil {
				http.Error(writer, err.Error(), http.StatusInternalServerError)
				return
			}
			tasks = append(tasks, curTask)
		}
		if err := rows.Err(); err != nil {
			http.Error(writer, err.Error(), http.StatusInternalServerError)
			return
		}
		writer.Header().Set("Content-Type", "application/json")
		writer.WriteHeader(http.StatusOK)

		if err := json.NewEncoder(writer).Encode(tasks); err != nil {
			http.Error(writer, err.Error(), http.StatusInternalServerError)
			return
		}

	}
}

func GetTaskByID(database *sql.DB) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		query := "SELECT id , title, text, author, urgent FROM task WHERE id=$1"
		variables := mux.Vars(request)
		var (
			err  error
			task Task
		)

		task.Id, err = strconv.ParseInt(variables["id"], 10, 64)

		if err != nil {
			http.Error(writer, err.Error(), http.StatusBadRequest)
			return
		}

		row := database.QueryRow(query, task.Id)

		if err := row.Err(); err != nil {
			if err == sql.ErrNoRows {
				http.Error(writer, err.Error(), http.StatusNotFound)
			} else {
				http.Error(writer, err.Error(), http.StatusInternalServerError)
			}
			return
		}

		if err := row.Scan(&task.Title, &task.Text, &task.Author, &task.Urgent); err != nil {
			http.Error(writer, err.Error(), http.StatusInternalServerError)
			return
		}

		writer.Header().Set("Content-Type", "application/json")
		writer.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(writer).Encode(task); err != nil {
			http.Error(writer, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

func PostTask(database *sql.DB) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		var task Task

		if err := json.NewDecoder(request.Body).Decode(&task); err != nil {
			http.Error(writer, err.Error(), http.StatusBadRequest)
			return
		}

		defer request.Body.Close()

		query := `INSERT INTO task tittle=$1, text=$2, author=$3, urgent=$4`
		result, err := database.Exec(query, task.Title, task.Text, task.Author, task.Urgent)

		if err != nil {
			http.Error(writer, err.Error(), http.StatusInternalServerError)
			return
		}
		Id, err := result.LastInsertId()
		if err != nil {
			http.Error(writer, err.Error(), http.StatusInternalServerError)
			return
		}

		task.Id = Id

		numRows, err := result.RowsAffected()
		if err != nil {
			http.Error(writer, err.Error(), http.StatusInternalServerError)
			return
		}

		writer.Header().Set("Content-Type", "application/json")
		if numRows == 0 {
			writer.WriteHeader(http.StatusOK)
		} else {
			writer.WriteHeader(http.StatusCreated)
		}

		if err := json.NewEncoder(writer).Encode(task); err != nil {
			http.Error(writer, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

func main() {

	driverName := "postgres"
	databaseName := "user=postgres password=1980 dbname=mydb host=10.0.2.15 port=5432 sslmode=disable"
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
	router.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
		fmt.Fprintf(writer, "hello world from %s\n", request.URL.Path)
	})

	router.HandleFunc("/task", GetAllTask(database)).Methods("GET")
	router.HandleFunc("/task/{id}", GetTaskByID(database)).Methods("GET")

	http.ListenAndServe(":8080", router)
}
