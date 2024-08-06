package main

import (
	"github.com/gin-gonic/gin"
)

func main() {
    r := gin.Default() // Create a default Gin router

    // Define a route and handler
    r.GET("/ping", func(c *gin.Context) {
        c.JSON(200, gin.H{
            "message": "pong",
        })
    })

    // Run the server on port 8080
    r.Run() // default listens on :8080
}