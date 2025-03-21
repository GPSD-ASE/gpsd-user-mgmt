package tests

import (
	"encoding/json"
	"fmt"
	"gpsd-user-mgmt/src/db"
	"gpsd-user-mgmt/src/router"
	"gpsd-user-mgmt/src/user"
	"math/rand"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"

	"gotest.tools/assert"
)

func successUpdate(r *router.Engine) func(*testing.T) {
	return func(t *testing.T) {
		testUsers := []user.User{{
			UserName: "Test",
			DeviceID: "123",
			Role:     "reporter",
		}, {
			UserName: "Test2",
			DeviceID: "1234",
			Role:     "reporter",
		},
		}

		for i, _ := range testUsers {
			id, _ := user.AddUser(testUsers[i])
			testUsers[i].UserId = id
		}
		defer db.EmptyDatabase()

		for i, _ := range testUsers {
			testUsers[i].UserName = randomString(30)
			testUsers[i].DeviceID = strconv.Itoa(rand.Int())
		}

		for _, testUser := range testUsers {
			w := httptest.NewRecorder()
			payload, _ := json.Marshal(testUser)
			url := fmt.Sprintf("%s/%d", USER_API, testUser.UserId)
			req, _ := http.NewRequest(
				"PUT",
				url,
				strings.NewReader(string(payload)),
			)
			r.ServeHTTP(w, req)
			assert.Equal(t, 200, w.Code)

			var body map[string]interface{}
			err := json.Unmarshal(w.Body.Bytes(), &body)
			assert.NilError(t, err)
			assert.Equal(t, body["message"], "User updated successfully")

			userBody := body["user"].(map[string]interface{})

			assert.Equal(t, userBody["name"], testUser.UserName)
			assert.Equal(t, userBody["role"], testUser.Role)
			assert.Equal(t, userBody["devID"], testUser.DeviceID)
		}
	}
}

func notFoundUpdate(r *router.Engine) func(*testing.T) {
	return func(t *testing.T) {
		testUser := user.User{
			UserName: "Test",
			DeviceID: "123",
			Role:     "reporter",
		}

		w := httptest.NewRecorder()
		url := fmt.Sprintf("%s/%d", USER_API, 0)
		payload, _ := json.Marshal(testUser)
		req, _ := http.NewRequest(
			"PUT",
			url,
			strings.NewReader(string(payload)),
		)
		r.ServeHTTP(w, req)
		assert.Equal(t, 404, w.Code)

		var body map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &body)
		assert.NilError(t, err)

		assert.Equal(t, body["error"].(string), "User not found")
	}
}
