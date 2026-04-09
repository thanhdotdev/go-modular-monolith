package ordermemory

import (
	"context"
	orderdomain "project-example/internal/modules/order/domain"
	"project-example/internal/shared/collection"
	"project-example/internal/shared/ptr"
)

type Repository struct {
	orders map[string]orderdomain.Order
}

func NewRepository(seed []orderdomain.Order) *Repository {
	return &Repository{
		orders: collection.IndexBy(seed, func(order orderdomain.Order) string {
			return order.ID
		}),
	}
}

func SeedOrders() []orderdomain.Order {
	return []orderdomain.Order{
		{
			ID:           "ord-001",
			DiscountCode: "WELCOME10",
			CustomerName: "Alice",
			Status:       "pending",
			TotalAmount:  125000,
		},
		{
			ID:           "ord-002",
			CustomerName: "Bob",
			Status:       "paid",
			TotalAmount:  340000,
		},
	}
}

func (r *Repository) FindByID(_ context.Context, id string) (*orderdomain.Order, error) {
	order, ok := r.orders[id]
	if !ok {
		return nil, orderdomain.ErrOrderNotFound
	}

	return ptr.Of(order), nil
}
