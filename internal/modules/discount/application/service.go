package discountapplication

import (
	"context"
	"errors"
	"strings"

	discountdomain "project-example/internal/modules/discount/domain"
)

var ErrInvalidDiscountCode = errors.New("discount code is required")

type Service struct {
	repo discountdomain.Repository
}

func NewService(repo discountdomain.Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) GetDiscount(ctx context.Context, code string) (*DiscountDTO, error) {
	code = strings.TrimSpace(code)
	if code == "" {
		return nil, ErrInvalidDiscountCode
	}

	discount, err := s.repo.FindByCode(ctx, code)
	if err != nil {
		return nil, err
	}

	return &DiscountDTO{
		Code:   discount.Code,
		Type:   discount.Type,
		Value:  discount.Value,
		Active: discount.Active,
	}, nil
}
