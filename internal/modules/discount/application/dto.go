package discountapplication

type DiscountDTO struct {
	Code   string `json:"code"`
	Type   string `json:"type"`
	Value  int64  `json:"value"`
	Active bool   `json:"active"`
}
