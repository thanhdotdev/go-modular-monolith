package orderdomain

type Order struct {
	ID           string
	DiscountCode string
	CustomerName string
	Status       string
	TotalAmount  int64
}
