package main

import (
	"github.com/gin-gonic/gin"

	"payment-service/internal/http"
)

func main() {

	router := gin.Default()

	router.POST("/v1/payments", http.Posting)

	router.GET("/v1/payments/:id", http.Getting)

	router.Run(":8080")
}
