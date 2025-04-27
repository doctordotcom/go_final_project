package handlers

import (
	"database/sql"
	"encoding/json"
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
			http.Error(w, "Метод не разрешен", http.StatusMethodNotAllowed)
			return
		}

		// Десериализация JSON-запроса в структуру Task
		var task models.Task
		if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
			writeError(w, "Ошибка десериализации JSON", http.StatusBadRequest)
			return
		}

		// Проверка наличия обязательного поля Title
		if task.Title == "" {
			writeError(w, "Не указан заголовок задачи", http.StatusBadRequest)
			return
		}

		// Получение текущей даты в формате YYYYMMDD
		today := time.Now().Format("20060102")
		todayDate, err := time.Parse("20060102", today)
		if err != nil {
			writeError(w, "Ошибка парсинга текущей даты", http.StatusInternalServerError)
			return
		}

		var taskDate time.Time

		// Проверка и корректировка поля Date
		if task.Date == "" {
			// Если дата не указана, используется сегодняшняя
			task.Date = today
		} else {
			// Проверка корректности формата даты 20060102
			taskDate, err = time.Parse("20060102", task.Date)
			if err != nil {
				writeError(w, "Неправильный формат даты, должно быть YYYYMMDD", http.StatusBadRequest)
				return
			}
			// Если указанная дата меньше текущей
			if taskDate.Format("20060102") < today {
				// Если правило повторения пустое, берётся сегодняшняя дата
				if task.Repeat == "" {
					task.Date = today
				} else {
					// Вычисление следующей даты с учетом правила повторения
					nextDateStr, err := rules.NextDate(todayDate, task.Date, task.Repeat)
					if err != nil {
						writeError(w, "Ошибка вычисления следующей даты", http.StatusBadRequest)
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
			writeError(w, "Ошибка обновления задачи", http.StatusInternalServerError)
			return
		}
		// Проверка, была ли задача найдена и обновлена
		rowsAffected, err := res.RowsAffected()
		if err != nil {
			writeError(w, "Ошибка проверки обновленной задачи", http.StatusInternalServerError)
			return
		}
		if rowsAffected == 0 {
			writeError(w, "Задача не найдена", http.StatusNotFound)
			return
		}
		// Возвращение пустого JSON {}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("{}"))
	}
}
