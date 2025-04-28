package handlers

import (
	"database/sql"
	"net/http"
	"strconv"
)

// DeleteTaskHandler - удаление ненужной задачи
func DeleteTaskHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8") // Установка заголовка ответа для JSON

		// Проверка метода Delete
		if r.Method != http.MethodDelete {
			writeError(w, "Method is not allowed", http.StatusMethodNotAllowed)
			return
		}
		// Получение параметра id из URL
		idStr := r.URL.Query().Get("id")
		if idStr == "" { // Проверка, что идентификатор указан
			writeError(w, "ID not specified", http.StatusBadRequest)
			return
		}

		// Проверка на корректность id (можно ли преобразовать в число)
		if _, err := strconv.Atoi(idStr); err != nil {
			writeError(w, "Invalid ID", http.StatusBadRequest)
			return
		}

		// SQL-запрос для удаления задачи по идентификатору
		query := `DELETE FROM scheduler WHERE id = ?`
		result, err := db.Exec(query, idStr) // Запрос к базе данных
		if err != nil {
			writeError(w, "Issue update error", http.StatusInternalServerError)
			return
		}

		rowsAffected, err := result.RowsAffected()
		if err != nil {
			writeError(w, "Failed to get affected rows", http.StatusInternalServerError)
			return
		}

		if rowsAffected == 0 {
			writeError(w, "Task not found", http.StatusNotFound)
			return
		}

		// Возвращение пустого JSON {} после успешного удаления задачи
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("{}"))
	}
}
