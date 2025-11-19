package main

import (
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"

	"github.com/Badgain/book-discount/internal/handler"
	"github.com/Badgain/book-discount/internal/service"
)

func main() {
	// Получаем порт из переменной окружения или используем 8080 по умолчанию
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// Инициализация сервиса
	discountService := service.NewDiscountService()

	// Инициализация handler
	discountHandler := handler.NewDiscountHandler(discountService)

	// Настройка маршрутов
	http.HandleFunc("/api/v1/discount/calculate", discountHandler.CalculateDiscount)

	// Добавляем health check endpoint
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		if _, err := w.Write([]byte("OK")); err != nil {
			slog.Default().Error("unable to write message to response writer", "error", err.Error())
		}
	})

	// Запуск сервера
	addr := fmt.Sprintf(":%s", port)
	log.Printf("Starting server on %s", addr)
	if err := http.ListenAndServe(addr, nil); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
