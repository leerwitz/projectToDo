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

type HTTPError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func EnableCors(next http.Handler) http.Handler {
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

// GetAllTaskByTitle godoc
// @Summary Get all tasks by title
// @Description Дает содержимое обо всех задачах, название которых начитнается с аргумента
// title из query запроса, если он не указан, то выводит все задачи.
// @Tags handlers
// @Produce  json
// @Param title query string false "Фильтр по названию"
// @Success 200 {array} task.Task
// @Failure 500 {object} HTTPError
// @Router /task [get]
func GetAllTaskByTitle(database *sql.DB) http.HandlerFunc {
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

// GetTaskByID godoc
// @Summary Get task by ID
// @Description Получает задачи по ее айди.
// @Tags handlers
// @Produce  json
// @Param id path int true "Task ID"
// @Success 200 {object} task.Task
// @Failure 400 {object} HTTPError
// @Failure 404 {object} HTTPError
// @Failure 500 {object} HTTPError
// @Router /task/{id} [get]
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

// PostTask godoc
// @Summary Post task
// @Description Создает задачу по json из тела запроса
// @Tags handlers
// @Accept  json
// @Produce  json
// @Success 201 {object} task.Task
// @Failure 400 {object} HTTPError
// @Failure 404 {object} HTTPError
// @Failure 500 {object} HTTPError
// @Router /task [post]
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
		writer.WriteHeader(http.StatusCreated)

		if err := json.NewEncoder(writer).Encode(task); err != nil {
			http.Error(writer, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

//	DeleteTaskById godoc
//
// @Summary Delete task by ID
// @Description Удаляет задачу по ее айди.
// @Tags handlers
// @Param id path int true "Task ID"
// @Success 204 "No content"
// @Failure 400 {object} HTTPError
// @Failure 404 {object} HTTPError
// @Failure 500 {object} HTTPError
// @Router /task/{id} [delete]
func DeleteTaskById(database *sql.DB) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {

		variables := mux.Vars(request)
		id, err := strconv.ParseInt(variables["id"], 10, 64)

		if err != nil {
			http.Error(writer, err.Error(), http.StatusBadRequest)
			return
		}

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

// PatchTaskById godoc
// @Summary Patch task by ID
// @Description Обновляет задачу по ее айди, не удаляя старый и создавая новый объекты,
// если задачи с таким айди нет выбрасывает 404.
// @Tags handlers
// @Param id path int true "Task ID"
// @Success 204 "No content"
// @Failure 400 {object} HTTPError
// @Failure 404 {object} HTTPError
// @Failure 500 {object} HTTPError
// @Router /task/{id} [patch]
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
