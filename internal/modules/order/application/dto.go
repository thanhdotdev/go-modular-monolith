package orderapplication

type OrderDTO struct {
	ID           string `json:"id"`
	CustomerName string `json:"customerName"`
	Status       string `json:"status"`
	TotalAmount  int64  `json:"totalAmount"`
}
