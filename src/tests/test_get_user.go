package tests

import (
	"encoding/json"
	"fmt"
	"gpsd-user-mgmt/db"
	"gpsd-user-mgmt/router"
	"gpsd-user-mgmt/user"
	"net/http"
	"net/http/httptest"
	"testing"

	"gotest.tools/assert"
)

func successGet(r *router.Engine) func(*testing.T) {
	return func(t *testing.T) {
		testUsers := []user.User{{
			Name:  "Test",
			DevID: "123",
			Role:  "reporter",
		}, {
			Name:  "Test2",
			DevID: "1234",
			Role:  "reporter",
		},
		}

		for i, _ := range testUsers {
			id, _ := user.AddUser(testUsers[i])
			testUsers[i].Id = id
		}
		defer db.EmptyDatabase()

		for _, testUser := range testUsers {
			w := httptest.NewRecorder()
			url := fmt.Sprintf("%s/%d", USER_API, testUser.Id)
			req, _ := http.NewRequest(
				"GET",
				url,
				nil,
			)
			r.ServeHTTP(w, req)
			assert.Equal(t, 200, w.Code)

			var body map[string]interface{}
			err := json.Unmarshal(w.Body.Bytes(), &body)
			assert.NilError(t, err)

			assert.Equal(t, body["name"], testUser.Name)
			assert.Equal(t, body["role"], testUser.Role)
			assert.Equal(t, body["devID"], testUser.DevID)
		}
	}
}

func notFoundGet(r *router.Engine) func(*testing.T) {
	return func(t *testing.T) {
		w := httptest.NewRecorder()
		url := fmt.Sprintf("%s/%d", USER_API, 0)
		req, _ := http.NewRequest(
			"GET",
			url,
			nil,
		)
		r.ServeHTTP(w, req)

		assert.Equal(t, 404, w.Code)

		var body map[string]interface{}

		err := json.Unmarshal(w.Body.Bytes(), &body)
		assert.NilError(t, err)

		assert.Equal(t, body["error"].(string), "User not found")
	}
}
