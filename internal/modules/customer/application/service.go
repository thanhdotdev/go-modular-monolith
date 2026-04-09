package customerapplication

import (
	"context"
	"errors"
	customerdomain "project-example/internal/modules/customer/domain"
	"strings"
)

var ErrInvalidCustomerID = errors.New("customer id is required")

type Service struct {
	repo customerdomain.Repository
}

func NewService(repo customerdomain.Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) GetCustomer(ctx context.Context, id string) (*CustomerDTO, error) {
	id = strings.TrimSpace(id)
	if id == "" {
		return nil, ErrInvalidCustomerID
	}

	customer, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	return &CustomerDTO{
		ID:    customer.ID,
		Name:  customer.Name,
		Email: customer.Email,
		Tier:  customer.Tier,
	}, nil
}
