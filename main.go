package main

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"time"
)

const (
	schemaSQL = `
CREATE TABLE IF NOT EXISTS payments (
	time timestamp NOT NULL DEFAULT now(),
	amount float NOT NULL DEFAULT 0,
	person text NOT NULL DEFAULT '',
);

CREATE INDEX IF NOT EXISTS payments_person ON payments(person);
`
)

type DB struct {
	db *sql.DB
}

func NewDB(dbfile string) (*DB, error) {
	db, err := sql.Open("sqlite3", dbfile)
	if err != nil {
		return nil, err
	}

	if err = db.Ping(); err != nil {
		return nil, err
	}

	if _, err := db.Exec(schemaSQL); err != nil {
		return nil, err
	}

	return &DB{db: db}, nil
}

func main() {

	// TODO: SQLite

	// HTTP Server
	http.HandleFunc("/payment", paymentHandler)
	http.ListenAndServe(":8080", nil)
}

type Person struct {
	Name string `json:"name"`
}

type Payment struct {
	Person Person    `json:"person"`
	Amount int       `json:"amount"`
	Time   time.Time `json:"date"`
}

func paymentHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "POST":
		postMethod(w, r)
	case "GET":
		getMethod(w, r)
	case "UPDATE":
		updateMethod(w, r)
	case "DELETE":
		deleteMethod(w, r)
	default:
		http.Error(w, "Invalid HTTP method", http.StatusMethodNotAllowed)
		return
	}
}

func postMethod(w http.ResponseWriter, r *http.Request) {
	payment := Payment{}
	if err := json.NewDecoder(r.Body).Decode(&payment); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

}

func getMethod(w http.ResponseWriter, r *http.Request) {

}

func updateMethod(w http.ResponseWriter, r *http.Request) {

}

func deleteMethod(w http.ResponseWriter, r *http.Request) {

}
