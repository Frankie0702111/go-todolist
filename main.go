package main

import (
	"go-todolist/router"
)

// @title Gin swagger
// @version 1.0
// @description Gin swagger

// @host localhost:8642
// @BasePath /api/v1
// schemes http
func main() {
	router.SetupRouter()
}
