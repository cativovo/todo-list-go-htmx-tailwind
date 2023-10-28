package main

import (
	"database/sql"
	"errors"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/lib/pq"
)

type (
	Handler       func(w http.ResponseWriter, r *http.Request)
	HandlerWithDB func(db *sql.DB, w http.ResponseWriter, r *http.Request) error
)

type Todo struct {
	Id        string
	TaskName  string
	UpdatedAt string
	Completed bool
}

// pages
func homePage(db *sql.DB, w http.ResponseWriter, _ *http.Request) error {
	rows, err := db.Query("SELECT id, task_name, updated_at, completed from todo ORDER BY completed ASC, updated_at DESC")
	if err != nil {
		return err
	}

	todos := []Todo{}

	for rows.Next() {
		todo := Todo{}
		err = rows.Scan(&todo.Id, &todo.TaskName, &todo.UpdatedAt, &todo.Completed)
		if err != nil {
			return err
		}

		todos = append(todos, todo)
	}

	tmpl := template.Must(template.ParseFiles("web/templates/base.html", "web/templates/todo.html"))

	err = tmpl.Execute(w, todos)

	if err != nil {
		return err
	}

	return nil
}

// partials
func addTodo(db *sql.DB, w http.ResponseWriter, r *http.Request) error {
	taskName := r.PostFormValue("taskName")

	tx, err := db.Begin()
	if err != nil {
		return err
	}

	stmt, err := tx.Prepare(`
      INSERT INTO
      todo (task_name, updated_at)
      VALUES ($1, $2)
      RETURNING id, updated_at
      `)
	if err != nil {
		return err
	}

	defer stmt.Close()

	id := ""
	updated_at := ""
	timestamp := pq.FormatTimestamp(time.Now())

	{
		err := stmt.QueryRow(taskName, timestamp).Scan(&id, &updated_at)
		if err != nil {
			return err
		}
	}

	{
		err := tx.Commit()
		if err != nil {
			return err
		}
	}

	w.WriteHeader(http.StatusCreated)
	todoTmpl := template.Must(template.ParseFiles("web/templates/todo.html"))
	err = todoTmpl.Execute(w, map[string]interface{}{
		"Id":        id,
		"TaskName":  taskName,
		"UpdatedAt": updated_at,
	})

	return err
}

func deleteTodo(db *sql.DB, w http.ResponseWriter, r *http.Request) error {
	segments := strings.Split(r.URL.Path, "/")
	id := segments[len(segments)-1]

	tx, err := db.Begin()
	if err != nil {
		return err
	}

	stmt, err := tx.Prepare("DELETE FROM todo WHERE id = $1")
	if err != nil {
		return err
	}

	defer stmt.Close()

	result, err := stmt.Exec(id)
	if err != nil {
		return err
	}

	log.Println(result)

	{
		err := tx.Commit()
		if err != nil {
			return err
		}
	}

	w.WriteHeader(http.StatusOK)

	return err
}

func updateTodo(db *sql.DB, w http.ResponseWriter, r *http.Request) error {
	segments := strings.Split(r.URL.Path, "/")
	id := segments[len(segments)-1]
	taskName := r.PostFormValue("taskName")
	completed := r.PostFormValue("completed") == "true"

	tx, err := db.Begin()
	if err != nil {
		return err
	}

	taskNameResult := ""
	completedResult := false
	updatedAt := ""
	timestamp := pq.FormatTimestamp(time.Now())

	// TODO: find a better way to do this
	if taskName != "" {
		stmt, err := tx.Prepare(`
	     UPDATE todo
       SET task_name = $1, updated_at = $2
       WHERE id = $3
       RETURNING task_name, completed, updated_at
	     `)
		if err != nil {
			return err
		}

		defer stmt.Close()

		err = stmt.QueryRow(taskName, timestamp, id).Scan(&taskNameResult, &completedResult, &updatedAt)
		if err != nil {
			return err
		}
	} else {
		stmt, err := tx.Prepare(`
	     UPDATE todo
       SET completed = $1, updated_at = $2
       WHERE id = $3
       RETURNING task_name, completed, updated_at
	     `)
		if err != nil {
			return err
		}

		defer stmt.Close()

		err = stmt.QueryRow(completed, timestamp, id).Scan(&taskNameResult, &completedResult, &updatedAt)
		if err != nil {
			return err
		}
	}

	{
		err := tx.Commit()
		if err != nil {
			return err
		}
	}

	tmpl := template.Must(template.ParseFiles("web/templates/todo.html"))
	log.Println(taskNameResult, completedResult)

	w.WriteHeader(http.StatusOK)
	tmpl.Execute(w, map[string]interface{}{
		"Id":        id,
		"TaskName":  taskNameResult,
		"Completed": completedResult,
		"UpdatedAt": updatedAt,
	})

	return nil
}

// utils
func dbInit(db *sql.DB) {
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
}

func createHandler(method string, db *sql.DB, handler HandlerWithDB) Handler {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != method {
			http.Error(w, "Not Found", http.StatusNotFound)
			return
		}

		w.Header().Set("Content-Type", "text/html")

		err := handler(db, w, r)
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}
}

func main() {
	// DB connection
	dbUsername := os.Getenv("DB_USERNAME")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbHost := os.Getenv("DB_HOST")
	dbName := os.Getenv("DB_NAME")
	dbSslMode := os.Getenv("DB_SSLMODE")
	// postgres://<dbUsername>:<dbPassword>@<dbHost>/<dbName>?sslmode=<dbSslMode>
	connStr := fmt.Sprintf("postgres://%s:%s@%s/%s?sslmode=%s", dbUsername, dbPassword, dbHost, dbName, dbSslMode)

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal("DB connection error", err)
	}

	pingErr := db.Ping()

	if pingErr != nil {
		log.Fatal(pingErr)
	}

	log.Println("Connected to db!")

	dbInit(db)

	mux := http.NewServeMux()
	log.Println("running")

	// assets
	mux.Handle("/public/", http.StripPrefix("/public/", http.FileServer(http.Dir("./public"))))
	// pages
	mux.HandleFunc("/", createHandler(http.MethodGet, db, homePage))
	// partials
	mux.HandleFunc("/todo/add", createHandler(http.MethodPost, db, addTodo))
	mux.HandleFunc("/todo/delete/", createHandler(http.MethodDelete, db, deleteTodo))
	mux.HandleFunc("/todo/update/", createHandler(http.MethodPatch, db, updateTodo))

	httpErr := http.ListenAndServe("127.0.0.1:4000", mux)

	if httpErr != nil {
		if errors.Is(httpErr, http.ErrServerClosed) {
			log.Fatal("Server closed")
		} else {
			log.Fatal(httpErr)
		}
	}
}
