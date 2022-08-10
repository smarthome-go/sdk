## Changelog for v0.19.0

- Added username & token label fetching when logging in via an authentication token
- This means you can obtain the username even though a token was used to authenticate
- In order to get the username and token label, two additional public functions have been implemented (*demonstrated in the code example below*)
- Removed usages of the deprecated `ioutil` package
- This release is required because one error change inside the error "`ENUMS`"
- This is a significant change and should be made public.
- Because of the above facts, this release exists

### Code Example (*Obtaining username + token label*)
```go
package main

import (
	"fmt"

	"github.com/smarthome-go/sdk"
)

const URL = "http://smarthome.box"

func main() {
	// === Example 1: Obtaining username & token-label === //
	c, err := sdk.NewConnection(URL, sdk.AuthMethodQueryToken /* sdk.AuthMethodCookieToken */)
	if err != nil {
		panic(err.Error())
	}
	// Login with your user authentication token, to obtain one, visit `http://your-smarthome.box/profile`
	if err := c.TokenLogin("807b2eded585803ff287c295afe23d1a"); err != nil {
		panic(err.Error())
	}
	fmt.Println(c.GetUsername())         // This function works on all auth methods except `None`
	fmt.Println(c.GetTokenClientLabel()) // This function only works when token auth is used
}
```
