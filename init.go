package smarthome_sdk

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/url"
)

// Creates a new connection
// First argument specifies the base URL of the target Smarthome-server
// Second argument specifies how to handle authentication
func New(smarthomeURL string, authMethod AuthMethod) (*Connection, error) {
	u, err := url.Parse(smarthomeURL)
	if err != nil {
		return nil, ErrInvalidURL
	}
	return &Connection{
		SmarthomeURL:  u,
		AuthMethod:    authMethod,
		SessionCookie: &http.Cookie{},
	}, nil
}

// If the authentication mode is set to `AuthMethodNone`, both arguments can be set to nil
// Otherwise, username and password are required to login
func (c *Connection) Init(username string, password string) error {

	u := c.SmarthomeURL
	u.Path = "health"
	// Check if the URL is working
	res, err := http.Get(u.String())
	if err != nil {
		return ErrConnFailed
	}
	if res.StatusCode == http.StatusServiceUnavailable {
		return ErrServiceUnavailable
	}
	if c.AuthMethod == AuthMethodNone {
		c.ready = true
		return nil
	}
	// If authentication is set to either `AuthMethodCookie` or `AuthMethodQuery`
	// A login request is sent in order to validate the provided credentials
	cookie, err := c.doLogin()
	if err != nil {
		return err
	}
	c.SessionCookie = cookie
	c.ready = true
	return nil
}

func (c *Connection) doLogin() (*http.Cookie, error) {
	u := c.SmarthomeURL
	u.Path = "/api/login"

	body, err := json.Marshal(struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}{
		Username: c.Username,
		Password: c.Password,
	})
	if err != nil {
		return nil, err
	}

	r, err := http.NewRequest(http.MethodPost, u.String(), bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}

	client := &http.Client{}

	res, err := client.Do(r)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	switch res.StatusCode {
	case 204:
		for _, cookie := range res.Cookies() {
			if cookie.Name == "session" {
				return cookie, nil
			}
		}
		return nil, ErrNoCookiesSent
	case 401:
		return nil, ErrInvalidCredentials
	case 500:
		return nil, ErrInternalServerError
	case 503:
		return nil, ErrServiceUnavailable
	default:
		return nil, ErrUnknown
	}
}
