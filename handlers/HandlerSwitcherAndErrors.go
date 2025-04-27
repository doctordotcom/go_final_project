package handlers

import (
	"database/sql"
	"fmt"
	"net/http"
)

// SwitcherHandler - переключатель между хендлерами
func SwitcherHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			PostTaskHandler(db)(w, r)
		case http.MethodGet:
			GetTaskHandler(db)(w, r)
		case http.MethodPut:
			PutTaskHandler(db)(w, r)
		case http.MethodDelete:
			DeleteTaskHandler(db)(w, r)
		default:
			http.Error(w, "Метод не разрешен", http.StatusMethodNotAllowed)
		}
	}
}

// Код обработки ошибок для вывода в верном формате (совет от наставника)
func writeError(w http.ResponseWriter, message string, code int) {
	http.Error(w, fmt.Sprintf(`{"error":"%s"}`, message), code)
}
