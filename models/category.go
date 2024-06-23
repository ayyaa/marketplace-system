package models

import "time"

type Category struct {
	CategoryID     uint      `gorm:"primaryKey;autoIncrement"`
	CategoryUUID   string    `gorm:"type:uuid;not null"`
	CategoryName   string    `gorm:"type:varchar(255);not null"`
	CategorySlug   string    `gorm:"type:varchar(255);not null"`
	CategoryStatus string    `gorm:"type:varchar(10);not null;check:category_status IN ('active', 'deleted')"`
	CreatedAt      time.Time `gorm:"autoCreateTime"`
	UpdatedAt      time.Time `gorm:"autoUpdateTime"`

	// Adding a Products slice to include all products in this category
	Products []Product `gorm:"foreignKey:CategoryID"`
}
