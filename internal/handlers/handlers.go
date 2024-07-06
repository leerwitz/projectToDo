package handlers

import (
	"database/sql"
	"encoding/json"
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
		var tasks []Task
		var err error
		patternTitle := request.URL.Query().Get(`title`)

		if patternTitle != `` {
			tasks, err = GetByTitle(database, patternTitle)
		} else {
			tasks, err = GetAll(database)
		}

		if err != nil {
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