package main

import (
	"net/http"
	"os"

	"github.com/Hassanhs-97/Data-Processing/handlers"
	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	r.POST("/process-data", handlers.CreateData)

	port := os.Getenv("PORT")
	http.ListenAndServe(":"+port, r)
}
