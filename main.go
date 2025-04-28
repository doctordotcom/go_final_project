package main

import (
	"fmt"
	"m/handlers"
	"m/middleware"
	"m/scheduler"
	"net/http"
	"os"

	"github.com/joho/godotenv"

	_ "modernc.org/sqlite"
)

func main() {
	// Загрузка переменных окружения из .env файла
	err := godotenv.Load()
	if err != nil {
		// log.Fatal("Error loading .env file")
		fmt.Println("Error loading .env file")
		os.Exit(1)
	}

	// Открытие соединения с базой данных
	db, err := scheduler.OpenDatabase()
	if err != nil {
		// log.Fatal(err)
		fmt.Println("Error opening database:", err)
		os.Exit(1)
	}

	// Закрытие соединения после работы с базой данных
	defer db.Close()

	// Обработка запросов на вход в систему
	http.Handle("/api/signin", http.HandlerFunc(handlers.SignInHandler))

	// Обработка запросов на получение следующей даты
	http.HandleFunc("/api/nextdate", handlers.GetNextDateHandler)

	// Обработка запросов на управление задачами с авторизацией
	http.Handle("/api/task", middleware.AuthMiddleware(http.HandlerFunc(handlers.SwitcherHandler(db))))
	http.Handle("/api/task/done", middleware.AuthMiddleware(http.HandlerFunc(handlers.DoneTaskHandler(db))))
	http.Handle("/api/tasks", middleware.AuthMiddleware(http.HandlerFunc(handlers.GetTaskListHandler(db))))

	// Указание на директорию для статических файлов
	webDir := "web"

	// Обработка запросов к статическим файлам
	http.Handle("/", http.FileServer(http.Dir(webDir)))
	fmt.Println("Сервер запущен")

	// Получение порта из переменных окружения
	port := os.Getenv("TODO_PORT")
	if port == "" {
		port = "7540"
	}

	// Формирование адреса для сервера
	address := fmt.Sprintf(":%s", port)

	// Логирование на всякий случай
	// log.Printf("Starting server on port %s...\n", port)
	fmt.Printf("Starting server on port %s...\n", port)

	// Запуска сервера и логирование ошибки, если она возникнет
	// log.Fatal(http.ListenAndServe(address, nil))
	fmt.Println(http.ListenAndServe(address, nil))

}
