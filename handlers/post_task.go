package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"m/global"
	"m/models"
	"m/rules"
	"net/http"
	"time"
)

// PostTaskHandler - обработчик для добавления задачи
func PostTaskHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		if r.Method != http.MethodPost {
			http.Error(w, "Method is not allowed", http.StatusMethodNotAllowed)
			return
		}
		// Десериализация JSON-запроса в структуру Task
		var task models.Task
		if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
			writeError(w, "Error deserializing JSON", http.StatusBadRequest)
			return
		}

		// Проверка обязательного поля Title
		if task.Title == "" {
			writeError(w, "The issue title is not specified", http.StatusBadRequest)
			return
		}

		// Получение текущей даты
		today := time.Now().Format(global.DateFormat)
		todayDate, err := time.Parse(global.DateFormat, today)
		if err != nil {
			writeError(w, "Current date parsing error", http.StatusInternalServerError)
			return
		}

		var taskDate time.Time
		// Проверка даты и установка
		if task.Date == "" {
			// Если дата не указана, используем сегодняшнюю
			task.Date = today
		} else {
			// Проверка формата 20060102
			var err error
			taskDate, err = time.Parse(global.DateFormat, task.Date)
			if err != nil {
				writeError(w, "The date format is incorrect, it must be YYYYMMDD", http.StatusBadRequest)
				return
			}

			// Если дата меньше сегодня
			if taskDate.Format(global.DateFormat) < today {
				// Если правило повторения пустое, берётся сегодняшняя дата
				if task.Repeat == "" {
					task.Date = today
				} else {
					// Функция NextDate вычисляет следующую дату с учетом repeat
					nextDateStr, err := rules.NextDate(todayDate, task.Date, task.Repeat)
					if err != nil {
						writeError(w, "Error calculating the next date", http.StatusBadRequest)
						return
					}
					task.Date = nextDateStr
				}
			}
		}

		// Добавление новой задачи в таблицу базы данных
		query := `INSERT INTO scheduler (date, title, comment, repeat) VALUES (?, ?, ?, ?)`
		res, err := db.Exec(query, task.Date, task.Title, task.Comment, task.Repeat)
		if err != nil {
			writeError(w, "Error adding an issue", http.StatusInternalServerError)
			return
		}

		// Получение ID созданной записи
		id, err := res.LastInsertId()
		if err != nil {
			writeError(w, "Couldn't get a new issue ID", http.StatusInternalServerError)
			return
		}

		// Возвращение ответа с кодом 201 (Created) в формате JSON
		w.WriteHeader(http.StatusCreated)
		fmt.Fprintf(w, `{"id":"%d"}`, id)
	}
}
