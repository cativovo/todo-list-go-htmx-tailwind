package rest

import (
	"net/http"
	"text/template"

	"github.com/cativovo/todo-list-go-htmx-tailwind/pkg/todo"
	"github.com/go-chi/chi/v5"
)

type Pages struct {
	todoService todo.Service
}

func NewPages(r *chi.Mux, t todo.Service) {
	p := Pages{
		todoService: t,
	}

	r.Get("/", p.homePage)
}

func (p *Pages) homePage(w http.ResponseWriter, _ *http.Request) {
	todos := p.todoService.GetTodos()

	tmpl := template.Must(template.ParseFiles("web/template/base.html", "web/template/todo.html"))

	err := tmpl.Execute(w, todos)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	}
}
