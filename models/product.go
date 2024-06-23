package models

import (
	"time"
)

type Product struct {
	ProductID     int       `gorm:"primaryKey;autoIncrement"`
	ProductUUID   string    `gorm:"type:uuid;not null"`
	CategoryID    int       `gorm:"not null"`
	ProductName   string    `gorm:"type:varchar(255);not null"`
	ProductSlug   string    `gorm:"type:varchar(255);not null"`
	Price         float64   `gorm:"type:numeric(10,2);not null"`
	StockQuantity int       `gorm:"not null"`
	Description   string    `gorm:"type:text"`
	ProductStatus string    `gorm:"type:varchar(10);not null;check:product_status IN ('active', 'deleted')"`
	CreatedAt     time.Time `gorm:"autoCreateTime"`
	UpdatedAt     time.Time `gorm:"autoUpdateTime"`
}
