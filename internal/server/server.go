package server

import (
	"log"
	"net/http"
	"time"

	"github.com/Der1an0/6sprint/internal/handlers"
)

// Создаем структуру сервера

type Server struct {
	logger *log.Logger
	server *http.Server
}

// Создаем новый сервер
func New(logger *log.Logger) *Server {
	// Создаем роутер
	router := createRouter()

	// Создаем сервер
	httpServer := &http.Server{
		Addr:         ":8080",
		Handler:      router,
		ErrorLog:     logger,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  15 * time.Second,
	}

	return &Server{
		logger: logger,
		server: httpServer,
	}
}

// Создаем и настраиваем роутер
func createRouter() *http.ServeMux {
	router := http.NewServeMux()

	// Регистрируем хендлеры
	router.HandleFunc("/", handlers.IndexHandler)
	router.HandleFunc("/upload", handlers.UploadHandler)

	return router
}

// Запускаем сервер
func (s *Server) Start() error {
	s.logger.Printf("Запуск сервера на порту %s", s.server.Addr)
	return s.server.ListenAndServe()
}

// Останавливаем сервер
func (s *Server) Stop() error {
	s.logger.Println("Остановка сервера")
	return s.server.Close()
}
