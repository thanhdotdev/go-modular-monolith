package orderpostgres

import orderdomain "project-example/internal/modules/order/domain"

func toDomain(model OrderModel) *orderdomain.Order {
	return &orderdomain.Order{
		ID:           model.ID,
		CustomerName: model.CustomerName,
		Status:       model.Status,
		TotalAmount:  model.TotalAmount,
	}
}
