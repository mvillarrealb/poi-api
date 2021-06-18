package api

import (
	"fmt"
	"log"
	"net/http"
	"poi-api/api/common"
	"poi-api/api/models"
	"poi-api/api/routes"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Api struct {
	router *mux.Router
}

func New(host string, port int, username string, password string) *Api {
	api := &Api{router: mux.NewRouter().StrictSlash(true)}
	db := api.databaseInit(host, port, username, password)
	poiRoutes := &routes.POIRoutes{DB: db, Geocoder: &common.Geocoder{DB: db}}
	categoryRoutes := &routes.CategoryRoutes{DB: db}
	poiRoutes.Init(api.router)
	categoryRoutes.Init(api.router)

	return api
}

func (a *Api) databaseInit(host string, port int, username string, password string) *gorm.DB {
	dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=poi_manager sslmode=disable TimeZone=America/Lima",
		host, port, username, password)
	database, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}
	database.AutoMigrate(&models.Poi{}, &models.Category{}, &models.Area{})
	return database
}

func (app *Api) Run(addr string) {
	handledRouter := handlers.CORS(
		handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type", "Authorization"}),
		handlers.AllowedMethods([]string{"GET", "POST", "PUT", "HEAD", "OPTIONS"}),
		handlers.AllowedOrigins([]string{"*"}),
	)(app.router)
	if err := http.ListenAndServe(addr, handledRouter); err != nil {
		log.Fatal(err)
	} else {
		log.Println("Successfully started http server")
	}
}
