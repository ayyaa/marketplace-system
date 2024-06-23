package models

import "time"

type CustomerRequest struct {
	FullName string `json:"full_name" validate:"required,min=2,max=255"`
	Email    string `json:"email"  validate:"required,email,max=255"`
	Phone    string `json:"phone" validate:"required"`
	Password string `json:"password" validate:"required,min=8,max=255"`
	Address  string `json:"address"  validate:"omitempty"`
}

// Customer struct definition
type Customer struct {
	CustomerID     int       `gorm:"primaryKey;autoIncrement"`
	CustomerUUID   string    `gorm:"type:uuid;not null"`
	FullName       string    `gorm:"type:varchar(255);not null"`
	Email          string    `gorm:"type:varchar(255);unique;not null"`
	Phone          string    `gorm:"type:varchar(20);unique;not null"`
	Password       string    `gorm:"type:varchar(255);not null"`
	Address        string    `gorm:"type:text"`
	CustomerStatus string    `gorm:"type:varchar(10);not null;check:customer_status IN ('active', 'deleted')"`
	CreatedAt      time.Time `gorm:"autoCreateTime"`
	UpdatedAt      time.Time `gorm:"autoUpdateTime"`
}

type LoginRequest struct {
	// LoginReq defines model for LoginReq.
	Password string `json:"password"`
	Phone    string `json:"phone"`
}

type LoginResponse struct {
	Token string `json:"token"`
}
