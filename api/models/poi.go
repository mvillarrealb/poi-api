package models

import (
	"time"
)

type Poi struct {
	ID         uint      `gorm:"primary_key" json:"id"`
	Name       string    `json:"name" validate:"required"`
	Area1Name  string    `json:"area1Name" validate:"-"`
	Area2Name  string    `json:"area2Name" validate:"-"`
	Area3Name  string    `json:"area3Name" validate:"-"`
	Latitude   float64   `json:"latitude" validate:"required,latitude"`
	Longitude  float64   `json:"longitude" validate:"required,longitude"`
	CategoryID int       `json:"categoryId" validate:"required,number"`
	Category   Category  `json:"category" validate:"-" gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	CreatedAt  time.Time `json:"createdAt"`
	UpdatedAt  time.Time `json:"updatedAt"`
}
