package main

import (
	"net/http"
	"os"

	"github.com/Hassanhs-97/Data-Processing/handlers"
	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	dataProcessingGroup := r.Group("/process-data")
	dataProcessingGroup.POST("/", handlers.CreateData)

	port := os.Getenv("PORT")

	if port == "" {
		port = "8080"
	}
	http.ListenAndServe(":"+port, r)
}
