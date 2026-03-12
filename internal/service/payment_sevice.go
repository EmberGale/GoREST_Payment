package service

import (
	"GoREST_Payment/internal/model"
	"GoREST_Payment/internal/repository"
)

type PaymentService interface {
	CreatePayment(payment *model.Payment) (int64, error)
	GetPaymentsByPerson(name string) ([]model.Payment, error)
	UpdatePayment(payment *model.Payment) (int64, error)
	DeletePayment(id int) (int64, error)
}

type paymentService struct {
	repo repository.PaymentRepository
}

func NewPaymentService(repo repository.PaymentRepository) PaymentService {
	return &paymentService{repo: repo}
}

func (s *paymentService) CreatePayment(payment *model.Payment) (int64, error) {
	return s.repo.Create(payment)
}

func (s *paymentService) GetPaymentsByPerson(name string) ([]model.Payment, error) {
	return s.repo.GetByPerson(name)
}

func (s *paymentService) UpdatePayment(payment *model.Payment) (int64, error) {
	return s.repo.Update(payment)
}

func (s *paymentService) DeletePayment(id int) (int64, error) {
	return s.repo.Delete(id)
}
