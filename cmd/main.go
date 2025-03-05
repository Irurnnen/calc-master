package main

import (
	"log"
	"net/http"

	"github.com/Irurnnen/calc-master/internal/handlers"
)

// Глобальные переменные

func main() {

	// Регистрируем HTTP-обработчики
	http.HandleFunc("/api/v1/calculate", handlers.HandleCalculate)
	http.HandleFunc("/api/v1/expressions", handlers.HandleGetExpressions)
	http.HandleFunc("/api/v1/expressions/", handlers.HandleGetExpressionByID)
	http.HandleFunc("/internal/task", handlers.HandleTask)

	// Запускаем сервер
	log.Println("Starting server on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
