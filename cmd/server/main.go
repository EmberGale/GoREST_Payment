package main

import (
	"database/sql"
	"fmt"
	"net/http"

	"GoREST_Payment/internal/handler"
	"GoREST_Payment/internal/repository"
	"GoREST_Payment/internal/service"

	_ "github.com/glebarez/go-sqlite"
)

const (
	schemaSQL = `CREATE TABLE IF NOT EXISTS payments (
		Id INTEGER PRIMARY KEY,
		Person TEXT,
		Amount REAL,
		Date TEXT
	);`
	dbFile = "payments.db"
)

func initDB(dbPath string) (*sql.DB, error) {
	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	if _, err := db.Exec(schemaSQL); err != nil {
		return nil, err
	}

	return db, nil
}

func main() {
	// Инициализация БД
	db, err := initDB(dbFile)
	if err != nil {
		panic(fmt.Sprintf("Failed to initialize database: %v", err))
	}
	defer db.Close()

	// Внедрение зависимостей: DB → repository → service → handler → router
	paymentRepo := repository.NewPaymentRepository(db)
	paymentService := service.NewPaymentService(paymentRepo)
	paymentHandler := handler.NewPaymentHandler(paymentService)
	router := handler.NewRouter(paymentHandler)

	// Запуск HTTP-сервера
	fmt.Println("Server starting on :8080")
	if err := http.ListenAndServe(":8080", router); err != nil {
		panic(fmt.Sprintf("Failed to start server: %v", err))
	}
}
