package service

import (
	"errors"
	"github.com/google/uuid"
	"github.com/laterius/service_architecture_hw3/app/internal/domain"
	"github.com/laterius/service_architecture_hw3/app/internal/repo"
)

var (
	ErrOrderAlreadyCreated = errors.New("order with this idempotence key is already created")
	//ErrKeyAlreadyExists    = errors.New("idempotence key is already exist")
	//success                = false
	IdempotenceKeys = make(map[string]bool)
)

type OrderService interface {
	OrderReader
	OrderCreator
}

func NewOrderService(repo repo.OrderRepo, keyRepo repo.IdempotenceKeyRepo) *orderService {
	return &orderService{
		reader:             repo,
		creator:            repo,
		idempotenceKeyRepo: keyRepo,
	}
}

type orderService struct {
	reader             repo.OrderReader
	creator            repo.OrderCreator
	idempotenceKeyRepo repo.IdempotenceKeyRepo
}

type OrderReader interface {
	Get() ([]*domain.Order, error)
}

type OrderCreator interface {
	Create(idempotenceKey uuid.UUID, creator *domain.Order) error
}

func (s *orderService) Get() ([]*domain.Order, error) {

	orders, err := s.reader.Get()
	if err != nil {
		if err.Error() == "record not found" {
			return nil, domain.ErrOrderNotFound
		}
	}

	return orders, err
}

func (s *orderService) Create(idempotenceKey uuid.UUID, order *domain.Order) error {

	key, err := s.idempotenceKeyRepo.Get(idempotenceKey)
	if err != nil {
		if err.Error() != "record not found" {
			return err
		}

	}

	if key.ID == idempotenceKey {
		return ErrOrderAlreadyCreated
	}

	err = s.idempotenceKeyRepo.Create(&domain.IdempotenceKey{
		ID:     idempotenceKey,
		Status: domain.OrderStatusCreated,
	})
	if err != nil {
		return err
	}

	err = s.creator.Create(order)
	if err != nil {
		return err
	}

	return nil
}

type OrderData struct {
	UserID int `json:"userId" schema:"userId"`
	Status int `json:"status" schema:"status"`
	Amount int `json:"amount" schema:"total_amount"`
}

type Order struct {
	Id string `json:"id"`
	OrderData
}

func (order *Order) FromDomain(d *domain.Order) *Order {
	order.Id = d.ID.String()
	order.UserID = d.UserID
	order.Status = int(d.Status)
	order.Amount = d.Amount

	return order
}

func (o *Order) ToDomain() *domain.Order {
	id, err := uuid.Parse(o.Id)
	if err != nil {
		panic(err)
	}

	order := &domain.Order{
		ID:     id,
		UserID: o.UserID,
		Status: domain.OrderStatus(o.Status),
		Amount: o.Amount,
	}

	return order
}
