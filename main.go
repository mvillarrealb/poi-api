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
