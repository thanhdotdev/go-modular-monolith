package orderpostgres

type OrderModel struct {
	ID           string `gorm:"primaryKey;type:varchar(64)"`
	DiscountCode string `gorm:"column:discount_code;type:varchar(64);not null"`
	CustomerName string `gorm:"column:customer_name;type:varchar(255);not null"`
	Status       string `gorm:"type:varchar(50);not null"`
	TotalAmount  int64  `gorm:"column:total_amount;not null"`
}

func (OrderModel) TableName() string {
	return "orders"
}
