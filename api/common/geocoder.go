package common

import (
	"poi-api/api/models"

	"github.com/uber/h3-go"
	"gorm.io/gorm"
)

type Geocoder struct {
	DB *gorm.DB
}

func (geocoder *Geocoder) GetAddress(latitude float64, longitude float64) models.Area {
	var area models.Area
	h3Index := geocoder.createH3Index(latitude, longitude)
	geocoder.DB.Where("h3_index", h3.ToString(h3Index)).First(&area)
	return area
}

func (geocoder *Geocoder) createH3Index(latitude float64, longitude float64) h3.H3Index {
	geo := h3.GeoCoord{
		Latitude:  float64(latitude),
		Longitude: float64(longitude),
	}
	resolution := 10
	return h3.FromGeo(geo, resolution)
}
