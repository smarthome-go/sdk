package sdk

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/sessions"
	"github.com/stretchr/testify/assert"
)

func TestInit(t *testing.T) {
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

	r.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	ts := httptest.NewServer(r)
	defer ts.Close()
	// End test server

	c, err := NewConnection(ts.URL, AuthMethodCookiePassword)
	assert.NoError(t, err)
	assert.NoError(t, c.UserLogin("test", "test"))
	assert.Equal(t, c.sessionCookie.Name, "session")
	assert.True(t, c.ready)

	ts.Close()
	c3, err := NewConnection("http://not-reachable.local", AuthMethodNone)
	assert.NoError(t, err)
	assert.Error(t, c3.Connect())
	assert.EqualError(t, c3.Connect(), ErrConnFailed.Error())
}
