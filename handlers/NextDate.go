package handlers

import (
	"m/rules"
	"net/http"
	"time"
)

// GetNextDateHandler - обработчик для эндпоинта "/api/nextdate" (получение следующей даты)
func GetNextDateHandler(w http.ResponseWriter, r *http.Request) {

	// Проверка метода Get
	if r.Method != http.MethodGet {
		http.Error(w, "Метод не разрешен", http.StatusMethodNotAllowed)
		return
	}

	// Получение параметров URL-запроса
	nowStr := r.URL.Query().Get("now")
	dateStr := r.URL.Query().Get("date")
	repeat := r.URL.Query().Get("repeat")

	// Проверка наличия всех обязательных параметров
	if nowStr == "" || dateStr == "" || repeat == "" {
		http.Error(w, "Параметры now, date и repeat обязательны", http.StatusBadRequest)
		return
	}

	// Парсинг now (преобразование строки now в объект времени в соотетствии с форматом "20060102")
	Now, err := time.Parse("20060102", nowStr)
	if err != nil {
		http.Error(w, "Ошибка парсинга параметра now", http.StatusBadRequest)
		return
	}

	// Парсинг dateStr (преобразование строки date в объект времени)
	Date, err := time.Parse("20060102", dateStr)
	if err != nil {
		http.Error(w, "Ошибка парсинга параметра date", http.StatusBadRequest)
		return
	}

	// Вычисление следующей даты с использованием текущей даты, исходной даты и правил повторения
	NextDate, err := rules.NextDate(Now, Date.Format("20060102"), repeat)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Преобразование даты в байтовый массив и вывод в ответ
	w.Write([]byte(NextDate))
}
