package orderapplication

import (
	"context"
	"errors"
	orderdomain "project-example/internal/modules/order/domain"
	applogger "project-example/internal/platform/logger"
	"strings"

	"go.uber.org/zap"
)

var ErrInvalidOrderID = errors.New("order id is required")

type Service struct {
	repo orderdomain.Repository
}

func NewService(repo orderdomain.Repository) *Service {
	return &Service{repo: repo}
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

	return &OrderDTO{
		ID:           order.ID,
		CustomerName: order.CustomerName,
		Status:       order.Status,
		TotalAmount:  order.TotalAmount,
	}, nil
}
