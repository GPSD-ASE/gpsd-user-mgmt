package tests

import (
	"encoding/json"
	"gpsd-user-mgmt/src/db"
	"gpsd-user-mgmt/src/router"
	"gpsd-user-mgmt/src/user"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"gotest.tools/assert"
)

func successCreate(r *router.Engine) func(*testing.T) {
	return func(t *testing.T) {
		testUsers := []user.User{{
			UserName: "Test",
			DeviceID: "123",
			Role:     "reporter",
		}, {
			UserName: "Test2",
			DeviceID: "1234",
			Role:     "admin",
		}}
		defer db.EmptyDatabase()

		for _, testUser := range testUsers {

			w := httptest.NewRecorder()
			payload, _ := json.Marshal(testUser)
			req, _ := http.NewRequest(
				"POST",
				USER_API,
				strings.NewReader(string(payload)),
			)
			r.ServeHTTP(w, req)

			assert.Equal(t, 200, w.Code)

			var body map[string]interface{}
			err := json.Unmarshal(w.Body.Bytes(), &body)
			assert.NilError(t, err)
			assert.Equal(t, body["message"], "User created successfully")

			userBody := body["user"].(map[string]interface{})

			assert.Equal(t, userBody["name"], testUser.UserName)
			assert.Equal(t, userBody["role"], testUser.Role)
			assert.Equal(t, userBody["devID"], testUser.UserName)
		}
	}
}
