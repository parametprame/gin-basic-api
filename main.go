package main

import (
	"example/go-rest-api/api"
	"example/go-rest-api/db"

	"github.com/gin-gonic/gin"
)

func main() {

	db := db.InitDB()

	handler := api.NewHandler(db)

	r := gin.New()

	api.InitRoutes(r, handler)

	r.Run()
}
