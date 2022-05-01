package smarthome_sdk

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

type Location struct {
	Filename string `json:"filename"`
	Line     uint   `json:"line"`
	Column   uint   `json:"column"`
	Index    uint   `json:"index"`
}

type HomescriptError struct {
	ErrorType string   `json:"errorType"`
	Location  Location `json:"location"`
	Message   string   `json:"message"`
}

type HomescriptResponse struct {
	Success  bool              `json:"success"`
	Exitcode int               `json:"exitCode"`
	Message  string            `json:"message"`
	Output   string            `json:"output"`
	Errors   []HomescriptError `json:"error"`
}

// Executes a string of homescript code on the Smarthome-server
// Returns a Homescript-response and an error
// The error is meant to indicate a failure of communication, not a failure of execution
func (c *Connection) RunHomescript(code string, timeout time.Duration) (response HomescriptResponse, err error) {
	if !c.ready {
		return HomescriptResponse{}, ErrNotInitialized
	}
	body, err := json.Marshal(struct {
		Code string `json:"code"`
	}{
		Code: code,
	})
	if err != nil {
		return HomescriptResponse{}, err
	}
	req, err := c.prepareRequest("/api/homescript/run/live", Post, body)
	if err != nil {
		return HomescriptResponse{}, err
	}
	client := &http.Client{
		Timeout: timeout,
	}
	res, err := client.Do(req)
	if err != nil {
		return HomescriptResponse{}, ErrConnFailed
	}
	defer res.Body.Close()

	switch res.StatusCode {
	// Either the script has executed successfully or it has terminated abnormally
	case 200, 500:
		resBody, err := ioutil.ReadAll(res.Body)
		if err != nil {
			return HomescriptResponse{}, ErrReadResponseBody
		}
		var parsedBody HomescriptResponse
		if err := json.Unmarshal(resBody, &parsedBody); err != nil {
			return HomescriptResponse{}, ErrReadResponseBody
		}
		return parsedBody, nil
	case 401:
		return HomescriptResponse{}, ErrUnauthorized
	case 403:
		return HomescriptResponse{}, ErrPermissionDenied
	}
	return HomescriptResponse{}, fmt.Errorf("unknown response code: %s", res.Status)

}
