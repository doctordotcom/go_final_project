package handlers

import (
	"database/sql"
	"encoding/json"
	"m/global"
	"m/models"
	"m/rules"
	"net/http"
	"time"
)

// PutTaskHandler - обработчик для передачи обновлённых полей задачи после её редактирования
func PutTaskHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")

		// Проверка метода PUT
		if r.Method != http.MethodPut {
			http.Error(w, "Method is not allowed", http.StatusMethodNotAllowed)
			return
		}

		// Десериализация JSON-запроса в структуру Task
		var task models.Task
		if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
			writeError(w, "Error deserializing JSON", http.StatusBadRequest)
			return
		}

		// Проверка наличия обязательного поля Title
		if task.Title == "" {
			writeError(w, "The issue title is not specified", http.StatusBadRequest)
			return
		}

		// Получение текущей даты в формате YYYYMMDD
		today := time.Now().Format(global.DateFormat)
		todayDate, err := time.Parse(global.DateFormat, today)
		if err != nil {
			writeError(w, "Current date parsing error", http.StatusInternalServerError)
			return
		}

		var taskDate time.Time

		// Проверка и корректировка поля Date
		if task.Date == "" {
			// Если дата не указана, используется сегодняшняя
			task.Date = today
		} else {
			// Проверка корректности формата даты 20060102
			taskDate, err = time.Parse(global.DateFormat, task.Date)
			if err != nil {
				writeError(w, "The date format is incorrect, it must be YYYYMMDD", http.StatusBadRequest)
				return
			}
			// Если указанная дата меньше текущей
			if taskDate.Format(global.DateFormat) < today {
				// Если правило повторения пустое, берётся сегодняшняя дата
				if task.Repeat == "" {
					task.Date = today
				} else {
					// Вычисление следующей даты с учетом правила повторения
					nextDateStr, err := rules.NextDate(todayDate, task.Date, task.Repeat)
					if err != nil {
						writeError(w, "Error calculating the next date", http.StatusBadRequest)
						return
					}
					task.Date = nextDateStr
				}
			}
		}
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
