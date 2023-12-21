package main

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

func BenchmarkMoveCourierHandler(b *testing.B) {
	// Создаем экземпляр Gin-роутера
	router := gin.Default()

	// Добавляем обработчик POST "/move-courier"
	router.POST("/move-courier", func(c *gin.Context) {
		// Эмулируем JSON-запрос
		c.Request, _ = http.NewRequest("POST", "/move-courier", nil)

		// Эмулируем JSON-тело запроса
		c.Request.Header.Set("Content-Type", "application/json")
		c.Request.Body = http.NoBody

		// Вызываем обработчик
		router.ServeHTTP(c.Writer, c.Request)
	})

	// Эмулируем запросы
	req, _ := http.NewRequest("POST", "/move-courier", nil)

	// Запускаем бенчмарк
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		// Создаем HTTP запрос
		rec := httptest.NewRecorder()
		router.ServeHTTP(rec, req)
	}
}
