package sdk

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/sessions"
	"github.com/stretchr/testify/assert"
)

func TestInit(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		sessionStore := sessions.NewCookieStore([]byte("key"))
		session, _ := sessionStore.Get(r, "session")
		session.Values["valid"] = true
		session.Values["username"] = "test"
		assert.NoError(t, session.Save(r, w))
		w.WriteHeader(http.StatusNoContent)
	}))
	defer ts.Close()
	c, err := NewConnection(ts.URL, AuthMethodCookie)
	assert.NoError(t, err)
	assert.NoError(t, c.Connect("test", "test"))
	assert.Equal(t, c.SessionCookie.Name, "session")
	assert.True(t, c.ready)

	c2, err := NewConnection(ts.URL, AuthMethodNone)
	assert.NoError(t, err)
	assert.NoError(t, c2.Connect("test", "test"))
	assert.Equal(t, &http.Cookie{}, c2.SessionCookie)
	assert.True(t, c.ready)

	ts.Close()
	c3, err := NewConnection("http://not-reachable.local", AuthMethodNone)
	assert.NoError(t, err)
	assert.EqualError(t, c3.Connect("", ""), ErrConnFailed.Error())
	assert.Error(t, c3.Connect("", ""))
}
