package smarthome_sdk

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
)

type DBStatus struct {
	OpenConnections int `json:"openConnections"`
	InUse           int `json:""`
	Idle            int `json:""`
}

type PowerJob struct {
	Id         int64  `json:"id"`
	SwitchName string `json:"switchName"`
	Power      bool   `json:"power"`
}

type JobResult struct {
	Id    int64 `json:"id"`
	Error error `json:"error"`
}

// Is returned when the debug information is requested
type DebugInfoData struct {
	ServerVersion          string      `json:"version"`
	DatabaseOnline         bool        `json:"databaseOnline"`
	DatabaseStats          DBStatus    `json:"databaseStats"`
	CpuCores               uint8       `json:"cpuCores"`
	Goroutines             uint16      `json:"goroutines"`
	GoVersion              string      `json:"goVersion"`
	MemoryUsage            uint16      `json:"memoryUsage"`
	PowerJobCount          uint16      `json:"powerJobCount"`
	PowerJobWithErrorCount uint16      `json:"lastPowerJobErrorCount"`
	PowerJobs              []PowerJob  `json:"powerJobs"`
	PowerJobResults        []JobResult `json:"powerJobResults"`
}

// Retrieves debugging information from the smarthome server
/** Errors
- nil
- ErrNotInitialized
- ErrConnFailed
- ErrReadResponseBody
- ErrInvalidCredentials
- ErrServiceUnavailable
- PrepareRequest errors
- ErrPermissionDenied
*/
func (c *Connection) GetDebugInfo() (info DebugInfoData, err error) {
	if !c.ready {
		return DebugInfoData{}, ErrNotInitialized
	}
	req, err := c.prepareRequest("/api/debug", Get, nil)
	if err != nil {
		return DebugInfoData{}, err
	}
	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return DebugInfoData{}, ErrConnFailed
	}
	switch res.StatusCode {
	case 200:
		resBody, err := ioutil.ReadAll(res.Body)
		if err != nil {
			return DebugInfoData{}, ErrReadResponseBody
		}
		var parsedBody DebugInfoData
		if err := json.Unmarshal(resBody, &parsedBody); err != nil {
			return DebugInfoData{}, ErrReadResponseBody
		}
		return parsedBody, nil
	case 401:
		return DebugInfoData{}, ErrInvalidCredentials
	case 403:
		return DebugInfoData{}, ErrPermissionDenied
	case 503:
		return DebugInfoData{}, ErrServiceUnavailable
	}
	return DebugInfoData{}, nil
}
