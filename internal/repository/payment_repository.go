package repository

import "GoREST_Payment/internal/model"

type PaymentRepository interface {
	Create(payment *model.Payment) (int64, error)
	GetByPerson(name string) ([]model.Payment, error)
	Update(payment *model.Payment) (int64, error)
	Delete(id int) (int64, error)
}
