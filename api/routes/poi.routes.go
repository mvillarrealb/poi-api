package routes

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"poi-api/api/common"
	"poi-api/api/models"
	"strconv"

	"github.com/go-playground/validator"
	"github.com/gorilla/mux"
	"gorm.io/gorm"
)

type POIRoutes struct {
	DB       *gorm.DB
	Geocoder *common.Geocoder
}

func (api *POIRoutes) Init(router *mux.Router) {
	router.HandleFunc("/geo_code", api.reverseGeoCode).Methods("GET")
	router.HandleFunc("/poi", api.create).Methods("POST")
	router.HandleFunc("/poi", api.findAll).Methods("GET")
	router.HandleFunc("/poi/{id:[0-9]+}", api.findById).Methods("GET")
	router.HandleFunc("/poi/{id:[0-9]+}", api.delete).Methods("DELETE")
}

func (api *POIRoutes) reverseGeoCode(w http.ResponseWriter, req *http.Request) {
	query := req.URL.Query()
	latitude, errlat := strconv.ParseFloat(query.Get("latitude"), 64)
	longitude, errlon := strconv.ParseFloat(query.Get("longitude"), 64)

	if errlat != nil || errlon != nil {
		common.RespondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	area := api.Geocoder.GetAddress(latitude, longitude)

	common.RespondWithJSON(w, 200, area)
}

func (api *POIRoutes) create(w http.ResponseWriter, req *http.Request) {
	var poi models.Poi
	var category models.Category
	decoder := json.NewDecoder(req.Body)
	if err := decoder.Decode(&poi); err != nil {
		common.RespondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	validate := validator.New()
	err := validate.Struct(poi)

	if err != nil {
		validationErrors := err.(validator.ValidationErrors)
		common.RespondWithErrorDetails(w, http.StatusBadRequest, "Invalid input parameters", validationErrors.Error())
		return
	}

	categoryErr := api.DB.First(&category, poi.CategoryID).Error

	if categoryErr != nil {
		common.RespondWithError(w, 404, "Requested Category was not found")
		return
	}

	poi.Category = category

	area := api.Geocoder.GetAddress(poi.Latitude, poi.Longitude)

	poi.Area1Name = area.Area1Name
	poi.Area2Name = area.Area2Name
	poi.Area3Name = area.Area3Name

	if err := api.DB.Create(&poi).Error; err != nil {
		log.Fatal(err)
		common.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	common.RespondWithJSON(w, 201, poi)
}

func (api *POIRoutes) findAll(w http.ResponseWriter, req *http.Request) {
	var pois []models.Poi
	api.DB.Scopes(models.Paginate(req)).Joins("Category").Find(&pois)
	common.RespondWithJSON(w, 200, pois)
}

func (api *POIRoutes) findById(w http.ResponseWriter, req *http.Request) {
	var poi models.Poi
	vars := mux.Vars(req)
	poiID := vars["id"]
	err := api.DB.Joins("Category").First(&poi, poiID).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		common.RespondWithError(w, 404, "Requested Point of Interest was not found")
	} else {
		common.RespondWithJSON(w, 200, poi)
	}
}

func (api *POIRoutes) delete(w http.ResponseWriter, req *http.Request) {
	var poi models.Poi
	vars := mux.Vars(req)
	poiID := vars["id"]
	err := api.DB.Delete(&poi, poiID).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		common.RespondWithError(w, 404, "Requested Point of Interest was not found")
	} else {
		common.RespondWithJSON(w, 202, poi)
	}
}
