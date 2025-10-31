package main

import (
	"log"
	"os"

	"github.com/Der1an0/6sprint/internal/server"
)

func main() {
	// Создаем логгер
	logger := log.New(os.Stdout, "MORSE_SERVICE: ", log.Ldate|log.Ltime|log.Lshortfile)

	// Создаем сервер
	srv := server.New(logger)

	// Запускаем сервер
	if err := srv.Start(); err != nil {
		logger.Fatalf("Server startup error: %v", err)
	}
}
