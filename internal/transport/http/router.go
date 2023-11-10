package http

import (
	"net/http"

	"github.com/Yury132/Golang-Task-2/internal/transport/http/handlers"
	"github.com/gorilla/mux"
)

func InitRoutes(h *handlers.Handler) *mux.Router {
	r := mux.NewRouter()

	r.HandleFunc("/uploads", h.Upload).Methods(http.MethodPost)
	r.HandleFunc("/health", h.Health).Methods(http.MethodGet)
	// Получаем информацию о картинках
	r.HandleFunc("/get-data", h.GetData).Methods(http.MethodGet)
	// Получаем информацию о картинках по id
	r.HandleFunc("/uploads/{id:[0-9]+}", h.GetDataId).Methods(http.MethodGet)

	return r
}
