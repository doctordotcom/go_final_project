package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

// Определение секретноого ключа для подписи
var secretKey = []byte("secret_key")

// SignInHandler - обработчик для входа в систему
func SignInHandler(w http.ResponseWriter, r *http.Request) {
	var body struct {
		Password string `json:"password"` // Структура для хранения пароля из запроса
	}
	// Декодирование JSON-данных из тела запроса в структуру body
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, `{"error": "invalid request"}`, http.StatusBadRequest)
		return
	}

	// Получение ожидаемого пароля из переменных окружения
	expectedPassword := os.Getenv("TODO_PASSWORD")
	// Сравнение введенного пароля с ожидаемым
	if body.Password != expectedPassword {
		http.Error(w, `{"error": "invalid password"}`, http.StatusUnauthorized)
		return
	}

	// Хеширование пароля для создания токена
	hash := fmt.Sprintf("%x", body.Password)

	// Создание нового JWT с хешем пароля и временем истечения
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"hash": hash,                                 // хеш пароля добавлен в токен
		"exp":  time.Now().Add(8 * time.Hour).Unix(), // время жизни токена 8 часов
	})

	// Подписание токена с использованием секретного ключа
	tokenString, err := token.SignedString(secretKey)
	if err != nil {
		http.Error(w, `{"error": "could not create token"}`, http.StatusInternalServerError)
		return
	}

	// Установка заголовка ответа для JSON
	w.Header().Set("Content-Type", "application/json")
	// Кодирование токена в JSON и отправка его в ответе
	json.NewEncoder(w).Encode(map[string]string{"token": tokenString})
}
