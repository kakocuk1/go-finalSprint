package api

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// Правки по ревью, читаем конфигурацию один раз при старте
// password хранит пароль, считанный один раз при старте
var password string

// InitAuth считывает аутентификацию при старте сервера
func InitAuth() {
	password = os.Getenv("TODO_PASSWORD")
}

// passwordHash возвращает SHA256 хэш пароля в виде строки
func passwordHash(p string) string {
	h := sha256.Sum256([]byte(p))
	return fmt.Sprintf("%x", h)
}

// signInHandler обрабатывает POST-запрос на вход пользователя.
// Он проверяет пароль, и если он верный, создает JWT-токен с временем жизни 1 час и возвращает его в JSON-ответе.
func signInHandler(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Password string `json:"password"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "JSON decode error")
		return
	}

	if req.Password != password {
		writeError(w, http.StatusUnauthorized, "Incorrect password")
		return
	}

	// формируем JWT-токен, кладём хэш пароля
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"hash": passwordHash(password),
		"exp":  time.Now().Add(8 * time.Hour).Unix(), // время жизни куки 8 часов
	})

	tokenStr, err := token.SignedString([]byte(password)) // JWT подписывается самим паролем,как секретным ключом
	if err != nil {
		writeError(w, http.StatusInternalServerError, "Token generation error")
		return
	}

	writeJson(w, map[string]string{"token": tokenStr})
}

// auth - middleware для проверки JWT-токена в куки.
func auth(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if password == "" {
			// если пароль не зада, аутентификация не требуется(отключена)
			next(w, r)
			return
		}

		// читаем JWT-токен из куки
		cookie, err := r.Cookie("token") // Значение куки берётся с помощью r.Cookie("token")
		if err != nil {
			http.Error(w, "Authentification required", http.StatusUnauthorized) // 401 Unauthorized
			return
		}

		// парсим и проверяем JWT токен
		token, err := jwt.Parse(cookie.Value, func(t *jwt.Token) (interface{}, error) {
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("Unexpected signing method: %v", t.Header["alg"])
			}
			return []byte(password), nil
		})
		if err != nil || !token.Valid {
			http.Error(w, "Authentification required", http.StatusUnauthorized)
			return
		}

		// проверяем хэш пароля внутри токена
		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			http.Error(w, "Authentification required", http.StatusUnauthorized)
			return
		}
		if claims["hash"] != passwordHash(password) {
			http.Error(w, "Authentification required", http.StatusUnauthorized)
			return
		}
		// если всё ок, вызываем следующий обработчик
		next(w, r)
	})
}
