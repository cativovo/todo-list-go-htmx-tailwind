package main

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"net/http"
	"text/template"

	_ "github.com/lib/pq"
)

type (
	Handler       func(w http.ResponseWriter, r *http.Request)
	HandlerWithDB func(db *sql.DB, w http.ResponseWriter, r *http.Request) error
)

// pages
func homePage(w http.ResponseWriter, _ *http.Request) {
	tmpl := template.Must(template.ParseFiles("web/templates/base.html"))
	tmpl.Execute(w, nil)
}

func todos(db *sql.DB, w http.ResponseWriter, _ *http.Request) error {
	rows, err := db.Query("SELECT * from todo")
	if err != nil {
		return err
	}

	for rows.Next() {
		id := ""
		task_name := ""
		updated_at := ""
		completed := false

		if err := rows.Scan(&id, &task_name, &updated_at, &completed); err != nil {
			return err
		}

		fmt.Println(id, task_name, updated_at, completed)
	}

	fmt.Fprint(w, "test")

	return nil
}

// partials
func handleAddTodo(db *sql.DB, w http.ResponseWriter, r *http.Request) error {
	todoText := r.PostFormValue("todo")

	tx, err := db.Begin()
	if err != nil {
		return err
	}

	stmt, err := tx.Prepare(`
      INSERT INTO
      todo (task_name)
      VALUES ($1)
      RETURNING id, updated_at
      `)
	if err != nil {
		return err
	}

	defer stmt.Close()

	id := ""
	updated_at := ""

	{
		err := stmt.QueryRow(todoText).Scan(&id, &updated_at)
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

	fmt.Println("Todo created:", id, todoText, updated_at)
	fmt.Fprint(w, todoText)

	return nil
}

// utils
func dbInit(db *sql.DB) {
	_, err := db.Exec(`
    CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

    CREATE TABLE IF NOT EXISTS todo (
      id uuid DEFAULT uuid_generate_v4() PRIMARY KEY,
      task_name VARCHAR(255) NOT NULL,
      updated_at TIMESTAMP DEFAULT NOW(),
      completed BOOLEAN DEFAULT false
    );
    `)
	if err != nil {
		log.Fatal("DB init failed")
	}

	fmt.Println("DB init success")
}

func createHandler(method string, db *sql.DB, handler HandlerWithDB) Handler {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != method {
			http.Error(w, "Not Found", http.StatusNotFound)
			return
		}

		err := handler(db, w, r)
		if err != nil {
			fmt.Println(err)
			fmt.Fprint(w, "Ooops something went wrong")
		}
	}
}

func main() {
	// DB connection
	connStr := "postgres://postgres:1234@127.0.0.1:5432/gotodo?sslmode=disable"

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal("DB connection error", err)
	}

	pingErr := db.Ping()

	if pingErr != nil {
		log.Fatal(pingErr)
	}

	fmt.Println("Connected to db!")

	dbInit(db)

	mux := http.NewServeMux()
	fmt.Println("running")

	// assets
	mux.Handle("/public/", http.StripPrefix("/public/", http.FileServer(http.Dir("./public"))))
	// pages
	mux.HandleFunc("/", homePage)
	mux.HandleFunc("/todos", createHandler(http.MethodGet, db, todos))
	// partials
	mux.HandleFunc("/add-todo", createHandler(http.MethodPost, db, handleAddTodo))

	httpErr := http.ListenAndServe("127.0.0.1:4000", mux)

	if httpErr != nil {
		if errors.Is(httpErr, http.ErrServerClosed) {
			log.Fatal("Server closed")
		} else {
			log.Fatal(httpErr)
		}
	}
}
