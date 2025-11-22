package middlewares

import (
	"bytes"
	"io"
	"log"
	"net/http"
)

func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("=> %s %s\n", r.Method, r.URL.Path) // логируем метод и путь
		bodyBytes, _ := io.ReadAll(r.Body)
		r.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
		log.Println(string(bodyBytes))
		next.ServeHTTP(w, r) // вызываем следующий обработчик
		log.Println("<= responded")
	})
}
