package sdk

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

type HTTPMethod string

const (
	Get    HTTPMethod = "GET"
	Post   HTTPMethod = "POST"
	Put    HTTPMethod = "PUT"
	Delete HTTPMethod = "DELETE"
)

// Used internally in order to act as a middleware to add authentication to a requested URI
func (c *Connection) prepareRequest(path string, method HTTPMethod, body interface{}) (*http.Request, error) {
	// Creates a local copy of the smarthome base URL, then sets the path
	u := c.SmarthomeURL
	u.Path = path

	// If the authentication mode is set to `AuthMethodQueryPassword`, encode username and password and attach it to the URL
	if c.authMethod == AuthMethodQueryPassword {
		query := u.Query()
		query.Set("username", c.userPasswordData.Username)
		query.Set("password", c.userPasswordData.Password)
		u.RawQuery = query.Encode()
	} else if c.authMethod == AuthMethodQueryToken {
		query := u.Query()
		query.Set("token", c.token)
		u.RawQuery = query.Encode()
	}

	// If a body is specified, encode it to JSON
	encodedBody := make([]byte, 0)
	if body != nil {
		b, err := json.Marshal(body)
		if err != nil {
			return nil, err
		}
		encodedBody = b
	}

	// Creates the request
	r, err := http.NewRequest(string(method), u.String(), bytes.NewBuffer(encodedBody))
	if err != nil {
		return nil, err
	}

	// If the authentication mode is set to `AuthMethodCookiePassword`, add the cookie to the request
	if c.authMethod == AuthMethodCookiePassword {
		r.AddCookie(c.sessionCookie)
	}

	// Set `Content-Type` and `User-Agent`
	r.Header.Set("Content-Type", "application/json")
	r.Header.Set("User-Agent", fmt.Sprintf("SmarthomeSDK/%s", Version))
	return r, nil
}
