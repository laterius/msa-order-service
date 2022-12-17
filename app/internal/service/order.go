package service

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

const Host = "http://localhost:8000"

type Service struct {
	db *gorm.DB
}

func NewService(db *gorm.DB) Service {
	return Service{db: db}
}

//Реализация методов обращения в базу данных

type Reservations struct {
	OrderId uuid.UUID `json:"orderId" gorm:"type:uuid; not null"`
	GoodId  uuid.UUID `json:"goodId" gorm:"type:uuid; not null"`
}

type Good struct {
	ID     uuid.UUID `json:"id"`
	Amount int       `json:"amount"`
}

type Order struct {
	Id uuid.UUID `json:"id" gorm:"type:uuid; unique; primary_key;"`
}

//type ID struct {
//	value string
//}
//
//func (v *ID) GetValue() string {
//	return v.value
//}
//
//func createID() ID {
//	value := uuid.NewString()
//
//	return ID{
//		value,
//	}
//}

// CreateOrder returns new Order
func (s *Service) CreateOrder() Order {
	return Order{
		uuid.New(),
	}
}

func (s *Service) Store(order Order) error {
	err := s.db.Create(Order{
		Id: order.Id,
	}).Error

	if err != nil {
		return err
	}

	return nil
}

func (s *Service) Delete(order Order) error {
	return s.db.Delete(&Order{}, order.Id).Error
}
