package smarthome_sdk

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

// Creates a new connection
// First argument specifies the base URL of the target Smarthome-server
// Second argument specifies how to handle authentication
func NewConnection(smarthomeURL string, authMethod AuthMethod) (*Connection, error) {
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
func (c *Connection) Connect(username string, password string) error {
	c.Username = username
	c.Password = password

	status, err := c.HealthCheck()
	if err != nil {
		return err
	}
	// Check if the healthcheck failed to an extend that authentication will not be possible
	if status == StatusDegraded || status == StatusUnknown {
		return ErrServiceUnavailable
	}
	// If the connection does not use authentication, it can be marked as ready
	if c.AuthMethod == AuthMethodNone {
		c.ready = true
		return nil
	}
	// If the authentication mode is set to `AuthMethodQuery`, validate the user's credentials and mark the connection as ready
	if c.AuthMethod == AuthMethodQuery {
		_, err := c.doLogin()
		if err != nil {
			return err
		}
		c.ready = true
		return nil
	}
	// If the authentication mode is set to `AuthMethodCookie`, validate the user's credentials and save the cookie
	if c.AuthMethod == AuthMethodCookie {
		cookie, err := c.doLogin()
		if err != nil {
			return err
		}
		c.SessionCookie = cookie
		c.ready = true
		return nil
	}
	// Unreachable
	return nil
}

// Used internally to send a login request
// When the authentication mode is set to `AuthMethodCookie`, the response cookie is saved
// However, for `AuthMethodQuery`, it serves the purpose of validating the provided credentials beforehand
// If the authentication mode is sey set to `AuthMethodNone`, the function call is omitted
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
		return nil, fmt.Errorf("unknown response code: %s", res.Status)
	}
}
