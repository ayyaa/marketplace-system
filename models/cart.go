package models

import (
	"time"
)

type Cart struct {
	CartID     int          `gorm:"primaryKey"`
	CartUUID   string       `gorm:"type:uuid;not null"`
	CustomerID int          `gorm:"not null"`
	CartStatus string       `gorm:"type:varchar(10);not null;check:cart_status IN ('active', 'converted', 'abandoned')"`
	CreatedAt  time.Time    `gorm:"autoCreateTime"`
	UpdatedAt  time.Time    `gorm:"autoUpdateTime"`
	Details    []CartDetail `gorm:"foreignKey:CartID"`
}

type CartRequest struct {
	ProductSlug string `json:"product_slug" validate:"required"`
	Quantity    int    `json:"quantity" validate:"required,min=1"`
}

type DeleteCartRequest struct {
	ProductSlug string `json:"product_slug" validate:"required"`
}

type ActionCart struct {
	ProductSlug string `json:"product_slug" validate:"required"`
	ProductID   int    `json:"product_id" validate:"required"`
	Quantity    int    `json:"quantity" validate:"required,min=1"`
	CustomerID  int    `json:"customer_id" validate:"required"`
	CartID      int    `json:"cart_id" validate:"required"`
}
