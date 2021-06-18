package models

type Area struct {
	ID        uint   `gorm:"primary_key" json:"id"`
	H3Index   string `json:"h3Index"`
	Area1Name string `json:"area1Name"`
	Area2Name string `json:"area2Name"`
	Area3Name string `json:"area3Name"`
}
