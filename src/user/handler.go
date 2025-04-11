package user

import (
	"encoding/json"
	"errors"
	"fmt"
	"gpsd-user-mgmt/src/auth"
	"io"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

func Get(c *gin.Context) {
	id := c.Param("id")
	userId, err := strconv.Atoi(id)
	if err != nil {
		c.JSON(400, gin.H{
			"error": err.Error(),
		})
		return
	}
	user, err := GetUser(userId)

	if err != nil {
		if errors.Is(err, NotFound{}) {
			c.JSON(http.StatusNotFound, gin.H{
				"error": err.Error(),
			})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, user)
}

func List(c *gin.Context) {
	var errorStrings []string
	limit, err := strconv.Atoi(c.DefaultQuery("limit", "15"))
	if err != nil {
		errString := fmt.Sprintf("Improper limit parameter: %v", limit)
		errorStrings = append(errorStrings, errString)
		slog.Error(errString)
	}

	offset, err := strconv.Atoi(c.DefaultQuery("offset", "0"))
	if err != nil {
		errString := fmt.Sprintf("Improper offset parameter: %v", offset)
		errorStrings = append(errorStrings, errString)
		slog.Error(errString)
	}

	if len(errorStrings) != 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"errors": errorStrings,
		})
		return
	}

	users, err := GetUsers(limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Internal Server Error",
			"e":     err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"users":  users,
		"limit":  limit,
		"offset": offset,
	})
}

func Create(c *gin.Context) {
	requestBody, err := io.ReadAll(c.Request.Body)
	if err != nil {
		c.JSON(400, gin.H{
			"error": err.Error(),
		})
		return
	}

	var user User
	err = json.Unmarshal(requestBody, &user)
	if err != nil {
		c.JSON(400, gin.H{
			"error": err.Error(),
		})
		return
	}

	err = user.AddUser()
	if err != nil {
		c.JSON(400, gin.H{
			"error": err.Error(),
		})
		return
	}
	c.JSON(200, gin.H{
		"message": "User created successfully",
		"user":    user,
	})
}

func Edit(c *gin.Context) {
	id := c.Param("id")
	userId, err := strconv.Atoi(id)
	if err != nil {
		c.JSON(400, gin.H{
			"error": err.Error(),
		})
		return
	}

	requestBody, err := io.ReadAll(c.Request.Body)
	if err != nil {
		c.JSON(400, gin.H{
			"error": err.Error(),
		})
		return
	}

	var user User
	err = json.Unmarshal(requestBody, &user)
	if err != nil {
		c.JSON(400, gin.H{
			"error": err.Error(),
		})
		return
	}

	err = UpdateUser(userId, user)
	if err != nil {
		if errors.Is(err, NotFound{}) {
			c.JSON(http.StatusNotFound, gin.H{
				"error": err.Error(),
			})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	user.UserId = userId
	c.JSON(200, gin.H{
		"message": "User updated successfully",
		"user":    user,
	})
}

func Delete(c *gin.Context) {
	id := c.Param("id")
	userId, err := strconv.Atoi(id)
	if err != nil {
		c.JSON(404, gin.H{
			"error": err.Error(),
		})
		return
	}

	err = DeleteUser(userId)
	if err != nil {
		if errors.Is(err, NotFound{}) {
			c.JSON(http.StatusNotFound, gin.H{
				"error": err.Error(),
			})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(200, gin.H{
		"message": "User deleted successfully",
	})
}

func SignIn(c *gin.Context) {
	requestBody, err := io.ReadAll(c.Request.Body)
	if err != nil {
		c.JSON(400, gin.H{
			"error": err.Error(),
		})
		return
	}

	var payload signInRequest
	err = json.Unmarshal(requestBody, &payload)
	if err != nil {
		c.JSON(400, gin.H{
			"error": err.Error(),
		})
		return
	}

	user, err := GetUserFromName(payload.Username)
	if err != nil {
		if errors.Is(err, NotFound{}) {
			c.JSON(http.StatusNotFound, gin.H{
				"error": err.Error(),
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}
	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(payload.Password))
	if err != nil {
		slog.Error(err.Error())
		c.JSON(http.StatusForbidden, gin.H{
			"error": err.Error(),
		})
		return
	}

	token, err := auth.CreateToken(payload.Username)
	if err != nil {
		slog.Error(err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Signed in successfully",
		"user":    user,
		"token":   token,
	})
}

func SignOut(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "Signed out successfully",
	})
}

func Verify(c *gin.Context) {
	token := c.Request.Header.Get("Authorization")

	if token == "" {
		slog.Info("No token")
		c.JSON(http.StatusForbidden, gin.H{
			"message": Unauthorized{}.Error(),
		})
		return
	}

	err := auth.VerifyToken(token)
	if err != nil {
		c.JSON(http.StatusForbidden, gin.H{
			"message": Unauthorized{}.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Valid token",
	})
}
