package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	_ "modernc.org/sqlite"
)

const (
	schemaSQL = `CREATE TABLE IF NOT EXISTS payments (Id INTEGER PRIMARY KEY, Person TEXT, Amount float, Date TEXT);`
)

var paymentDB *sql.DB

const dbfile = "payments.db"

func NewDB(dbpath string) error {
	var err error
	paymentDB, err = sql.Open("sqlite", dbpath)
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
	Id     int     `json:"id"`
	Person string  `json:"Person"`
	Amount float32 `json:"Amount"`
	Date   string  `json:"Date"`
}

type Person struct {
	Name string `json:"Person"`
}

func paymentHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "POST":
		fmt.Printf("\nPost called\n")
		postMethod(w, r)
	case "GET":
		fmt.Printf("\nGet called\n")
		getMethod(w, r)
	case "PATCH":
		fmt.Printf("\nUpdate called\n")
		updateMethod(w, r)
	case "DELETE":
		fmt.Printf("\nDelete called\n")
		deleteMethod(w, r)
	default:
		http.Error(w, "Invalid HTTP method", http.StatusMethodNotAllowed)
		return
	}
}

func postMethod(w http.ResponseWriter, r *http.Request) {
	payment := Payment{}
	if err := json.NewDecoder(r.Body).Decode(&payment); err != nil {
		fmt.Print("error decoding payment")
		http.Error(w, err.Error(), http.StatusBadRequest)
		fmt.Printf(err.Error())
		return
	}

	println("\n Unmarshaled payment")
	println(payment.Person, payment.Amount, payment.Date)
	println("---\n")

	sql_q := `INSERT INTO payments (person, amount, date) VALUES (?, ?, ?)`

	var result, err = paymentDB.Exec(sql_q, payment.Person, payment.Amount, payment.Date)
	if err != nil {
		println(err.Error())
		return
	}
	println(result.LastInsertId())
	w.WriteHeader(http.StatusCreated)
	var id, _ = result.LastInsertId()
	res := "Payment created with ID: " + strconv.FormatInt(id, 10)
	w.Write([]byte(res))
}

func getMethod(w http.ResponseWriter, r *http.Request) {
	person := Person{}
	if err := json.NewDecoder(r.Body).Decode(&person); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	println(person.Name)

	result, err := paymentDB.Query("SELECT * FROM payments WHERE person = $1", person.Name)
	if err != nil {
		return
	}
	println(result)

	var paymentsQuery Payment
	var paymentsRes []Payment
	for result.Next() {
		var id int64
		err := result.Scan(&paymentsQuery.Id, &paymentsQuery.Person, &paymentsQuery.Amount, &paymentsQuery.Date)
		if err != nil {
			println(err.Error())
			return
		}
		println(id, paymentsQuery.Person, paymentsQuery.Amount, paymentsQuery.Date)
		println("---")
		paymentsRes = append(paymentsRes, paymentsQuery)
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(paymentsRes)
	if err != nil {
		return
	}
}

func updateMethod(w http.ResponseWriter, r *http.Request) {
	payment := Payment{}
	if err := json.NewDecoder(r.Body).Decode(&payment); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		println(err.Error())
		return
	}

	println(payment.Id)

	result, err := paymentDB.Exec("UPDATE payments set Person = $2, Amount = $3, Date = $4 where id = $1", payment.Id, payment.Person, payment.Amount, payment.Date)
	if err != nil {
		println(err.Error())
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	rowsN, _ := result.RowsAffected()
	res := "Updated rows: " + strconv.FormatInt(rowsN, 10)
	w.Write([]byte("Payment updated successfully: " + res))
}

func deleteMethod(w http.ResponseWriter, r *http.Request) {
	payment := Payment{}
	if err := json.NewDecoder(r.Body).Decode(&payment); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		println(err.Error())
		return
	}

	println(payment.Id)

	result, err := paymentDB.Exec("DELETE from payments where Id = $1", payment.Id)
	if err != nil {
		println(err.Error())
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	rowsN, _ := result.RowsAffected()
	res := "Deleted rows: " + strconv.FormatInt(rowsN, 10)
	w.Write([]byte("Payment deleted successfully: " + res))
}
