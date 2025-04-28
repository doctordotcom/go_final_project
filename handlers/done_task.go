package handlers

import (
	"database/sql"
	"encoding/json"
	"m/models"
	"m/rules"
	"net/http"
	"strconv"
	"time"
)

// DoneTaskHandler - отмечает задачу выполненной и переносит, либо удаляет её
func DoneTaskHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		// Проверка метода Post
		if r.Method != http.MethodPost {
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
		err := db.QueryRow(sqlQuery, idStr).Scan(&task.ID, &task.Date, &task.Title, &task.Comment, &task.Repeat)
		if err != nil {
			if err == sql.ErrNoRows {
				// Если задача не найдена, возвращается ошибка
				w.WriteHeader(http.StatusNotFound)
				json.NewEncoder(w).Encode(map[string]string{"error": "Issue not found"})
			} else {
				// Ошибка базы данных
				writeError(w, "Database query execution error", http.StatusInternalServerError)
			}
			return
		}

		// Разовая задача
		if task.Repeat == "" {
			q := `DELETE FROM scheduler WHERE id = ?`
			_, err := db.Exec(q, task.ID)
			if err != nil {
				writeError(w, "Issue update error", http.StatusInternalServerError)
				return
			}
			// Возвращение пустого JSON {} после удаления разовой задачи
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("{}"))
			return
		}

		// Периодическая задача
		nextDate, err := rules.NextDate(time.Now(), task.Date, task.Repeat)
		if err != nil {
			writeError(w, "Date error: "+err.Error(), http.StatusInternalServerError)
		}
		task.Date = nextDate
		// Проверка существования задачи и обновление её
		query := `UPDATE scheduler SET date = ?, title = ?, comment = ?, repeat = ? WHERE id = ?`
		res, err := db.Exec(query, task.Date, task.Title, task.Comment, task.Repeat, task.ID)
		if err != nil {
			writeError(w, "Issue update error", http.StatusInternalServerError)
			return
		}
		// Проверка, была ли задача найдена и обновлена
		rowsAffected, err := res.RowsAffected()
		if err != nil {
			writeError(w, "Error checking the updated issue", http.StatusInternalServerError)
			return
		}
		if rowsAffected == 0 {
			writeError(w, "Issue not found", http.StatusNotFound)
			return
		}
		// Возвращение пустого JSON {}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("{}"))
	}
}
