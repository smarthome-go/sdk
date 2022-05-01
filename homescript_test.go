package smarthome_sdk

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gorilla/sessions"
	"github.com/stretchr/testify/assert"
)

func TestHomescript(t *testing.T) {
	r := http.NewServeMux()

	r.HandleFunc("/api/login", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		sessionStore := sessions.NewCookieStore([]byte("key"))
		session, _ := sessionStore.Get(r, "session")
		session.Values["valid"] = true
		session.Values["username"] = "test"
		assert.NoError(t, session.Save(r, w))
		w.WriteHeader(http.StatusNoContent)
	})

	r.HandleFunc("/api/homescript/run/live", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		assert.NoError(t, json.NewEncoder(w).Encode(HomescriptResponse{
			Success:  true,
			Exitcode: 0,
			Message:  "",
			Output:   "",
			Errors:   make([]HomescriptError, 0),
		}))
	})

	ts := httptest.NewServer(r)
	defer ts.Close()

	// Authenticate
	c, err := New(ts.URL, AuthMethodCookie)
	assert.NoError(t, err)
	assert.NoError(t, c.Init("test", "test"))
	assert.Equal(t, c.SessionCookie.Name, "session")
	assert.True(t, c.ready)

	// Test Homescript
	res, err := c.RunHomescript("print('_')", time.Second)
	assert.NoError(t, err)
	assert.Equal(t, HomescriptResponse{
		Success:  true,
		Exitcode: 0,
		Message:  "",
		Output:   "",
		Errors:   make([]HomescriptError, 0),
	}, res)
}
