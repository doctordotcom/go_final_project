package handlers

import (
	"database/sql"
	"encoding/json"
	"m/models"
	"net/http"
	"strconv"
)

// GetTaskHandler - обработчик для получения полей задачи перед её редактированием
func GetTaskHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")

		// Проверка, что метод запроса - GET
		if r.Method != http.MethodGet {
			http.Error(w, "Method is not allowed", http.StatusMethodNotAllowed)
			return
		}
		// Получение параметра id из URL
		idStr := r.URL.Query().Get("id")
		if idStr == "" {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]string{"error": "ID not specified"})
			return
		}

		// Проверка на корректный id
		if _, err := strconv.Atoi(idStr); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]string{"error": "Invalid ID"})
			return
		}

		// Подготовка SQL-запроса
		sqlQuery := "SELECT id, date, title, comment, repeat FROM scheduler WHERE id = ?"
		var task models.Task

		// Выполнение запроса
		// Метод QueryRow выполняет запрос и ожидает ровно одну строку результата
		// Метод Scan связывает поля результата с соответствующими полями структуры task
		err := db.QueryRow(sqlQuery, idStr).Scan(&task.ID, &task.Date, &task.Title, &task.Comment, &task.Repeat)
		if err != nil {
			// Проверка на случай, если задача не найдена в базе данных
			if err == sql.ErrNoRows {
				// Установка HTTP-статуса "Not Found" и формирование ответа с сообщением об ошибке
				w.WriteHeader(http.StatusNotFound)
				json.NewEncoder(w).Encode(map[string]string{"error": "Issue not found"})
			} else {
				// Вызов вспомогательной функции для записи ошибки
				writeError(w, "Ошибка выполнения запроса к базе данных", http.StatusInternalServerError)
			}
			return
		}

		// Если задача успешно найдена, формируется ответ с данными
		w.WriteHeader(http.StatusOK)
		// Сериализация объекта task в JSON и отправка в ответ
		json.NewEncoder(w).Encode(task)
	}
}
