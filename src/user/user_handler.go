package user

import (
	"encoding/json"
	"io"
	"log"

	"github.com/gin-gonic/gin"
)

func List(c *gin.Context) {
	response := gin.H{}
	for _, v := range users {
		response[v.Name] = v.DevID
	}

	c.JSON(200, response)
}

func Create(c *gin.Context) {
	response := gin.H{}

	requestBody, err := io.ReadAll(c.Request.Body)
	if err != nil {
		log.Printf("Create Error : %v\n", err.Error())
		c.JSON(400, gin.H{
			"Error": err.Error(),
		})
	}

	var jsonBody User
	err = json.Unmarshal(requestBody, &jsonBody)
	if err != nil {
		log.Printf("Create Error : %v\n", err.Error())
		c.JSON(400, gin.H{
			"Error": err.Error(),
		})
	}

	users = append(users, jsonBody)
	response[jsonBody.Name] = jsonBody.DevID
	c.JSON(200, response)
}

func Edit(c *gin.Context) {
	c.JSON(200, gin.H{
		"API": "edit",
	})
}
