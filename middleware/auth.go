package middleware

import (
	"fmt"
	"net/http"
	"os"

	"github.com/golang-jwt/jwt/v4"
)

// Секретный ключ для подписи и проверки JWT-токенов
var secretKey = []byte("secret_key")

// AuthMiddleware - функция для проверки аутентификации пользователя через JWT-токен
func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		// Получение пароля из переменной окружения TODO_PASSWORD
		password := os.Getenv("TODO_PASSWORD")

		// Если переменная окружения пуста, пропускается проверка и вызывается следующий обработчик
		if password == "" {
			next.ServeHTTP(w, r)
			return
		}

		// Получение cookie с именем "token" из запроса
		cookie, err := r.Cookie("token")
		if err != nil {
			http.Error(w, `{"error": "unauthorized"}`, http.StatusUnauthorized)
			return
		}

		// Получение строки токена из значения cookie
		tokenString := cookie.Value
		// Парсинг JWT-токена с помощью функции проверки подписи
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			// Проверка метода подписи — должен быть HMAC
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, http.ErrNoLocation
			}
			// Возвращается секретный ключ для проверки подписи
			return secretKey, nil
		})

		// Проверка результата парсинга и валидности токена
		if err != nil || !token.Valid {
			http.Error(w, `{"error": "unauthorized"}`, http.StatusUnauthorized)
			return
		}

		// Получение claims из токена в виде MapClaims
		m := token.Claims.(jwt.MapClaims)
		// Проверка, что хеш в токене совпадает с хешем из пароля
		if m["hash"] != fmt.Sprintf("%x", password) {
			http.Error(w, `{"error": "unauthorized"}`, http.StatusUnauthorized)
			return
		}

		// Если все проверки прошли успешно, управление передаётся следующему обработчику
		next.ServeHTTP(w, r)
	})
}
