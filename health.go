package sdk

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type HealthStatus uint

const (
	StatusUnknown           HealthStatus = iota //Default return value if a request error occurs
	StatusHealthy                               // All systems are in a operational state
	StatusPartiallyDegraded                     // Some systems are degraded, for example one out of more Hardware nodes have failed recently
	StatusDegraded                              // The database connection has failed, Smarthome cannot be used in this state
)

type VersionResponse struct {
	Version   string `json:"version"`
	GoVersion string `json:"goVersion"`
}

// Can be used in order to check if the Smarthome server is reachable and responds
// Returns an ENUM type indicating the overall health status of the server
/** Errors
- nil
- ErrConnFailed
- ErrReadResponseBody
*/
func (c *Connection) HealthCheck() (status HealthStatus, err error) {
	u := c.SmarthomeURL
	u.Path = "/health"

	// Check if the base URL is working and the server is reachable
	res, err := http.Get(u.String())
	if err != nil {
		return StatusUnknown, ErrConnFailed
	}
	switch res.StatusCode {
	case 200:
		return StatusHealthy, nil
	case 502:
		return StatusPartiallyDegraded, nil
	case 503:
		return StatusDegraded, nil
	}
	return StatusUnknown, fmt.Errorf("unknown response code: %s", res.Status)
}

// Can be used to retrieve the current version of the Smarthome server
/** Errors
- nil
- ErrConnFailed
- ErrReadResponseBody
*/
func (c *Connection) Version() (version VersionResponse, err error) {
	u := c.SmarthomeURL
	u.Path = "/api/version"

	res, err := http.Get(u.String())
	if err != nil {
		return VersionResponse{}, ErrConnFailed
	}

	decoder := json.NewDecoder(res.Body)
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&version); err != nil {
		return VersionResponse{}, ErrReadResponseBody
	}

	return version, nil
}
