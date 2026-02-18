package main

import (
	"payment-service/internal/request"
)


func main() {

	router := request.InitRouter()

	router.Run(":8080")
}
