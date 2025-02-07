package main

import (
	"fmt"
	"gpsd-user-mgmt/src/user"
	"os"

	"github.com/gin-gonic/gin"
)

func setupRouter() *gin.Engine {
	router := gin.Default()

	v1 := router.Group("/api/v1")
	{
		v1.GET("/list", user.List)
		v1.POST("/create", user.Create)

		v1.GET("/edit", func(c *gin.Context) {
			c.JSON(200, gin.H{
				"API": "edit",
			})
		})

		v1.GET("/signin", func(c *gin.Context) {
			c.JSON(200, gin.H{
				"API": "signin",
			})
		})

		v1.GET("/signout", func(c *gin.Context) {
			c.JSON(200, gin.H{
				"API": "signout",
			})
		})
	}

	return router
}

func main() {
	router := setupRouter()

	// port := os.Getenv("PORT")
	// if port == "" {
	// 	fmt.Println("Unable to read port")
	// 	return
	// }
	// address := fmt.Sprintf(":%s", port)

	router.Run(":5500")
}
