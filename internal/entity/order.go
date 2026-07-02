package entity

import "errors"

var (
	ErrIDRequired   = errors.New("id is required")
	ErrInvalidPrice = errors.New("invalid price")
	ErrInvalidTax   = errors.New("invalid tax")
)

type Order struct {
	ID         string
	Price      float64
	Tax        float64
	FinalPrice float64
}

func NewOrder(id string, price float64, tax float64) (*Order, error) {
	order := &Order{
		ID:    id,
		Price: price,
		Tax:   tax,
	}
	if err := order.Validate(); err != nil {
		return nil, err
	}
	order.CalculateFinalPrice()
	return order, nil
}

func (o *Order) Validate() error {
	if o.ID == "" {
		return ErrIDRequired
	}
	if o.Price <= 0 {
		return ErrInvalidPrice
	}
	if o.Tax <= 0 {
		return ErrInvalidTax
	}
	return nil
}

func (o *Order) CalculateFinalPrice() {
	o.FinalPrice = o.Price + o.Tax
}

type OrderRepositoryInterface interface {
	Save(order *Order) error
	List() ([]*Order, error)
}
