package models

import "time"

type Category struct {
	ID           uint      `gorm:"primary_key" json:"id"`
	CategoryName string    `json:"categoryName" validate:"required"`
	CreatedAt    time.Time `json:"createdAt"`
	UpdatedAt    time.Time `json:"updatedAt"`
}
