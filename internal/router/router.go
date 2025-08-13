package router

import (
	"lotest/internal/handlers"
	"lotest/internal/middleware"
	"net/http"
)

func NewRouter(h *handlers.Handler) *http.ServeMux {
	mux := http.NewServeMux()

	withLogging := func(next http.HandlerFunc) http.HandlerFunc {
		return middleware.LogHandler(next, h.Log)
	}

	mux.HandleFunc("/tasks", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			withLogging(h.CreateTask).ServeHTTP(w, r)
		case http.MethodGet:
			withLogging(h.GetAllTasks).ServeHTTP(w, r)
		default:
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		}
	})

	mux.HandleFunc("/task/", func(w http.ResponseWriter, r *http.Request) {
		withLogging(h.GetTaskByID).ServeHTTP(w, r)
	})

	return mux
}
