package user

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"testing-demo/internal/db"
	"testing-demo/internal/logging"
	"testing-demo/internal/middleware"
)

var ts *httptest.Server

func TestMain(m *testing.M) {
	logging.Configure("debug", "test", []string{"stdout"}, []string{"stdout"})

	err := db.Connect()
	if err != nil {
		logging.GetLogger().Panic(err)
	}

	r := gin.New()
	r.Use(middleware.ContextLogger(logging.GetLogger()))

	apiRouter := r.Group("/api")
	apiRouter.Use(middleware.ErrorHandler())
	RegisterHandlers(apiRouter)

	ts = httptest.NewServer(r)
	defer ts.Close()

	exitVal := m.Run()

	os.Exit(exitVal)
}

func TestUsersPut(t *testing.T) {
	client := &http.Client{}
	newYork, _ := time.LoadLocation("America/New_York")

	t.Run("The correct results are returned when the inputs are valid", func(t *testing.T) {
		uri := fmt.Sprintf("%s/%s", ts.URL, "api/users")
		body := []byte(`[
  {
    "user_id": 1,
    "name": "Joe Smith",
    "DOB": "1983-05-12",
    "created_on": 1642612034
  },
  {
    "user_id": 2,
    "name": "Jane Doe",
    "DOB": "1990-08-06",
    "created_on": 1642612034
  }
]`)
		req, err := http.NewRequest("PUT", uri, bytes.NewBuffer(body))
		require.NoError(t, err)

		resp, err := client.Do(req)
		require.NoError(t, err)
		require.Equal(t, http.StatusOK, resp.StatusCode)

		expected := map[float64]map[string]interface{}{
			1: {
				"user_id":    float64(1),
				"name":       "Joe Smith",
				"DOB":        "Thursday",
				"created_on": time.Date(2022, 1, 19, 12, 7, 14, 0, newYork).Format(time.RFC3339),
			},
			2: {
				"user_id":    float64(2),
				"name":       "Jane Doe",
				"DOB":        "Monday",
				"created_on": time.Date(2022, 1, 19, 12, 7, 14, 0, newYork).Format(time.RFC3339),
			},
		}
		var actual []map[string]interface{}
		decoder := json.NewDecoder(resp.Body)
		err = decoder.Decode(&actual)
		require.NoError(t, err)
		for _, au := range actual {
			assert.Equal(t, expected[au["user_id"].(float64)], au)
		}
	})

	t.Run("The correct error results are returned when the inputs are not valid", func(t *testing.T) {
		uri := fmt.Sprintf("%s/%s", ts.URL, "api/users")
		body := []byte(`[
  {
    "user_id": 1,
    "name": "Joe Smith",
    "DOB": "1983-05-12",
    "created_on": 1642612034
  },
  {
    "user_id": 2,
    "name": "Jane Doe",
    "DOB": "199008-06",
    "created_on": 1642612034
  }
]`)
		req, err := http.NewRequest("PUT", uri, bytes.NewBuffer(body))
		require.NoError(t, err)

		resp, err := client.Do(req)
		require.NoError(t, err)
		require.Equal(t, http.StatusBadRequest, resp.StatusCode)

		expected := "'199008-06' can't be parsed as a date"

		var actual string
		decoder := json.NewDecoder(resp.Body)
		err = decoder.Decode(&actual)
		require.NoError(t, err)
		assert.Equal(t, expected, actual)
	})
}

func TestUserDatabase(t *testing.T) {
	t.Run("Can successfully create a user", func(t *testing.T) {
		suffix := time.Now().UnixMilli()
		up := &UserProfile{
			FirstName: "First",
			LastName:  "Last",
			Username:  fmt.Sprintf("user-%d", suffix),
			City:      "Anytown",
			ZipCode:   "12345",
		}
		err := up.Create()
		require.NoError(t, err)
	})

	t.Run("Can add new passwords for a user", func(t *testing.T) {
		suffix := time.Now().UnixMilli()
		up := &UserProfile{
			FirstName: "Password",
			LastName:  "Testing",
			Username:  fmt.Sprintf("user-%d", suffix),
			City:      "Anytown",
			ZipCode:   "12345",
		}
		err := up.Create()
		require.NoError(t, err)

		ph1 := &PasswordHistory{
			UserID:   up.ID,
			Password: fmt.Sprintf("password-1-%d", suffix),
		}
		err = ph1.Create()
		require.NoError(t, err)

		ph2 := &PasswordHistory{
			UserID:   up.ID,
			Password: fmt.Sprintf("password-2-%d", suffix),
		}
		err = ph2.Create()
		require.NoError(t, err)

		history, err := up.PasswordHistory()
		require.NoError(t, err)
		assert.Len(t, history, 2)

		activeHistory, err := up.ActivePasswordHistory()
		require.NoError(t, err)
		assert.Len(t, activeHistory, 1)
	})
}
