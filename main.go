package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	_ "modernc.org/sqlite"
)

const (
	schemaSQL = `CREATE TABLE IF NOT EXISTS payments (id INTEGER PRIMARY KEY, person TEXT, Amount float, date datetime);`
)

var paymentDB *sql.DB

const dbfile = "/payments.db"

func NewDB(dbpath string) error {
	paymentDB, err := sql.Open("sqlite", "payments.db")
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

	fmt.Printf("%+v", time.Time{})
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
	Person string `json:"person"`
	Amount int    `json:"amount"`
	Time   string `json:"date"`
}

type Person struct {
	Name string `json:"name"`
}

func paymentHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "POST":
		fmt.Printf("Post called\n")
		postMethod(w, r)
	case "GET":
		fmt.Printf("Get called\n")
		getMethod(w, r)
	case "UPDATE":
		fmt.Printf("Update called\n")
		updateMethod(w, r)
	case "DELETE":
		fmt.Printf("Delete called\n")
		deleteMethod(w, r)
	default:
		http.Error(w, "Invalid HTTP method", http.StatusMethodNotAllowed)
		return
	}
}

func postMethod(w http.ResponseWriter, r *http.Request) {
	payment := Payment{}
	fmt.Printf("%+v", r.Body)
	if err := json.NewDecoder(r.Body).Decode(&payment); err != nil {
		fmt.Print("error decoding payment")
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	println(payment.Person, payment.Amount, payment.Time)
	fmt.Printf("%+v", payment)
	_, err := paymentDB.Exec("INSERT INTO payments (time, amount, person) VALUES ($1, $2, $3)", payment.Time, payment.Amount, payment.Person)
	if err != nil {
		return
	}
	w.WriteHeader(http.StatusCreated)
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

	var paymentsQuery Payment
	var paymentsRes []Payment
	for result.Next() {
		result.Scan(&paymentsQuery)
		paymentsRes = append(paymentsRes, paymentsQuery)
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(paymentsRes)
	if err != nil {
		return
	}
}

func updateMethod(w http.ResponseWriter, r *http.Request) {

}

func deleteMethod(w http.ResponseWriter, r *http.Request) {

}
