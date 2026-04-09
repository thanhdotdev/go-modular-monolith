package orderapplication

type DiscountDTO struct {
	Code   string `json:"code"`
	Type   string `json:"type"`
	Value  int64  `json:"value"`
	Active bool   `json:"active"`
}

type OrderDTO struct {
	ID           string       `json:"id"`
	DiscountCode string       `json:"discountCode,omitempty"`
	CustomerName string       `json:"customerName"`
	Status       string       `json:"status"`
	TotalAmount  int64        `json:"totalAmount"`
	Discount     *DiscountDTO `json:"discount,omitempty"`
}
