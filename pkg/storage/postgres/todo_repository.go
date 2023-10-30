package postgres

import (
	"log"
	"time"

	"github.com/cativovo/todo-list-go-htmx-tailwind/pkg/todo"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
)

type TodoRepository struct {
	db *sqlx.DB
}

func NewTodoRepository(db *sqlx.DB) *TodoRepository {
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

func (tr *TodoRepository) GetAllTodos() ([]todo.Todo, error) {
	var todos []todo.Todo

	rows, err := tr.db.Queryx("SELECT id, task_name, updated_at, completed from todo ORDER BY completed ASC, updated_at DESC")
	if err != nil {
		return todos, err
	}

	for rows.Next() {
		var t Todo

		err = rows.StructScan(&t)
		if err != nil {
			return todos, err
		}

		todos = append(todos, todo.Todo{
			Id:        t.Id,
			TaskName:  t.TaskName,
			UpdatedAt: t.UpdatedAt,
			Completed: t.Completed,
		})
	}

	return todos, nil
}

func (tr *TodoRepository) AddTodo(taskName string) (todo.Todo, error) {
	var t Todo
	timestamp := pq.FormatTimestamp(time.Now())

	if err := tr.db.Get(
		&t,
		`
    INSERT INTO todo (task_name, updated_at)
    VALUES ($1, $2)
    RETURNING id, updated_at
    `,
		taskName,
		timestamp,
	); err != nil {
		return todo.Todo{}, err
	}

	todo := todo.Todo{
		Id:        t.Id,
		TaskName:  taskName,
		UpdatedAt: t.UpdatedAt,
	}

	return todo, nil
}

func (tr *TodoRepository) UpdateTaskName(id, taskName string) (todo.Todo, error) {
	var t Todo
	timestamp := pq.FormatTimestamp(time.Now())

	if err := tr.db.Get(
		&t,
		`
    UPDATE todo
    SET task_name = $1, updated_at = $2
    WHERE id = $3
    RETURNING task_name, completed, updated_at
    `,
		taskName,
		timestamp,
		id,
	); err != nil {
		return todo.Todo{}, err
	}

	todo := todo.Todo{
		Id:        id,
		TaskName:  t.TaskName,
		UpdatedAt: t.UpdatedAt,
		Completed: t.Completed,
	}

	return todo, nil
}

func (tr *TodoRepository) UpdateCompleted(id string, completed bool) (todo.Todo, error) {
	var t Todo
	timestamp := pq.FormatTimestamp(time.Now())

	if err := tr.db.Get(
		&t,
		`
    UPDATE todo
    SET completed = $1, updated_at = $2
    WHERE id = $3
    RETURNING task_name, completed, updated_at
    `,
		completed,
		timestamp,
		id,
	); err != nil {
		return todo.Todo{}, err
	}

	todo := todo.Todo{
		Id:        id,
		TaskName:  t.TaskName,
		UpdatedAt: t.UpdatedAt,
		Completed: t.Completed,
	}

	return todo, nil
}

func (tr *TodoRepository) DeleteTodo(id string) (bool, error) {
	if _, err := tr.db.Exec("DELETE FROM todo WHERE id = $1", id); err != nil {
		return false, err
	}

	return true, nil
}
