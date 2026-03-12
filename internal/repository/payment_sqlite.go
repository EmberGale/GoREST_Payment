package repository

import (
	"GoREST_Payment/internal/model"
	"database/sql"
)

type sqlitePaymentRepository struct {
	db *sql.DB
}

func NewPaymentRepository(db *sql.DB) PaymentRepository {
	return &sqlitePaymentRepository{db: db}
}

func (r *sqlitePaymentRepository) Create(payment *model.Payment) (int64, error) {
	const query = `INSERT INTO payments (person, amount, date) VALUES (?, ?, ?)`

	result, err := r.db.Exec(query, payment.Person, payment.Amount, payment.Date)
	if err != nil {
		return 0, err
	}

	return result.LastInsertId()
}

func (r *sqlitePaymentRepository) GetByPerson(name string) ([]model.Payment, error) {
	const query = `SELECT id, person, amount, date FROM payments WHERE person = ?`

	rows, err := r.db.Query(query, name)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var payments []model.Payment
	for rows.Next() {
		var p model.Payment
		if err := rows.Scan(&p.Id, &p.Person, &p.Amount, &p.Date); err != nil {
			return nil, err
		}
		payments = append(payments, p)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return payments, nil
}

func (r *sqlitePaymentRepository) Update(payment *model.Payment) (int64, error) {
	const query = `UPDATE payments SET person = ?, amount = ?, date = ? WHERE id = ?`

	result, err := r.db.Exec(query, payment.Person, payment.Amount, payment.Date, payment.Id)
	if err != nil {
		return 0, err
	}

	return result.RowsAffected()
}

func (r *sqlitePaymentRepository) Delete(id int) (int64, error) {
	const query = `DELETE FROM payments WHERE id = ?`

	result, err := r.db.Exec(query, id)
	if err != nil {
		return 0, err
	}

	return result.RowsAffected()
}
