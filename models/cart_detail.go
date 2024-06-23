package models

import "time"

type CartDetail struct {
	CartDetailID     int    `gorm:"primaryKey"`
	CartDetailUUID   string `gorm:"type:uuid;not null"`
	CartID           int    `gorm:"not null"`
	ProductID        int    `gorm:"not null"`
	Quantity         int    `gorm:"not null"`
	CartDetailStatus string
	CreatedAt        time.Time `gorm:"autoCreateTime"`
	UpdatedAt        time.Time `gorm:"autoUpdateTime"`
}
