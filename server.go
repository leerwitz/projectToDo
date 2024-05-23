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

// func rowTaskById(database *sql.DB, writer *http.ResponseWriter, task *Task) (*sql.Result, error) {
// 	query := "SELECT id , title, text, author, urgent FROM task WHERE id=$1"
// 	row := database.QueryRow(query, task.Id)
// 	if err := row.Err(); err != nil {
// 		if err == sql.ErrNoRows {
// 			http.Error(*writer, err.Error(), http.StatusNotFound)
// 		} else {
// 			http.Error(*writer, err.Error(), http.StatusInternalServerError)
// 		}
// 		return
// 	}
// }

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

		if err := row.Scan(&task.Id, &task.Title, &task.Text, &task.Author, &task.Urgent); err != nil {
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

		query := `INSERT INTO task (title, text, author, urgent) VALUES ($1, $2, $3, $4) RETURNING id`
		err := database.QueryRow(query, task.Title, task.Text, task.Author, task.Urgent).Scan(&task.Id)

		if err != nil {
			http.Error(writer, err.Error(), http.StatusInternalServerError)
			return
		}

		writer.Header().Set("Content-Type", "application/json")
		writer.WriteHeader(http.StatusCreated)

		if err := json.NewEncoder(writer).Encode(task); err != nil {
			http.Error(writer, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

// func postTaskIntoDB(task *Task, database *sql.DB, writer *http.ResponseWriter) error {
// 	query := `INSERT INTO task (title, text, author, urgent) VALUES ($1, $2, $3, $4) RETURNING id`
// 	err := database.QueryRow(query, (*task).Title, (*task).Text, (*task).Author, (*task).Urgent).Scan(task.Id)

// 	if err != nil {
// 		http.Error(*writer, err.Error(), http.StatusInternalServerError)
// 	}
// 	return err
// }

func PutTaskById(database *sql.DB) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		variables := mux.Vars(request)

		var (
			task Task
			err  error
		)

		if err := json.NewDecoder(request.Body).Decode(&task); err != nil {
			http.Error(writer, err.Error(), http.StatusBadRequest)
			return
		}

		defer request.Body.Close()

		task.Id, err = strconv.ParseInt(variables["id"], 10, 64)

		if err != nil {
			http.Error(writer, err.Error(), http.StatusBadRequest)
			return
		}

		query := `UPDATE task SET title = $1, text = $2, author = $3, urgent = $4 WHERE id = $5`
		result, err := database.Exec(query, task.Title, task.Text, task.Author, task.Urgent, task.Id)

		if err != nil {
			http.Error(writer, err.Error(), http.StatusInternalServerError)
			return
		}

		numRows, err := result.RowsAffected()

		if err != nil {
			http.Error(writer, err.Error(), http.StatusInternalServerError)
			return
		}

		writer.Header().Set("Content-Type", "application/json")
		if numRows == 0 {
			query := `INSERT INTO task (title, text, author, urgent) VALUES ($1, $2, $3, $4) RETURNING id`
			err := database.QueryRow(query, task.Title, task.Text, task.Author, task.Urgent).Scan(&task.Id)

			if err != nil {
				http.Error(writer, err.Error(), http.StatusInternalServerError)
			}
			if err != nil {
				return
			}
			writer.WriteHeader(http.StatusCreated)
		} else {
			writer.WriteHeader(http.StatusOK)
		}

		if err := json.NewEncoder(writer).Encode(task); err != nil {
			http.Error(writer, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

func DeleteTaskById(database *sql.DB) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {

		variables := mux.Vars(request)
		id, err := strconv.ParseInt(variables["id"], 10, 64)
		if err != nil {
			http.Error(writer, err.Error(), http.StatusBadRequest)
			return
		}
		query := `DELETE FROM task WHERE id=$1`
		result, err := database.Exec(query, id)
		if err != nil {
			http.Error(writer, err.Error(), http.StatusInternalServerError)
			return
		}
		rowsNum, err := result.RowsAffected()

		if err != nil {
			http.Error(writer, err.Error(), http.StatusInternalServerError)
			return
		}

		writer.Header().Set("Content-Type", "application/json")
		if rowsNum == 0 {
			writer.WriteHeader(http.StatusNoContent)
		} else {
			writer.WriteHeader(http.StatusOK)
		}

	}
}

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
	router.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
		fmt.Fprintf(writer, "hello world from %s\n", request.URL.Path)
	})

	router.HandleFunc("/task", GetAllTask(database)).Methods("GET")
	router.HandleFunc("/task/{id}", GetTaskByID(database)).Methods("GET")
	router.HandleFunc("/task", PostTask(database)).Methods("POST")
	router.HandleFunc("/task/{id}", PutTaskById(database)).Methods("PUT")
	router.HandleFunc("/task/{id}", DeleteTaskById(database)).Methods("DELETE")

	http.ListenAndServe(":8080", router)
}
