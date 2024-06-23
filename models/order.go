package models

import "time"

type Order struct {
	OrderID            int       `gorm:"primaryKey"`
	OrderUUID          string    `gorm:"type:uuid;not null"`
	InvoiceNumber      string    `gorm:"type:varchar(255);not null"`
	CustomerID         int       `gorm:"not null"`
	CartID             int       `gorm:"not null"`
	OrderDate          time.Time `gorm:"autoCreateTime"`
	OrderPaymentType   string    `gorm:"type:varchar(20);not null;check:order_payment_type IN ('alfamart', 'indomaret', 'virtual_account', 'jenius')"`
	OrderPaymentStatus string    `gorm:"type:varchar(10);not null;check:order_payment_status IN ('unpaid', 'completed')"`
	OrderStatus        string    `gorm:"type:varchar(10);not null;check:order_status IN ('pending', 'scheduled')"`
	GrandTotal         float64   `gorm:"type:numeric(10,2);not null"`
	ExpiredAt          time.Time `gorm:"not null"`
	CreatedAt          time.Time `gorm:"autoCreateTime"`
	UpdatedAt          time.Time `gorm:"autoUpdateTime"`

	Details []OrderDetail `gorm:"foreignKey:OrderDetailID"`
}

type OrderDetail struct {
	OrderDetailID     int     `gorm:"primaryKey"`
	OrderDetailUUID   string  `gorm:"type:uuid;not null"`
	OrderID           int     `gorm:"not null"`
	ProductID         int     `gorm:"not null"`
	Quantity          int     `gorm:"not null"`
	Price             float64 `gorm:"type:numeric(10,2);not null"`
	FinalPrice        float64 `gorm:"type:numeric(10,2);not null"`
	OrderDetailStatus string
	CreatedAt         time.Time `gorm:"autoCreateTime"`
	UpdatedAt         time.Time `gorm:"autoUpdateTime"`
}

type CheckoutRequest struct {
	PaymentType string `json:"payment_type" validate:"required,oneof=alfamart indomaret virtual_account jenius"`
}

type Checkout struct {
	PaymentType string `json:"payment_type" validate:"required,oneof=alfamart indomaret virtual_account jenius"`
	CustomerID  int    `json:"customer_id" validate:"required"`
}
