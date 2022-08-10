package sdk

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"github.com/Masterminds/semver"
)

// Is sent by the server when token-login is used
// Required for setting the username when token authentication is used
type tokenLoginResponse struct {
	Username   string `json:"username"`
	TokenLabel string `json:"tokenLabel"`
}

// Creates a new connection
// First argument specifies the base URL of the target Smarthome-server
// Second argument specifies how to handle authentication
func NewConnection(
	smarthomeURL string,
	authMethod AuthMethod,
) (*Connection, error) {
	u, err := url.Parse(smarthomeURL)
	if err != nil {
		return nil, ErrInvalidURL
	}
	// Create and return a client
	return &Connection{
		SmarthomeURL:  u,
		authMethod:    authMethod,
		sessionCookie: &http.Cookie{},
	}, nil
}

// Can be used to connect when the authentication method is set to `None`
func (c *Connection) Connect() error {
	if c.authMethod != AuthMethodNone {
		return ErrInvalidFunctionAuthMethod
	}
	// Call the helper function
	return c.connectHelper()
}

// Can be used to connect when the authentication method is set to `Password-XXX`
func (c *Connection) UserLogin(username string, password string) error {
	if c.authMethod != AuthMethodQueryPassword && c.authMethod != AuthMethodCookiePassword {
		return ErrInvalidFunctionAuthMethod
	}
	// Set the internal credentials using the parameters
	c.credStore.Username = username
	c.credStore.Password = password
	// Call the helper function
	return c.connectHelper()
}

// Can be used to connect when the authentication method is set to `Token-XXX`
func (c *Connection) TokenLogin(token string) error {
	if c.authMethod != AuthMethodQueryToken && c.authMethod != AuthMethodCookieToken {
		return ErrInvalidFunctionAuthMethod
	}
	// Set the internal token to the parameter
	c.credStore.Token = token
	// Call the helper function
	return c.connectHelper()
}

// If the authentication mode is set to `AuthMethodNone`, both arguments can be set to nil
// Otherwise, username and password are required to login
func (c *Connection) connectHelper() error {

	// Retrieve the server's version
	version, err := c.Version()
	if err != nil {
		return err
	}

	// Set the version in the connection
	// Is already set here so it can be used in error messages as `c.SmarthomeVersion`
	c.SmarthomeVersion = version.Version
	c.SmarthomeGoVersion = version.GoVersion

	// Check Smarthome version compatibility
	supportedV, err := semver.NewConstraint(fmt.Sprintf("^%s", MinSmarthomeVersion))
	if err != nil {
		// This must not happen (tests)
		// If this happens, the best thing is to abort the connection
		return ErrInvalidVersion
	}

	currentV, err := semver.NewVersion(version.Version)
	if err != nil {
		// This must also not happen
		// If this happens, the best thing is to abort the connection
		return ErrInvalidVersion
	}

	// Perform the version comparison
	if !supportedV.Check(currentV) {
		// Would not be supported
		return ErrUnsupportedVersion
	}

	switch c.authMethod {
	// If the connection does not use authentication, it can be marked as ready
	case AuthMethodNone:
		c.ready = true
		return nil

	// If the authentication mode is set to `AuthMethodQueryToken`, validate the token and mark the connection as ready
	case AuthMethodQueryToken:
		_, tokenData, err := c.doLogin()
		if err != nil {
			return err
		}
		c.tokenClientName = tokenData.TokenLabel
		c.credStore.Username = tokenData.Username
		c.ready = true
		return nil
	// If the authentication mode is set to `AuthMethodQueryPassword`, validate the user's credentials and mark the connection as ready
	case AuthMethodQueryPassword:
		_, _, err := c.doLogin()
		if err != nil {
			return err
		}
		c.ready = true
		return nil
	// If the authentication mode is set to `AuthMethodCookieToken`, use the token to obtain a session cookie
	case AuthMethodCookieToken:
		_, tokenData, err := c.doLogin()
		if err != nil {
			return err
		}
		c.tokenClientName = tokenData.TokenLabel
		c.credStore.Username = tokenData.Username
		c.ready = true
		return nil
	// If the authentication mode is set to `AuthMethodCookiePassword`, use the user's credentials to obtain a session cookie
	case AuthMethodCookiePassword:
		cookie, _, err := c.doLogin()
		if err != nil {
			return err
		}
		c.sessionCookie = cookie
		c.ready = true
		return nil

	default:
		panic("unreachable")
	}
}

// Used internally to send a login request
// When the authentication mode is set to `AuthMethodCookie-XXX`, the response cookie is saved
// However, for `AuthMethodQuery-XXX`, it serves the purpose of validating the provided credentials beforehand
// If the authentication mode is sey set to `AuthMethodNone`, the function call is omitted
func (c *Connection) doLogin() (
	*http.Cookie,
	*tokenLoginResponse,
	error,
) {
	u := c.SmarthomeURL
	// The default path is the user login
	u.Path = "/api/login"
	// If authentication should use a token, change the path
	if c.authMethod == AuthMethodQueryToken || c.authMethod == AuthMethodCookieToken {
		u.Path = "/api/login/token"
	}

	var loginBody []byte
	var loginBodyErr error

	if c.authMethod == AuthMethodQueryPassword || c.authMethod == AuthMethodCookiePassword {
		loginBody, loginBodyErr = json.Marshal(struct {
			Username string `json:"username"`
			Password string `json:"password"`
		}{
			Username: c.credStore.Username,
			Password: c.credStore.Password,
		})
		if loginBodyErr != nil {
			return nil, nil, loginBodyErr
		}
	} else if c.authMethod == AuthMethodQueryToken || c.authMethod == AuthMethodCookieToken {
		loginBody, loginBodyErr = json.Marshal(struct {
			Token string `json:"token"`
		}{
			Token: c.credStore.Token,
		})
		if loginBodyErr != nil {
			return nil, nil, loginBodyErr
		}
	} else {
		panic("unreachable")
	}
	// Create a login request
	r, err := http.NewRequest(
		http.MethodPost,
		u.String(),
		bytes.NewBuffer(loginBody),
	)
	if err != nil {
		return nil, nil, err
	}
	// Perform the login request
	client := &http.Client{}
	res, err := client.Do(r)
	if err != nil {
		return nil, nil, err
	}
	defer res.Body.Close()
	switch res.StatusCode {
	case 200:
		if c.authMethod == AuthMethodQueryPassword || c.authMethod == AuthMethodCookiePassword {
			// should not happen: this is a bug
			panic("unreachable")
		}
		// Attempt to decode the response body
		resBody, err := io.ReadAll(res.Body)
		if err != nil {
			return nil, nil, ErrReadResponseBody
		}
		var parsedBody tokenLoginResponse
		if err := json.Unmarshal(resBody, &parsedBody); err != nil {
			return nil, nil, ErrReadResponseBody
		}
		var returnCookie *http.Cookie
		for _, cookie := range res.Cookies() {
			if cookie.Name == "session" {
				returnCookie = cookie
				break
			}
		}
		if returnCookie == nil {
			return nil, nil, ErrNoCookiesSent
		}
		return returnCookie, &parsedBody, nil
	case 204:
		for _, cookie := range res.Cookies() {
			if cookie.Name == "session" {
				return cookie, nil, nil
			}
		}
		return nil, nil, ErrNoCookiesSent
	case 401:
		return nil, nil, ErrInvalidCredentials
	case 500:
		return nil, nil, ErrInternalServerError
	case 503:
		return nil, nil, ErrServiceUnavailable
	default:
		return nil, nil, fmt.Errorf("unknown response code: %s", res.Status)
	}
}

// Works on every authentication method except `None`
func (c *Connection) GetUsername() (string, error) {
	if c.authMethod == AuthMethodNone {
		return "", ErrInvalidFunctionAuthMethod
	}
	return c.credStore.Username, nil
}

// Only works on token-based authentication methods
func (c *Connection) GetTokenClientLabel() (string, error) {
	if c.authMethod != AuthMethodQueryToken && c.authMethod != AuthMethodCookieToken {
		return "", ErrInvalidFunctionAuthMethod
	}
	return c.tokenClientName, nil
}
