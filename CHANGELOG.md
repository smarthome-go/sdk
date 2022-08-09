## Changelog for v0.17.0

- Added support for `token` authentication
- Completely refactored the way login and connection works
- Now compatible with Smarthome user authentication tokens and server version `0.0.54`

### Example Code
This code demonstrates the use and difference of the new authentication methods
```go
package main

import "github.com/smarthome-go/sdk"

const URL = "http://localhost:8082"

func main() {
	// === Example 1: Using token authentication === //
	c1, err := sdk.NewConnection(URL, sdk.AuthMethodQueryToken /* sdk.AuthMethodCookieToken */)
	if err != nil {
		panic(err.Error())
	}
	// Login with your user authentication token, to obtain one, visit `http://your-smarthome.box/profile`
	if err := c1.TokenLogin("650feaafc1487d18bd8c5a805363be96"); err != nil {
		panic(err.Error())
	}

	// === Example 2: Using username & password authentication === //
	c2, err := sdk.NewConnection(URL, sdk.AuthMethodQueryPassword /* sdk.AuthMethodCookiePassword */)
	if err != nil {
		panic(err.Error())
	}
	// Login with the usual username-password combination
	if err := c1.UserLogin("admin", "admin"); err != nil {
		panic(err.Error())
	}

	// => After login, each connection behaves idenically
	if err := c1.SetPower("s2", true); err != nil {
		panic(err.Error())
	}
	if err := c2.SetPower("s2", true); err != nil {
		panic(err.Error())
	}
}
```
