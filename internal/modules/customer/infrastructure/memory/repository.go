package customermemory

import (
	"context"
	customerdomain "project-example/internal/modules/customer/domain"
	"project-example/internal/shared/collection"
	"project-example/internal/shared/ptr"
)

type Repository struct {
	customers map[string]customerdomain.Customer
}

func NewRepository(seed []customerdomain.Customer) *Repository {
	return &Repository{
		customers: collection.IndexBy(seed, func(customer customerdomain.Customer) string {
			return customer.ID
		}),
	}
}

func SeedCustomers() []customerdomain.Customer {
	return []customerdomain.Customer{
		{
			ID:    "cus-001",
			Name:  "Alice",
			Email: "alice@example.com",
			Tier:  "gold",
		},
		{
			ID:    "cus-002",
			Name:  "Bob",
			Email: "bob@example.com",
			Tier:  "silver",
		},
	}
}

func (r *Repository) FindByID(_ context.Context, id string) (*customerdomain.Customer, error) {
	customer, ok := r.customers[id]
	if !ok {
		return nil, customerdomain.ErrCustomerNotFound
	}

	return ptr.Of(customer), nil
}
