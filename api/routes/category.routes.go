package routes

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"poi-api/api/common"
	"poi-api/api/models"

	"github.com/go-playground/validator"
	"github.com/gorilla/mux"
	"gorm.io/gorm"
)

type CategoryRoutes struct {
	DB *gorm.DB
}

func (api *CategoryRoutes) Init(router *mux.Router) {
	router.HandleFunc("/categories", api.create).Methods("POST")
	router.HandleFunc("/categories", api.findAll).Methods("GET")
	router.HandleFunc("/categories/{id:[0-9]+}", api.findById).Methods("GET")
	router.HandleFunc("/categories/{id:[0-9]+}", api.delete).Methods("DELETE")
}

func (api *CategoryRoutes) create(w http.ResponseWriter, req *http.Request) {
	var category models.Category
	decoder := json.NewDecoder(req.Body)
	if err := decoder.Decode(&category); err != nil {
		common.RespondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	validate := validator.New()
	err := validate.Struct(category)

	if err != nil {
		validationErrors := err.(validator.ValidationErrors)
		common.RespondWithErrorDetails(w, http.StatusBadRequest, "Invalid input parameters", validationErrors.Error())
		return
	}

	if err := api.DB.Create(&category).Error; err != nil {
		log.Fatal(err)
		common.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	common.RespondWithJSON(w, 201, category)
}

func (api *CategoryRoutes) findAll(w http.ResponseWriter, req *http.Request) {
	var categories []models.Category
	api.DB.Scopes(models.Paginate(req)).Find(&categories)
	common.RespondWithJSON(w, 200, categories)
}

func (api *CategoryRoutes) findById(w http.ResponseWriter, req *http.Request) {
	var category models.Category
	vars := mux.Vars(req)
	categoryID := vars["id"]
	err := api.DB.First(&category, categoryID).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		common.RespondWithError(w, 404, "Requested category was not found")
	} else {
		common.RespondWithJSON(w, 200, category)
	}
}

func (api *CategoryRoutes) delete(w http.ResponseWriter, req *http.Request) {
	var category models.Category
	vars := mux.Vars(req)
	categoryID := vars["id"]
	err := api.DB.Delete(&category, categoryID).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		common.RespondWithError(w, 404, "Requested category was not found")
	} else {
		common.RespondWithJSON(w, 202, category)
	}
}
