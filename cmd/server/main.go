package main

import (
	"github.com/gin-gonic/gin"

	"github.com/CedricThomas/22h31-FaisLesBacks/internal/pkg/middleware"
)

func main() {
	r := gin.Default()
	r.GET("/", middleware.Auth0("./dev-dgoly5h6.pem", []string{"casseur_flutter"}, "https://dev-dgoly5h6.eu.auth0.com/"), func(c *gin.Context) {
		c.JSON(200, gin.H{"test": "ok"})
	})
	r.Run(":9090")
}
