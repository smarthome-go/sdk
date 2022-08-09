package sdk

import (
	"net/http"
	"net/url"
)

type AuthMethod uint8

// Specifies how the library handles authentication on every request
const (
	/** No authentication will send every request without any form of user-authentication
	- Can be used in a context which does not require authentication, for example, listing the switches
	- Not reccomended in most cases due to the strict data protection of Smarthome
	*/
	AuthMethodNone AuthMethod = iota

	/** Cookie authentication relies on a cookie-store which sends a authentication cookie at every request
	- Faster response-time: The server does not need to revalidate the user's credentials on every request
	- Static: If the Smarthome server is restarted, the stored cookie becomes invalid and the communication breaks
	- Not recommended for long-running applications
	*/
	AuthMethodCookiePassword
	// Uses a token to connect instead of the username and password whilst using cookie authentication
	AuthMethodCookieToken

	/** URL-query authentication adds `?username=foo&password=bar` or `?token=foo` to every requested URL
	- Slower response-time: The server needs to revalidate the user's credentials on every request
	- Dynamic: The connection will remain in a working condition after the smarthome server has been restarted
	- Recommended for long-running applications
	*/
	AuthMethodQueryPassword
	// Uses a token to connect instead of the username and password whilst using query authentication
	AuthMethodQueryToken
)

type Connection struct {
	// Stores credentials used by the SDK
	credStore credStore
	// The base URL which will be used to create all request
	SmarthomeURL *url.URL
	// Stores which authentication mode will be used
	authMethod AuthMethod
	// The cookie-store which will be used if `AuthMethodCookie` is set
	// The store is written to once in the login function
	// Every request will access the store in order to include the cookie in the request
	sessionCookie *http.Cookie
	// Used internally to specify if the connection is ready to be used
	ready bool
	// Stores the version of the Smarthome server in order to avoid using the `Version` function multiple times
	SmarthomeVersion string
	// Stores the GO version on which the Smarthome server runs on
	SmarthomeGoVersion string
}

// Saves the username - password combination
type credStore struct {
	/** Token authentication
	- Will be used when `AuthMethod-XX-Token` is set
	*/
	Token string

	/** User authentication
	- Will be used when `AuthMethod-XXX-Password` is set
	*/
	Username string
	Password string
}
