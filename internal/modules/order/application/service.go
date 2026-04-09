package orderapplication

import (
	"context"
	"errors"
	discountapplication "project-example/internal/modules/discount/application"
	orderdomain "project-example/internal/modules/order/domain"
	applogger "project-example/internal/platform/logger"
	"strings"

	"go.uber.org/zap"
)

var ErrInvalidOrderID = errors.New("order id is required")

type discountGetter interface {
	GetDiscount(ctx context.Context, code string) (*discountapplication.DiscountDTO, error)
}

type Service struct {
	repo      orderdomain.Repository
	discounts discountGetter
}

func NewService(repo orderdomain.Repository, discounts discountGetter) *Service {
	return &Service{
		repo:      repo,
		discounts: discounts,
	}
}

func (s *Service) GetOrder(ctx context.Context, id string) (*OrderDTO, error) {
	id = strings.TrimSpace(id)
	applogger.FromContext(ctx).Debug("getting order", zap.String("order_id", id))
	if id == "" {
		return nil, ErrInvalidOrderID
	}

	order, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	dto := &OrderDTO{
		ID:           order.ID,
		DiscountCode: order.DiscountCode,
		CustomerName: order.CustomerName,
		Status:       order.Status,
		TotalAmount:  order.TotalAmount,
	}

	if s.discounts != nil && order.DiscountCode != "" {
		discount, err := s.discounts.GetDiscount(ctx, order.DiscountCode)
		if err != nil {
			applogger.FromContext(ctx).Warn(
				"failed to load discount",
				zap.String("order_id", order.ID),
				zap.String("discount_code", order.DiscountCode),
				zap.Error(err),
			)
		} else {
			dto.Discount = &DiscountDTO{
				Code:   discount.Code,
				Type:   discount.Type,
				Value:  discount.Value,
				Active: discount.Active,
			}
		}
	}

	return dto, nil
}
