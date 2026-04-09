package orderdomain

type Order struct {
	ID           string
	CustomerName string
	Status       string
	TotalAmount  int64
}
