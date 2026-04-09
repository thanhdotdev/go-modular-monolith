package orderapplication

import (
	"context"
	"errors"
	discountapplication "project-example/internal/modules/discount/application"
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

type fakeDiscounts struct {
	discount *discountapplication.DiscountDTO
	err      error
}

func (f fakeDiscounts) GetDiscount(_ context.Context, _ string) (*discountapplication.DiscountDTO, error) {
	if f.err != nil {
		return nil, f.err
	}

	return f.discount, nil
}

func TestGetOrder(t *testing.T) {
	service := NewService(fakeRepository{
		order: &orderdomain.Order{
			ID:           "ord-001",
			DiscountCode: "WELCOME10",
			CustomerName: "Alice",
			Status:       "pending",
			TotalAmount:  125000,
		},
	}, fakeDiscounts{
		discount: &discountapplication.DiscountDTO{
			Code:   "WELCOME10",
			Type:   "fixed",
			Value:  10000,
			Active: true,
		},
	})

	got, err := service.GetOrder(context.Background(), "ord-001")
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}

	if got.ID != "ord-001" {
		t.Fatalf("expected order id ord-001, got %s", got.ID)
	}

	if got.DiscountCode != "WELCOME10" {
		t.Fatalf("expected discount code WELCOME10, got %s", got.DiscountCode)
	}

	if got.Discount == nil || got.Discount.Code != "WELCOME10" {
		t.Fatalf("expected discount details to be loaded, got %#v", got.Discount)
	}
}

func TestGetOrderRequiresID(t *testing.T) {
	service := NewService(fakeRepository{}, nil)

	_, err := service.GetOrder(context.Background(), " ")
	if !errors.Is(err, ErrInvalidOrderID) {
		t.Fatalf("expected ErrInvalidOrderID, got %v", err)
	}
}
