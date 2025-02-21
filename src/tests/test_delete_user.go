package tests

import (
	"encoding/json"
	"fmt"
	"gpsd-user-mgmt/src/db"
	"gpsd-user-mgmt/src/router"
	"gpsd-user-mgmt/src/user"
	"net/http"
	"net/http/httptest"
	"testing"

	"gotest.tools/assert"
)

func successDelete(r *router.Engine) func(*testing.T) {
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
				"DELETE",
				url,
				nil,
			)
			r.ServeHTTP(w, req)
			assert.Equal(t, 200, w.Code)

			var body map[string]interface{}
			err := json.Unmarshal(w.Body.Bytes(), &body)
			assert.NilError(t, err)
			assert.Equal(t, body["message"], "User deleted successfully")
		}
	}
}

func notFoundDelete(r *router.Engine) func(*testing.T) {
	return func(t *testing.T) {

		w := httptest.NewRecorder()
		url := fmt.Sprintf("%s/%d", USER_API, 0)
		req, _ := http.NewRequest(
			"DELETE",
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
