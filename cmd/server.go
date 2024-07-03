package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	. "github.com/leerwitz/projectToDo/internal/task"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

func enableCors(next http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		writer.Header().Set("Access-Control-Allow-Origin", "*")
		writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS")
		writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if request.Method == http.MethodOptions {
			writer.WriteHeader(http.StatusNoContent)
			return
		}

		next.ServeHTTP(writer, request)
	})
}

func GetAllTask(database *sql.DB) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		tasks, err := GetAll(database)
		if err != nil {
			http.Error(writer, err.Error(), http.StatusInternalServerError)
			return
		}
		writer.Header().Set("Content-Type", "application/json")
		// writer.Header().Set("Access-Control-Allow-Origin", "http://127.0.0.1:5500")
		writer.WriteHeader(http.StatusOK)

		if err := json.NewEncoder(writer).Encode(tasks); err != nil {
			http.Error(writer, err.Error(), http.StatusInternalServerError)
			return
		}

	}
}

func GetTaskByID(database *sql.DB) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
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

		task, err = GetById(database, task.Id)
		if err != nil {
			if err == sql.ErrNoRows {
				http.Error(writer, err.Error(), http.StatusNotFound)
			} else {
				http.Error(writer, err.Error(), http.StatusInternalServerError)
			}
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

		err := Post(database, &task)
		if err != nil {
			http.Error(writer, err.Error(), http.StatusInternalServerError)
			return
		}

		writer.Header().Set("Content-Type", "application/json")
		// writer.Header().Set("Access-Control-Allow-Origin", "http://127.0.0.1:5500")
		writer.WriteHeader(http.StatusCreated)

		if err := json.NewEncoder(writer).Encode(task); err != nil {
			http.Error(writer, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

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

		numRows, err := Put(database, &task)

		if err != nil {
			http.Error(writer, err.Error(), http.StatusInternalServerError)
			return
		}

		writer.Header().Set("Content-Type", "application/json")
		if numRows == 0 {

			err := Post(database, &task)
			if err != nil {
				http.Error(writer, err.Error(), http.StatusInternalServerError)
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
		// query := `DELETE FROM task WHERE id=$1`
		// result, err := database.Exec(query, id)
		// if err != nil {
		// 	http.Error(writer, err.Error(), http.StatusInternalServerError)
		// 	return
		// }
		// rowsNum, err := result.RowsAffected()
		rowsNum, err := Delete(database, id)
		if err != nil {
			http.Error(writer, err.Error(), http.StatusInternalServerError)
			return
		}

		writer.Header().Set("Content-Type", "application/json")
		if rowsNum == 0 {
			http.Error(writer, "Task not found", http.StatusNotFound)
		} else {
			writer.WriteHeader(http.StatusNoContent)
		}

	}
}

func PatchTaskById(database *sql.DB) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		var (
			task Task
			err  error
		)

		variables := mux.Vars(request)

		if err := json.NewDecoder(request.Body).Decode(&task); err != nil {
			http.Error(writer, err.Error(), http.StatusBadRequest)
			return
		}

		defer request.Body.Close()

		if task.Id, err = strconv.ParseInt(variables["id"], 10, 64); err != nil {
			http.Error(writer, err.Error(), http.StatusBadRequest)
			return
		}

		paramsCount, numRows, err := Patch(database, &task)
		if err != nil {
			http.Error(writer, err.Error(), http.StatusInternalServerError)
			return
		}
		if paramsCount == 0 {
			http.Error(writer, "invalid query syntax", http.StatusBadRequest)
			return
		}
		if numRows == 0 {
			http.Error(writer, "invalid query syntax", http.StatusNotFound)
			return
		}

		writer.Header().Set("Content-Type", "application/json")
		writer.WriteHeader(http.StatusNoContent)
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
