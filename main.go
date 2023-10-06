package main

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"text/template"
)

// pages
func homePage(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("web/templates/base.html"))
	tmpl.Execute(w, nil)
}

// partials
func handleAddTodo(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Not Found", http.StatusNotFound)
		return
	}

	todoText := r.PostFormValue("todo")

	fmt.Println(todoText)

	fmt.Fprint(w, todoText)
}

func main() {
	mux := http.NewServeMux()
	fmt.Println("running")

	// assets
	mux.Handle("/public/", http.StripPrefix("/public/", http.FileServer(http.Dir("./public"))))
	// pages
	mux.HandleFunc("/", homePage)
	// partials
	mux.HandleFunc("/add-todo", handleAddTodo)

	err := http.ListenAndServe("127.0.0.1:4000", mux)

	if err != nil {
		if errors.Is(err, http.ErrServerClosed) {
			log.Fatal("Server closed")
		} else {
			log.Fatal(err)
		}
	}
}
