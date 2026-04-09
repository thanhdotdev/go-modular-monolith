package orderapplication

import (
	"context"
	"errors"
	orderdomain "project-example/internal/modules/order/domain"
	"testing"
)

type fakeRepository struct {
	order *orderdomain.Order
	err   error
}

func (f fakeRepository) FindByID(_ context.Context, _ string) (*orderdomain.Order, error) {
	if f.err != nil {
		return nil, f.err
	}

	return f.order, nil
}

func TestGetOrder(t *testing.T) {
	service := NewService(fakeRepository{
		order: &orderdomain.Order{
			ID:           "ord-001",
			CustomerName: "Alice",
			Status:       "pending",
			TotalAmount:  125000,
		},
	})

	got, err := service.GetOrder(context.Background(), "ord-001")
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}

	if got.ID != "ord-001" {
		t.Fatalf("expected order id ord-001, got %s", got.ID)
	}
}

func TestGetOrderRequiresID(t *testing.T) {
	service := NewService(fakeRepository{})

	_, err := service.GetOrder(context.Background(), " ")
	if !errors.Is(err, ErrInvalidOrderID) {
		t.Fatalf("expected ErrInvalidOrderID, got %v", err)
	}
}
