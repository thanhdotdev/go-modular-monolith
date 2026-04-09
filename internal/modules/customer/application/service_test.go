package customerapplication

import (
	"context"
	"errors"
	customerdomain "project-example/internal/modules/customer/domain"
	"testing"
)

type fakeRepository struct {
	customer *customerdomain.Customer
	err      error
}

func (f fakeRepository) FindByID(_ context.Context, _ string) (*customerdomain.Customer, error) {
	if f.err != nil {
		return nil, f.err
	}

	return f.customer, nil
}

func TestGetCustomer(t *testing.T) {
	service := NewService(fakeRepository{
		customer: &customerdomain.Customer{
			ID:    "cus-001",
			Name:  "Alice",
			Email: "alice@example.com",
			Tier:  "gold",
		},
	})

	got, err := service.GetCustomer(context.Background(), "cus-001")
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}

	if got.ID != "cus-001" {
		t.Fatalf("expected customer id cus-001, got %s", got.ID)
	}
}

func TestGetCustomerRequiresID(t *testing.T) {
	service := NewService(fakeRepository{})

	_, err := service.GetCustomer(context.Background(), " ")
	if !errors.Is(err, ErrInvalidCustomerID) {
		t.Fatalf("expected ErrInvalidCustomerID, got %v", err)
	}
}
