package models

import "time"

type CartDetail struct {
	CartDetailID uint      `gorm:"primaryKey"`
	CartID       uint      `gorm:"not null"`
	ProductID    uint      `gorm:"not null"`
	Quantity     int       `gorm:"not null"`
	CreatedAt    time.Time `gorm:"autoCreateTime"`
	UpdatedAt    time.Time `gorm:"autoUpdateTime"`
}
