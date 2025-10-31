package handlers

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/Der1an0/6sprint/internal/service"
)

// IndexHandler возвращает HTML форму
func IndexHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	// Читаем и отдаем HTML файл
	html, err := os.ReadFile("a:/6/6sprint/index.html")
	if err != nil {
		log.Printf("Ошибка чтения index.html: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write(html)
}

// UploadHandler обрабатывает загрузку файла
func UploadHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Парсим форму с файлом
	err := r.ParseMultipartForm(10 << 20) // 10MB max
	if err != nil {
		log.Printf("parsing error: %v", err)
		http.Error(w, "Cannot parse form", http.StatusInternalServerError)
		return
	}

	// Получаем файл из формы
	file, header, err := r.FormFile("myFile")
	if err != nil {
		log.Printf("Ошибка получения файла: %v", err)
		http.Error(w, "Cannot get file", http.StatusInternalServerError)
		return
	}
	defer file.Close()

	// Читаем данные из файла
	fileBytes, err := io.ReadAll(file)
	if err != nil {
		log.Printf("reading error: %v", err)
		http.Error(w, "cannot read file", http.StatusInternalServerError)
		return
	}

	content := string(fileBytes)
	log.Printf("Получен файл: %s, размер: %d байт", header.Filename, len(content))

	// Конвертируем данные
	converted, err := service.Convert(content)
	if err != nil {
		log.Printf("сonversion error: %v", err)
		http.Error(w, fmt.Sprintf("сonversion error: %v", err), http.StatusInternalServerError)
		return
	}

	// Создаем локальный модуль с результатами
	fileName := generateFilename(header.Filename)
	err = os.WriteFile(fileName, []byte(converted), 0755)
	if err != nil {
		log.Printf("file recording error: %v", err)
		http.Error(w, "cannot write result file", http.StatusInternalServerError)
		return
	}
	log.Printf("Результат сохранен в файл: %s", fileName)

	// Возвращаем результат
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Конвертация завершена!\n\nИсходный текст:\n%s\n\nРезультат:\n%s\n\nФайл сохранен: %s",
		content, converted, fileName)
}

// Генерируем имя файла
func generateFilename(original string) string {
	ext := filepath.Ext(original)
	if ext == "" {
		ext = ".txt"
	}
	timestamp := time.Now().UTC().Format("20060102_150405")
	return fmt.Sprintf("result_%s%s", timestamp, ext)
}
