package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"
)

// DeleteTaskHandler - удаление ненужной задачи
func DeleteTaskHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8") // Установка заголовка ответа для JSON
		// Проверка метода Delete
		if r.Method != http.MethodDelete {
			http.Error(w, "Метод не разрешен", http.StatusMethodNotAllowed)
			return
		}
		// Получение параметра id из URL
		idStr := r.URL.Query().Get("id")
		if idStr == "" { // Проверка, что идентификатор указан
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]string{"error": "Не указан идентификатор"})
			return
		}

		// Проверка на корректность id (можно ли преобразовать в число)
		if _, err := strconv.Atoi(idStr); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]string{"error": "Некорректный идентификатор"})
			return
		}

		// SQL-запрос для удаления задачи по идентификатору
		query := `DELETE FROM scheduler WHERE id = ?`
		_, err := db.Exec(query, idStr) // Запрос к базе данных
		if err != nil {
			writeError(w, "Ошибка обновления задачи", http.StatusInternalServerError)
			return
		}
		// Возвращение пустого JSON {} после успешного удаления задачи
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("{}"))
	}
}
