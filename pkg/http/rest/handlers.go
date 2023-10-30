package rest

import (
	"log"
	"net/http"
	"text/template"

	"github.com/cativovo/todo-list-go-htmx-tailwind/pkg/todo"
	"github.com/go-chi/chi/v5"
)

type Handlers struct {
	todoService todo.Service
}

func NewHandlers(r *chi.Mux, t todo.Service) {
	h := Handlers{
		todoService: t,
	}

	r.Route("/todo", func(r chi.Router) {
		r.Post("/add", h.addTodo)
		r.Patch("/update-completed/{todoId}", h.updateCompleted)
		r.Patch("/update-taskname/{todoId}", h.updateTaskName)
		r.Delete("/delete/{todoId}", h.deleteTodo)
	})
}

func (h *Handlers) addTodo(w http.ResponseWriter, r *http.Request) {
	taskName := r.PostFormValue("taskName")
	todo, err := h.todoService.AddTodo(taskName)
	if err != nil {
		log.Println(err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)

		return
	}

	tmpl := template.Must(template.ParseFiles("web/template/todo.html"))

	if err := tmpl.Execute(w, todo); err != nil {
		log.Println(err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)

		return
	}
}

func (h *Handlers) updateCompleted(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "todoId")
	completed := r.PostFormValue("completed") == "true"

	todo, err := h.todoService.UpdateCompleted(id, completed)
	if err != nil {
		log.Println(err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)

		return
	}

	tmpl := template.Must(template.ParseFiles("web/template/todo.html"))

	if err := tmpl.Execute(w, todo); err != nil {
		log.Println(err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)

		return
	}
}

func (h *Handlers) updateTaskName(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "todoId")
	taskName := r.PostFormValue("taskName")

	todo, err := h.todoService.UpdateTaskName(id, taskName)
	if err != nil {
		log.Println(err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)

		return
	}

	tmpl := template.Must(template.ParseFiles("web/template/todo.html"))

	if err := tmpl.Execute(w, todo); err != nil {
		log.Println(err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)

		return
	}
}

func (h *Handlers) deleteTodo(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "todoId")

	if _, err := h.todoService.DeleteTodo(id); err != nil {
		log.Println(err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)

		return
	}
}
