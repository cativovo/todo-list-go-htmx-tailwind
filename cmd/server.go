package main

import (
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/cativovo/todo-list-go-htmx-tailwind/pkg/http/rest"
	"github.com/cativovo/todo-list-go-htmx-tailwind/pkg/storage/postgres"
	"github.com/cativovo/todo-list-go-htmx-tailwind/pkg/todo"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func main() {
	r := chi.NewRouter()
	r.Use(middleware.Logger)

	repository := postgres.NewPostgresRepository()
	todoService := todo.NewService(repository)

	rest.NewPages(r, todoService)
	rest.NewHandlers(r, todoService)

	// Create a route along /public that will serve contents from
	// the ./public/ folder.
	workDir, _ := os.Getwd()
	filesDir := http.Dir(filepath.Join(workDir, "public"))
	FileServer(r, "/public", filesDir)
	http.ListenAndServe("127.0.0.1:3000", r)
}

// FileServer conveniently sets up a http.FileServer handler to serve
// static files from a http.FileSystem.
func FileServer(r chi.Router, path string, root http.FileSystem) {
	if strings.ContainsAny(path, "{}*") {
		panic("FileServer does not permit any URL parameters.")
	}

	if path != "/" && path[len(path)-1] != '/' {
		r.Get(path, http.RedirectHandler(path+"/", http.StatusMovedPermanently).ServeHTTP)
		path += "/"
	}
	path += "*"

	r.Get(path, func(w http.ResponseWriter, r *http.Request) {
		rctx := chi.RouteContext(r.Context())
		pathPrefix := strings.TrimSuffix(rctx.RoutePattern(), "/*")
		fs := http.StripPrefix(pathPrefix, http.FileServer(root))
		fs.ServeHTTP(w, r)
	})
}
