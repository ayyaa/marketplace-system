package models

import (
	"time"
)

type Cart struct {
	CartID     uint         `gorm:"primaryKey"`
	CartUUID   string       `gorm:"type:uuid;not null"`
	CustomerID uint         `gorm:"not null"`
	CartStatus string       `gorm:"type:varchar(10);not null;check:cart_status IN ('active', 'converted', 'abandoned')"`
	CreatedAt  time.Time    `gorm:"autoCreateTime"`
	UpdatedAt  time.Time    `gorm:"autoUpdateTime"`
	Details    []CartDetail `gorm:"foreignKey:CartID"`
}
