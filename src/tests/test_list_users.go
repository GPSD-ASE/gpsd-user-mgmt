package tests

import (
	"encoding/json"
	"fmt"
	"gpsd-user-mgmt/src/db"
	"gpsd-user-mgmt/src/router"
	"gpsd-user-mgmt/src/user"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"

	"gotest.tools/assert"
)

func successList(r *router.Engine) func(*testing.T) {
	return func(t *testing.T) {

		test_users := []user.User{{
			UserName: "Test",
			DeviceID: "123",
			Role:     "reporter",
		}, {
			UserName: "Test2",
			DeviceID: "1234",
			Role:     "reporter",
		},
		}

		for _, v := range test_users {
			user.AddUser(v)
		}
		defer db.EmptyDatabase()

		w := httptest.NewRecorder()
		req, _ := http.NewRequest(
			"GET",
			USER_API,
			nil,
		)
		r.ServeHTTP(w, req)

		assert.Equal(t, 200, w.Code)

		var body map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &body)
		assert.NilError(t, err)
		assert.Equal(t, int(body["limit"].(float64)), 15)
		assert.Equal(t, int(body["offset"].(float64)), 0)

		bodyUsers := body["users"].([]any)
		assert.Equal(t, len(bodyUsers), 2)

		for i := range len(bodyUsers) {
			bodyUser := bodyUsers[i].(map[string]interface{})
			assert.Equal(t, bodyUser["name"], test_users[i].UserName)
			assert.Equal(t, bodyUser["role"], test_users[i].Role)
			assert.Equal(t, bodyUser["devID"], test_users[i].DeviceID)
		}
	}
}

func successListQuery(r *router.Engine) func(*testing.T) {
	return func(t *testing.T) {

		testUsers := []user.User{{
			UserName: "Test",
			DeviceID: "123",
			Role:     "reporter",
		}, {
			UserName: "Test2",
			DeviceID: "1234",
			Role:     "reporter",
		}, {
			UserName: "Test3",
			DeviceID: "103",
			Role:     "reporter",
		}, {
			UserName: "Admin",
			DeviceID: "123004",
			Role:     "admin",
		},
		}

		for _, v := range testUsers {
			user.AddUser(v)
		}
		defer db.EmptyDatabase()

		query := map[string]interface{}{
			"limit":  3,
			"offset": 1,
		}
		var queryString string
		for k, v := range query {
			queryString += fmt.Sprintf("&%s=%v", k, v)
		}

		url := fmt.Sprintf("%s?%s", USER_API, queryString)
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", url, nil)
		r.ServeHTTP(w, req)

		assert.Equal(t, 200, w.Code)

		var body map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &body)
		assert.NilError(t, err)
		assert.Equal(t, int(body["limit"].(float64)), 3)
		assert.Equal(t, int(body["offset"].(float64)), 1)

		bodyUsers := body["users"].([]any)
		assert.Equal(t, len(bodyUsers), 3)

		slog.Debug(fmt.Sprintln(bodyUsers))

		for i := range len(bodyUsers) {
			bodyUser := bodyUsers[i].(map[string]interface{})

			testUserIdx := i + query["offset"].(int)
			assert.Equal(t, bodyUser["name"], testUsers[testUserIdx].UserName)
			assert.Equal(t, bodyUser["role"], testUsers[testUserIdx].Role)
			assert.Equal(t, bodyUser["devID"], testUsers[testUserIdx].DeviceID)
		}
	}
}

func successListEmpty(r *router.Engine) func(*testing.T) {
	return func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", USER_API, nil)
		r.ServeHTTP(w, req)

		assert.Equal(t, 200, w.Code)
		var body map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &body)
		assert.Equal(t, err, nil)
		assert.Equal(t, int(body["limit"].(float64)), 15)
		assert.Equal(t, int(body["offset"].(float64)), 0)
		assert.Equal(t, body["users"], nil)
	}
}
