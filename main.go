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

var paymentDB *sql.DB

const dbfile = "payments.db"

func NewDB(dbpath string) error {
	paymentDB, err := sql.Open("sqlite3", dbpath)
	if err != nil {
		return err
	}

	if err = paymentDB.Ping(); err != nil {
		return err
	}

	if _, err := paymentDB.Exec(schemaSQL); err != nil {
		return err
	}

	return nil
}

func main() {

	// TODO: SQLite
	err := NewDB(dbfile)
	if err != nil {
		panic(err)
	}

	// HTTP Server
	http.HandleFunc("/payment", paymentHandler)
	err = http.ListenAndServe(":8080", nil)
	if err != nil {
		return
	}
}

type Payment struct {
	Person string    `json:"person"`
	Amount int       `json:"amount"`
	Time   time.Time `json:"date"`
}

type Person struct {
	Name string `json:"name"`
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
	_, err := paymentDB.Exec("INSERT INTO payments (time, amount, person) VALUES ($1, $2, $3)", payment.Time, payment.Amount, payment.Person)
	if err != nil {
		return
	}

}

func getMethod(w http.ResponseWriter, r *http.Request) {
	person := Person{}
	if err := json.NewDecoder(r.Body).Decode(&person); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	result, err := paymentDB.Query("SELECT * FROM payments WHERE person = $1", person.Name)
	if err != nil {
		return
	}

	defer func(result *sql.Rows) {
		err := result.Close()
		if err != nil {
			return
		}
	}(result)

	r.Body = json.NewEncoder(result)
}

func updateMethod(w http.ResponseWriter, r *http.Request) {

}

func deleteMethod(w http.ResponseWriter, r *http.Request) {

}
