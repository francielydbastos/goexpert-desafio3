package entity

import "testing"

func TestNewOrder(t *testing.T) {
	order, err := NewOrder("1", 100.0, 10.0)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if order.FinalPrice != 110.0 {
		t.Errorf("expected final price 110, got %v", order.FinalPrice)
	}
}

func TestNewOrderValidation(t *testing.T) {
	if _, err := NewOrder("", 100.0, 10.0); err != ErrIDRequired {
		t.Errorf("expected ErrIDRequired, got %v", err)
	}
	if _, err := NewOrder("1", 0, 10.0); err != ErrInvalidPrice {
		t.Errorf("expected ErrInvalidPrice, got %v", err)
	}
	if _, err := NewOrder("1", 100.0, 0); err != ErrInvalidTax {
		t.Errorf("expected ErrInvalidTax, got %v", err)
	}
}
