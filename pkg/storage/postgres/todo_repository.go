package postgres

import (
	"database/sql"
	"log"
	"time"

	"github.com/cativovo/todo-list-go-htmx-tailwind/pkg/model"
	"github.com/lib/pq"
)

type TodoRepository struct {
	db *sql.DB
}

func NewTodoRepository(db *sql.DB) *TodoRepository {
	_, err := db.Exec(`
    CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

    CREATE TABLE IF NOT EXISTS todo (
      id uuid DEFAULT uuid_generate_v4() PRIMARY KEY,
      task_name VARCHAR(255) NOT NULL,
      updated_at TIMESTAMP,
      completed BOOLEAN DEFAULT false
    );
    `)
	if err != nil {
		log.Fatal("DB init failed")
	}

	log.Println("DB init success")

	return &TodoRepository{
		db: db,
	}
}

func (tr *TodoRepository) GetAllTodos() ([]model.Todo, error) {
	var todos []model.Todo

	rows, err := tr.db.Query("SELECT id, task_name, updated_at, completed from todo ORDER BY completed ASC, updated_at DESC")
	if err != nil {
		return todos, err
	}

	for rows.Next() {
		var todo model.Todo

		err = rows.Scan(&todo.Id, &todo.TaskName, &todo.UpdatedAt, &todo.Completed)
		if err != nil {
			return todos, err
		}

		todos = append(todos, todo)
	}

	return todos, nil
}

func (tr *TodoRepository) AddTodo(taskName string) (model.Todo, error) {
	todo := model.Todo{
		TaskName: taskName,
	}

	tx, err := tr.db.Begin()
	if err != nil {
		return todo, err
	}

	stmt, err := tx.Prepare(`
      INSERT INTO
      todo (task_name, updated_at)
      VALUES ($1, $2)
      RETURNING id, updated_at
      `)
	if err != nil {
		return todo, err
	}

	defer stmt.Close()

	timestamp := pq.FormatTimestamp(time.Now())

	{
		err := stmt.QueryRow(taskName, timestamp).Scan(&todo.Id, &todo.UpdatedAt)
		if err != nil {
			return todo, err
		}
	}

	{
		err := tx.Commit()
		if err != nil {
			return todo, err
		}
	}

	return todo, nil
}

func (tr *TodoRepository) UpdateTaskName(id, taskName string) (model.Todo, error) {
	todo := model.Todo{
		Id: id,
	}

	tx, err := tr.db.Begin()
	if err != nil {
		return todo, err
	}

	stmt, err := tx.Prepare(`
	     UPDATE todo
       SET task_name = $1, updated_at = $2
       WHERE id = $3
       RETURNING task_name, completed, updated_at
	     `)
	if err != nil {
		return todo, err
	}

	defer stmt.Close()

	timestamp := pq.FormatTimestamp(time.Now())

	err = stmt.QueryRow(taskName, timestamp, id).Scan(&todo.TaskName, &todo.Completed, &todo.UpdatedAt)
	if err != nil {
		return todo, err
	}

	{
		err := tx.Commit()
		if err != nil {
			return todo, err
		}
	}

	return todo, nil
}

func (tr *TodoRepository) UpdateCompleted(id string, completed bool) (model.Todo, error) {
	todo := model.Todo{
		Id: id,
	}

	tx, err := tr.db.Begin()
	if err != nil {
		return todo, err
	}

	stmt, err := tx.Prepare(`
	     UPDATE todo
       SET completed = $1, updated_at = $2
       WHERE id = $3
       RETURNING task_name, completed, updated_at
	     `)
	if err != nil {
		return todo, err
	}

	defer stmt.Close()

	timestamp := pq.FormatTimestamp(time.Now())

	err = stmt.QueryRow(completed, timestamp, id).Scan(&todo.TaskName, &todo.Completed, &todo.UpdatedAt)
	if err != nil {
		return todo, err
	}

	{
		err := tx.Commit()
		if err != nil {
			return todo, err
		}
	}

	return todo, nil
}

func (tr *TodoRepository) DeleteTodo(id string) (bool, error) {
	tx, err := tr.db.Begin()
	if err != nil {
		return false, err
	}

	stmt, err := tx.Prepare("DELETE FROM todo WHERE id = $1")
	if err != nil {
		return false, err
	}

	defer stmt.Close()

	if _, err := stmt.Exec(id); err != nil {
		return false, err
	}

	if err := tx.Commit(); err != nil {
		return false, err
	}

	return true, nil
}
