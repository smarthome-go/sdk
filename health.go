package sdk

import (
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

// Can be used in order to check if the Smarthome server is reachable and responds
// Returns an ENUM type indicating the overall health status of the server
func (c *Connection) HealthCheck() (status HealthStatus, err error) {
	u := c.SmarthomeURL
	u.Path = "health"

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
