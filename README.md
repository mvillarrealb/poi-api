# Points of Interest api with GO

A restful api for creating Points of interest using uber H3 :) and Go

# Exposed Resources

The following Resources & Operations are exposed in our Api:

Path|Verb|Description
---|---|---
/categories|GET|List all categories
/categories|POST|Create a category
/categories/{id}|GET|Get a category by id
/categories/{id}|PUT| Update a category
/categories/{id}|DELETE| Delete a Category
/geo_code|GET|Get the area of a given latitude, longitude
/poi|GET|List all Points of interest
/poi|POST|Create a point of interest
/poi/{id}|GET|Get a single point of interest
/poi/{id}|DELETE|Delete a point of interest

# Module And dependencies

The following dependencies were used for the api:

* github.com/go-playground/validator
* github.com/gorilla/mux
* gorm.io/driver/postgres
* github.com/gorilla/handlers
* gorm.io/gorm
* github.com/uber/h3-go

# Directory Structure

This is the directory structure to support our api

```sh
|
|-|Dockerfile
|-|go.mod #Go module file
|-|main.go #Application entry point
|-|api # Api related directory
|---|common #Common functions
|---|models # GORM models
|---|routes # Mux routes
|---|api.go # Main api
```
# Api entrypoint

Our api entrypoint will create a new Api reference passing as parameters the following environment variables

Environment|Description
---|---
DATABASE_PORT|Database instance port
DATABASE_HOST| Database hostname/ip address
DATABASE_USERNAME|Database user
DATABASE_PASSWORD|Database password

```go
//main.go
package main

import (
	"os"
	"poi-api/api"
	"strconv"
)

func main() {
	port, err := strconv.Atoi(os.Getenv("DATABASE_PORT"))
	if err != nil {
		panic("Invalid Database port")
	} else {
		app := api.New(
			os.Getenv("DATABASE_HOST"),
			port,
			os.Getenv("DATABASE_USERNAME"),
			os.Getenv("DATABASE_PASSWORD"))
		app.Run(":8080")
	}
}

```

# Run the go application locally

```sh
# It's required a a postgresql database instance

DATABASE_PORT=5432 DATABASE_HOST=127.0.0.1 DATABASE_USERNAME=postgres DATABASE_PASSWORD=casa1234 go run main.go 
```

# Load Some data

In order to use the geocode endpoint and create a referenced point of interest load the [data file](./data/PE-Lima.sql), this includes all divisions on the city of Lima. I've indexed every area of the city with h3 hexagons using 11 as resolution for the index.



# Build & Run Docker Image locally

```sh
docker build -t poi-api:0.0.1 .

docker run --name poi-api --network=host -e DATABASE_PASSWORD=casa1234 poi-api:0.0.1
```
---

# Run the Project on GCP Cloud Run

To run the project on Cloud Run [follow the tutorial](./articles/serverless_api.md)


# Use the postman collection

Use the [poi-e2e.postman_collection.json](./poi-e2e.postman_collection.json) collection to run some tests and play with the api