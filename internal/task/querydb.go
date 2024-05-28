package task

import (
	"database/sql"
	"fmt"
	"strings"

	_ "github.com/lib/pq"
)

type Task struct {
	Id     int64  `json:"id"`
	Title  string `json:"title,omitempty"`
	Text   string `json:"text,omitempty"`
	Author string `json:"author,omitempty"`
	Urgent *bool  `json:"urgent,omitempty"`
}

func GetAll(database *sql.DB) ([]Task, error) {
	query := "SELECT id , title, text, author, urgent FROM task"
	rows, err := database.Query(query)
	if err != nil {
		return nil, err
	}

	defer rows.Close()
	var tasks []Task

	for rows.Next() {
		var curTask Task
		if err := rows.Scan(&curTask.Id, &curTask.Title, &curTask.Text, &curTask.Author, &curTask.Urgent); err != nil {
			return nil, err
		}
		tasks = append(tasks, curTask)
	}

	return tasks, err
}

func GetById(database *sql.DB, id int64) (Task, error) {
	query := "SELECT id , title, text, author, urgent FROM task WHERE id=$1"
	var task Task
	task.Id = id
	row := database.QueryRow(query, task.Id)

	if err := row.Err(); err != nil {
		return task, err
	}

	err := row.Scan(&task.Id, &task.Title, &task.Text, &task.Author, &task.Urgent)

	return task, err
}

func Post(database *sql.DB, task *Task) error {
	query := `INSERT INTO task (title, text, author, urgent) VALUES ($1, $2, $3, $4) RETURNING id`
	err := database.QueryRow(query, task.Title, task.Text, task.Author, task.Urgent).Scan(&task.Id)

	return err
}

func Put(database *sql.DB, task *Task) (int64, error) {
	query := `UPDATE task SET title = $1, text = $2, author = $3, urgent = $4 WHERE id = $5`
	result, err := database.Exec(query, task.Title, task.Text, task.Author, task.Urgent, task.Id)

	if err != nil {
		return 0, err
	}
	numRows, err := result.RowsAffected()

	return numRows, err
}

func Delete(database *sql.DB, id int64) (int64, error) {
	query := `DELETE FROM task WHERE id=$1`
	result, err := database.Exec(query, id)
	if err != nil {
		return 0, err
	}
	rowsNum, err := result.RowsAffected()

	return rowsNum, err
}

func Patch(database *sql.DB, task *Task) (int64, int64, error) {
	var (
		updates     []string
		params      []interface{}
		paramsCount int64
	)

	if task.Title != "" {
		paramsCount++
		updates = append(updates, fmt.Sprintf("title = $%d", paramsCount))
		params = append(params, task.Title)
	}

	if task.Text != "" {
		paramsCount++
		updates = append(updates, fmt.Sprintf("text = $%d", paramsCount))
		params = append(params, task.Text)
	}

	if task.Author != "" {
		paramsCount++
		updates = append(updates, fmt.Sprintf("author = $%d", paramsCount))
		params = append(params, task.Author)
	}

	if task.Urgent != nil {
		paramsCount++
		updates = append(updates, fmt.Sprintf("author = $%d", paramsCount))
		params = append(params, task.Urgent)
	}

	if paramsCount == 0 {
		return 0, 0, nil
	}
	paramsCount++
	query := "UPDATE task SET " + strings.Join(updates, ", ") + fmt.Sprintf(" WHERE id = $%d", paramsCount)
	fmt.Println(query)
	params = append(params, task.Id)

	result, err := database.Exec(query, params...)
	if err != nil {
		return paramsCount, 0, err
	}
	numRows, err := result.RowsAffected()

	return paramsCount, numRows, err
}
