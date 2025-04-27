package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"m/models"
	"net/http"
	"time"
)

// Максимальное количество задач, возвращаемых в ответе
const maxTasks = 50

// GetTaskListHandler - обработчик для получения списка задач
func GetTaskListHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		// Проверка метода GET
		if r.Method != http.MethodGet {
			writeError(w, "Метод не разрешен", http.StatusMethodNotAllowed)
			return
		}

		// Получение параметра поиска search из запроса
		search := r.URL.Query().Get("search")
		// Ограничение количества возвращаемых задач
		limit := maxTasks
		// Строка SQL-запроса
		var sqlQuery string
		// Слайс аргументов для подготовленного запроса
		var args []interface{}

		// Проверка, является ли search датой и преобразование формата
		if isDate(search) {
			// Преобразование строки даты в нужный формат для сравнения с базой
			formattedDate, _ := time.Parse("02.01.2006", search)
			// Формирование SQL-запроса для поиска задач по конкретной дате
			sqlQuery = "SELECT id, date, title, comment, repeat FROM scheduler WHERE date = ? LIMIT ?"
			// Добавление отформатированной даты и лимита к аргументам
			args = append(args, formattedDate.Format("20060102"), limit)
		} else {
			// Использование LIKE для поиска по заголовку и комментарию
			sqlQuery = "SELECT id, date, title, comment, repeat FROM scheduler WHERE title LIKE ? OR comment LIKE ? ORDER BY date ASC LIMIT ?"
			args = append(args, fmt.Sprintf("%%%s%%", search), fmt.Sprintf("%%%s%%", search), limit)
		}

		// Выполнение подготовленного запроса с аргументами к базе данных
		rows, err := db.Query(sqlQuery, args...)
		if err != nil {
			writeError(w, "Ошибка выполнения запроса к базе данных", http.StatusInternalServerError)
			return
		}
		// Гарантированное закрытие результата после использования
		defer rows.Close()

		// Создание слайса для хранения задач
		tasks := make([]models.Task, 0)
		// Итерация строк результата запроса
		for rows.Next() {
			var task models.Task
			// Считывание значений в структуру Task
			if err := rows.Scan(&task.ID, &task.Date, &task.Title, &task.Comment, &task.Repeat); err != nil {
				writeError(w, "Ошибка считывания данных задачи", http.StatusInternalServerError)
				return
			}
			// Добавление задач в слайс
			tasks = append(tasks, task)
		}

		// Проверка на наличие ошибок в процессе итерации строк
		if err := rows.Err(); err != nil {
			writeError(w, "Ошибка при переборе задач", http.StatusInternalServerError)
			return
		}

		// Установка заголовка Content-Type для ответа
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		// Сериализация слайса задач в JSON и отправка его в ответе
		if err := json.NewEncoder(w).Encode(map[string]interface{}{"tasks": tasks}); err != nil {
			writeError(w, "Ошибка кодирования JSON", http.StatusInternalServerError)
			return
		}
	}
}

// isDate - вспомогательная функция, проверяющая, является ли строка датой в формате 02.01.2006
func isDate(dateStr string) bool {
	_, err := time.Parse("02.01.2006", dateStr)
	return err == nil
}
