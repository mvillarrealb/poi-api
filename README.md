# Create & Deploy a Rest Api With Go + Cloud Run + Github Actions

# Creating GO Rest Api

We will be creating a POI(Point of interest) management api, the following operations will be supported in our api:

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

## Module And dependencies

To start our project we will create a module as follows:

```sh
mkdir poi-api && cd poi-api
go module init poi-api
```

The following dependencies must be added:

* github.com/go-playground/validator
* github.com/gorilla/mux
* gorm.io/driver/postgres
* github.com/gorilla/handlers
* gorm.io/gorm

## Directory Structure

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
## Api entrypoint

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
## Run the go application locally

```sh
# It's required a a postgresql database instance

DATABASE_PORT=5432 DATABASE_HOST=127.0.0.1 DATABASE_USERNAME=postgres DATABASE_PASSWORD=casa1234 go run main.go 
```

# Dockerfile

To create a containerized version of our api, we will use docker multi stage build and take advantage of Go static compilation to create a lightweight image.

```dockerfile
#Build step
FROM golang:1.15 as builder
RUN mkdir -p /poi-api/api
WORKDIR /poi-api
ADD api ./api
COPY go.mod go.sum main.go ./
#static compilation options for go
RUN go build -ldflags "-linkmode external -extldflags -static" -o main .


#Run step
#Scratch image is an empty image to add our binary, so the image will be as small as possible
FROM scratch
#Environments for dataase connection
ENV DATABASE_HOST="127.0.0.1" \
DATABASE_PORT="5432" \
DATABASE_USERNAME="postgres" \
DATABASE_PASSWORD="password"
#Copy binary from builder
COPY --from=builder /poi-api/main ./main
CMD ["./main"]
```

## Build && Run Docker Image locally

```sh
docker build -t poi-api:0.0.1 .
docker run --name poi-api --network=host -e DATABASE_PASSWORD=casa1234 poi-api:0.0.1
```
---

# Google Cloud project setup

We need to following steps to setup a GCP project with cloud Run enabled and cloud SQL database linked

```sh

# Create GCP project
gcloud projects create mvillarreal-demo-platform

# Set project as current running project
gcloud config set project mvillarreal-demo-platform

# Enable cloud Run api
gcloud services enable run.googleapis.com

# Enable container Registry
gcloud services enable containerregistry.googleapis.com

# Enable Cloud SQL services
gcloud services enable sqladmin.googleapis.com

# Enable networking services
gcloud services enable servicenetworking.googleapis.com

# Service account for Github actions
gcloud iam service-accounts create mvillarrealb-gha-saccount \
--description "Main service account for github actions" \
--display-name "mvillarreal-gha-saccount"


# Assign editor role for service account(for terraform)
gcloud projects add-iam-policy-binding mvillarreal-demo-platform \
--member serviceAccount:mvillarrealb-gha-saccount@mvillarreal-demo-platform.iam.gserviceaccount.com \
--role roles/editor 
  

# Export service account key for terraform(keep this in a safe place)
gcloud iam service-accounts keys create $(pwd)/service-account-key.json \
--iam-account mvillarrealb-gha-saccount@mvillarreal-demo-platform.iam.gserviceaccount.com

# Apply terraform
cd terraform && terraform apply

gcloud sql connect mvillarrealb-pg-sql --user=postgres

CREATE DATABASE poi_manager
```

# Github Actions Pipeline

```yaml
name: poi-api
on:
  push:
    branches:
    - master
env:
  REGION: us-east1-a
  PROJECT_ID: mvillarrealb-demo-platform
  BASE_IMAGE: gcr.io/mvillarrealb-demo-platform/poi-api
  DATABASE_INSTANCE: mvillarrealb-pg-sql
  SERVICE_NAME: poi-api
jobs:
  build:
    runs-on: ubuntu-latest
    steps:
    - name: Setup Project
      id: checkout
      uses: actions/checkout@master
    - name: Build & Publish Image
      uses: docker/build-push-action@v2
      id: build
      with:
        context: .
        push: true
        tags: ${{ github.sha }}
        image: $BASE_IMAGE
    - name: Deploy to Cloud Run
      id: deploy
      uses: google-github-actions/deploy-cloudrun@main
      with:
        region: $REGION
        service: $SERVICE_NAME
        image: $BASE_IMAGE/${{ github.sha }}
        credentials: ${{ secrets.GCP_SA_KEY }}
        env_vars: "DATABASE_USERNAME=${{ secrets.DATABASE_USERNAME }},DATABASE_PASSWORD=${{ secrets.DATABASE_PASSWORD }}"
        flags: "--add-cloudsql-instances '$PROJECT_ID:$REGION:$DATABASE_INSTANCE'"
    - name: Run e2e Test
      uses: matt-ball/newman-action@master
      id: test
      with:
        collection: postman_collection.json
        options:
          envVar:
            BASE_URL: "${{ steps.deploy.outputs.url }}"
```
