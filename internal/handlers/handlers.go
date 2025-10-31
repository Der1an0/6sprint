package handlers

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"time"

	"github.com/Der1an0/6sprint/internal/service"
)

// IndexHandler возвращает HTML форму
func IndexHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	// Пробуем найти index.html в разных местах
	var html []byte
	var err error

	// Сначала пробуем текущую директорию
	html, err = os.ReadFile("index.html")
	if err != nil {
		// Пробуем абсолютный путь
		_, filename, _, _ := runtime.Caller(0)
		dir := filepath.Dir(filename)
		possiblePaths := []string{
			"index.html",
			filepath.Join(dir, "index.html"),
			filepath.Join(dir, "..", "index.html"),
			filepath.Join(dir, "..", "..", "index.html"),
			filepath.Join(".", "index.html"),
		}

		for _, path := range possiblePaths {
			html, err = os.ReadFile(path)
			if err == nil {
				break
			}
		}
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
