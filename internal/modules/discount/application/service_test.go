package discountapplication

import (
	"context"
	"errors"
	"testing"

	discountdomain "project-example/internal/modules/discount/domain"
)

type fakeRepository struct {
	discount *discountdomain.Discount
	err      error
}

func (f fakeRepository) FindByCode(_ context.Context, _ string) (*discountdomain.Discount, error) {
	if f.err != nil {
		return nil, f.err
	}

	return f.discount, nil
}

func TestGetDiscount(t *testing.T) {
	service := NewService(fakeRepository{
		discount: &discountdomain.Discount{
			Code:   "WELCOME10",
			Type:   "fixed",
			Value:  10000,
			Active: true,
		},
	})

	got, err := service.GetDiscount(context.Background(), "WELCOME10")
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}

	if got.Code != "WELCOME10" {
		t.Fatalf("expected discount code WELCOME10, got %s", got.Code)
	}
}

func TestGetDiscountRequiresCode(t *testing.T) {
	service := NewService(fakeRepository{})

	_, err := service.GetDiscount(context.Background(), " ")
	if !errors.Is(err, ErrInvalidDiscountCode) {
		t.Fatalf("expected ErrInvalidDiscountCode, got %v", err)
	}
}
