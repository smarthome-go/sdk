# SDK
 A GO package which makes communication to a Smarthome server simple.  
 It can be seen as a `API` wrapper for some commonly-used functions of the Smarthome server's `API`.

## Example Code Using the SDK (`v0.17.0`)
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
