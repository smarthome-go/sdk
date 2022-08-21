package sdk

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"
)

// Is sent to the server in order to provide an equivalent to arguments
type HomescriptRunArgRequest struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

// Is used in order to request execution of a Homescript via its id
type RunHomescriptByIdRequest struct {
	Id   string                    `json:"id"`
	Args []HomescriptRunArgRequest `json:"args"`
}

// Is used in order to request execution of an arbitrary Homescript string
type RunHomescriptStringRequest struct {
	Code string                    `json:"code"`
	Args []HomescriptRunArgRequest `json:"args"`
}

// Specifies where the Homescript error occurred
type ErrorLocation struct {
	Filename string `json:"filename"`
	Line     uint   `json:"line"`
	Column   uint   `json:"column"`
	Index    uint   `json:"index"`
}

// Contains information about why a Homescript terminated
type HomescriptError struct {
	ErrorType string        `json:"errorType"`
	Location  ErrorLocation `json:"location"`
	Message   string        `json:"message"`
}

// Under normal conditions, Smarthome will return such a response
type HomescriptResponse struct {
	Success  bool              `json:"success"`
	Id       string            `json:"id"`
	Exitcode int               `json:"exitCode"`
	Message  string            `json:"message"`
	Output   string            `json:"output"`
	Errors   []HomescriptError `json:"error"`
}

// Represents a Homescript entity
type Homescript struct {
	Owner string         `json:"owner"`
	Data  HomescriptData `json:"data"`
}

// The most important data which a Homescript contains
type HomescriptData struct {
	Id                  string `json:"id"`
	Name                string `json:"name"`
	Description         string `json:"description"`
	QuickActionsEnabled bool   `json:"quickActionsEnabled"`
	SchedulerEnabled    bool   `json:"schedulerEnabled"`
	Code                string `json:"code"`
	MDIcon              string `json:"mdIcon"`
	Workspace           string `json:"workspace"`
}

// Used for creating a new script or modifying an existing one
type HomescriptRequest struct {
	Id                  string `json:"id"`
	Name                string `json:"name"`
	Description         string `json:"description"`
	QuickActionsEnabled bool   `json:"quickActionsEnabled"`
	SchedulerEnabled    bool   `json:"schedulerEnabled"`
	Code                string `json:"code"`
	MDIcon              string `json:"mdIcon"`
}

// Executes a string of Homescript code on the Smarthome-server
// Returns a Homescript-response and an error
// The error is meant to indicate a failure of communication, not a failure of execution
// Normally, a `ErrConnFailed` indicates that the server is not reachable, however if other requests work
// a `ErrConnFailed` could indicate a request-timeout. In this case, check if you need to increase the timeout
/** Errors
- nil
- ErrNotInitialized
- ErrConnFailed
- ErrReadResponseBody
- ErrInvalidCredentials
- ErrPermissionDenied
- PrepareRequest errors
- Unknown
*/
func (c *Connection) RunHomescriptCode(code string, args map[string]string, timeout time.Duration) (response HomescriptResponse, err error) {
	if !c.ready {
		return HomescriptResponse{}, ErrNotInitialized
	}
	argsTemp := make([]HomescriptRunArgRequest, 0)
	for key, value := range args {
		argsTemp = append(argsTemp, HomescriptRunArgRequest{
			Key:   key,
			Value: value,
		})
	}
	req, err := c.prepareRequest("/api/homescript/run/live", Post, RunHomescriptStringRequest{
		Code: code,
		Args: argsTemp,
	})
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
		resBody, err := io.ReadAll(res.Body)
		if err != nil {
			return HomescriptResponse{}, ErrReadResponseBody
		}
		var parsedBody HomescriptResponse
		if err := json.Unmarshal(resBody, &parsedBody); err != nil {
			return HomescriptResponse{}, ErrReadResponseBody
		}
		return parsedBody, nil
	case 401:
		return HomescriptResponse{}, ErrInvalidCredentials
	case 403:
		return HomescriptResponse{}, ErrPermissionDenied
	}
	return HomescriptResponse{}, fmt.Errorf("unknown response code: %s", res.Status)
}

// Runs Homescript by id on the Smarthome-server
// Returns a Homescript-response and an error
// The error is meant to indicate a failure of communication, not a failure of execution
// Normally, a `ErrConnFailed` indicates that the server is not reachable, however if other requests work
// a `ErrConnFailed` could indicate a request-timeout. In this case, check if you need to increase the timeout
/** Errors
- nil
- ErrNotInitialized
- ErrConnFailed
- ErrReadResponseBody
- ErrInvalidCredentials
- ErrPermissionDenied
- PrepareRequest errors
- Unknown
*/
func (c *Connection) RunHomescriptById(id string, args map[string]string, timeout time.Duration) (response HomescriptResponse, err error) {
	if !c.ready {
		return HomescriptResponse{}, ErrNotInitialized
	}
	argsTemp := make([]HomescriptRunArgRequest, 0)
	for key, value := range args {
		argsTemp = append(argsTemp, HomescriptRunArgRequest{
			Key:   key,
			Value: value,
		})
	}
	req, err := c.prepareRequest("/api/homescript/run", Post, RunHomescriptByIdRequest{
		Id:   id,
		Args: argsTemp,
	})
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
		resBody, err := io.ReadAll(res.Body)
		if err != nil {
			return HomescriptResponse{}, ErrReadResponseBody
		}
		var parsedBody HomescriptResponse
		if err := json.Unmarshal(resBody, &parsedBody); err != nil {
			return HomescriptResponse{}, ErrReadResponseBody
		}
		return parsedBody, nil
	case 401:
		return HomescriptResponse{}, ErrInvalidCredentials
	case 403:
		return HomescriptResponse{}, ErrPermissionDenied
	case 422:
		return HomescriptResponse{}, ErrUnprocessableEntity
	}
	return HomescriptResponse{}, fmt.Errorf("unknown response code: %s", res.Status)
}

// Lints a string of Homescript code on the Smarthome-server
// Returns a Homescript-response and an error
// The error is meant to indicate a failure of communication, not a failure of linting
/** Errors
- nil
- ErrNotInitialized
- ErrConnFailed
- ErrReadResponseBody
- ErrInvalidCredentials
- ErrPermissionDenied
- PrepareRequest errors
- Unknown
*/
func (c *Connection) LintHomescriptCode(code string, args map[string]string, timeout time.Duration) (response HomescriptResponse, err error) {
	if !c.ready {
		return HomescriptResponse{}, ErrNotInitialized
	}
	argsTemp := make([]HomescriptRunArgRequest, 0)
	for key, value := range args {
		argsTemp = append(argsTemp, HomescriptRunArgRequest{
			Key:   key,
			Value: value,
		})
	}
	req, err := c.prepareRequest("/api/homescript/lint/live", Post, RunHomescriptStringRequest{
		Code: code,
		Args: argsTemp,
	})
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
	// Either the script has been linted successfully or errors were detected
	case 200, 500:
		resBody, err := io.ReadAll(res.Body)
		if err != nil {
			return HomescriptResponse{}, ErrReadResponseBody
		}
		var parsedBody HomescriptResponse
		if err := json.Unmarshal(resBody, &parsedBody); err != nil {
			return HomescriptResponse{}, ErrReadResponseBody
		}
		return parsedBody, nil
	case 401:
		return HomescriptResponse{}, ErrInvalidCredentials
	case 403:
		return HomescriptResponse{}, ErrPermissionDenied
	}
	return HomescriptResponse{}, fmt.Errorf("unknown response code: %s", res.Status)
}

// Lints a Homescript by id
// Returns a Homescript-response and an error
// The error is meant to indicate a failure of communication, not a failure of linting
/** Errors
- nil
- ErrNotInitialized
- ErrConnFailed
- ErrReadResponseBody
- ErrInvalidCredentials
- ErrPermissionDenied
- PrepareRequest errors
- Unknown
*/
func (c *Connection) LintHomescriptById(id string, args map[string]string, timeout time.Duration) (response HomescriptResponse, err error) {
	if !c.ready {
		return HomescriptResponse{}, ErrNotInitialized
	}
	argsTemp := make([]HomescriptRunArgRequest, 0)
	for key, value := range args {
		argsTemp = append(argsTemp, HomescriptRunArgRequest{
			Key:   key,
			Value: value,
		})
	}
	req, err := c.prepareRequest("/api/homescript/lint", Post, RunHomescriptByIdRequest{
		Id:   id,
		Args: argsTemp,
	})
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
	// Either the script has been linted successfully or errors were detected
	case 200, 500:
		resBody, err := io.ReadAll(res.Body)
		if err != nil {
			return HomescriptResponse{}, ErrReadResponseBody
		}
		var parsedBody HomescriptResponse
		if err := json.Unmarshal(resBody, &parsedBody); err != nil {
			return HomescriptResponse{}, ErrReadResponseBody
		}
		return parsedBody, nil
	case 401:
		return HomescriptResponse{}, ErrInvalidCredentials
	case 403:
		return HomescriptResponse{}, ErrPermissionDenied
	case 422:
		return HomescriptResponse{}, ErrUnprocessableEntity
	}
	return HomescriptResponse{}, fmt.Errorf("unknown response code: %s", res.Status)

}

// Creates a new Homescript which is owned by the current user
/** Errors
- nil
- ErrNotInitialized
- ErrConnFailed
- ErrReadResponseBody
- ErrInvalidCredentials
- ErrPermissionDenied
- PrepareRequest errors
- ErrUnprocessableEntity (conflicting id / invalid data)
- Unknown
*/
func (c *Connection) CreateHomescript(data HomescriptRequest) error {
	if !c.ready {
		return ErrNotInitialized
	}
	req, err := c.prepareRequest("/api/homescript/add", Post, data)
	if err != nil {
		return err
	}
	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return ErrConnFailed
	}
	defer res.Body.Close()
	switch res.StatusCode {
	case 200:
		return nil
	case 401:
		return ErrInvalidCredentials
	case 422:
		return ErrUnprocessableEntity
	case 403:
		return ErrPermissionDenied
	}
	return fmt.Errorf("unknown response code: %s", res.Status)
}

// Modifies an existing Homescript which is owned by the current user
/** Errors
- nil
- ErrNotInitialized
- ErrConnFailed
- ErrReadResponseBody
- ErrInvalidCredentials
- ErrPermissionDenied
- PrepareRequest errors
- ErrUnprocessableEntity (invalid id / in valid data)
- Unknown
*/
func (c *Connection) ModifyHomescript(data HomescriptRequest) error {
	if !c.ready {
		return ErrNotInitialized
	}
	req, err := c.prepareRequest("/api/homescript/modify", Put, data)
	if err != nil {
		return err
	}
	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return ErrConnFailed
	}
	defer res.Body.Close()
	switch res.StatusCode {
	case 200:
		return nil
	case 401:
		return ErrInvalidCredentials
	case 422:
		return ErrUnprocessableEntity
	case 403:
		return ErrPermissionDenied
	}
	return fmt.Errorf("unknown response code: %s", res.Status)
}

// Deletes an existing Homescript which is owned by the current user
/** Errors
- nil
- ErrNotInitialized
- ErrConnFailed
- ErrReadResponseBody
- ErrInvalidCredentials
- ErrPermissionDenied
- PrepareRequest errors
- ErrUnprocessableEntity (invalid id)
- ErrConflict (dependent automations)
- Unknown
*/
func (c *Connection) DeleteHomescript(id string) error {
	if !c.ready {
		return ErrNotInitialized
	}
	req, err := c.prepareRequest("/api/homescript/delete", Delete, struct {
		Id string `json:"id"`
	}{id})
	if err != nil {
		return err
	}
	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return ErrConnFailed
	}
	defer res.Body.Close()
	switch res.StatusCode {
	case 200:
		return nil
	case 401:
		return ErrInvalidCredentials
	case 409:
		return ErrConflict
	case 422:
		return ErrUnprocessableEntity
	case 403:
		return ErrPermissionDenied
	}
	return fmt.Errorf("unknown response code: %s", res.Status)
}

// Returns the metadata of a given homescript which is owned by the current user
/** Errors
- nil
- ErrNotInitialized
- ErrConnFailed
- ErrReadResponseBody
- ErrInvalidCredentials
- ErrPermissionDenied
- PrepareRequest errors
- ErrUnprocessableEntity (invalid id)
- Unknown
*/
func (c *Connection) GetHomescript(id string) (Homescript, error) {
	if !c.ready {
		return Homescript{}, ErrNotInitialized
	}
	req, err := c.prepareRequest(fmt.Sprintf("/api/homescript/get/%s", url.PathEscape(id)), Get, nil)
	if err != nil {
		return Homescript{}, err
	}
	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return Homescript{}, ErrConnFailed
	}
	defer res.Body.Close()
	switch res.StatusCode {
	case 200:
		resBody, err := io.ReadAll(res.Body)
		if err != nil {
			return Homescript{}, ErrReadResponseBody
		}
		var parsedBody Homescript
		if err := json.Unmarshal(resBody, &parsedBody); err != nil {
			return Homescript{}, ErrReadResponseBody
		}
		return parsedBody, nil
	case 401:
		return Homescript{}, ErrInvalidCredentials
	case 422:
		return Homescript{}, ErrUnprocessableEntity
	case 403:
		return Homescript{}, ErrPermissionDenied
	}
	return Homescript{}, fmt.Errorf("unknown response code: %s", res.Status)
}

// Returns a slice of Homescripts
/** Errors
- nil
- ErrNotInitialized
- ErrConnFailed
- ErrReadResponseBody
- ErrInvalidCredentials
- ErrPermissionDenied
- PrepareRequest errors
- Unknown
*/
func (c *Connection) ListHomescript() ([]Homescript, error) {
	if !c.ready {
		return nil, ErrNotInitialized
	}
	req, err := c.prepareRequest("/api/homescript/list/personal", Get, nil)
	if err != nil {
		return nil, err
	}
	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return nil, ErrConnFailed
	}
	defer res.Body.Close()
	switch res.StatusCode {
	case 200:
		resBody, err := io.ReadAll(res.Body)
		if err != nil {
			return nil, ErrReadResponseBody
		}
		var parsedBody []Homescript
		if err := json.Unmarshal(resBody, &parsedBody); err != nil {
			return nil, ErrReadResponseBody
		}
		return parsedBody, nil
	case 401:
		return nil, ErrInvalidCredentials
	case 422:
		return nil, ErrUnprocessableEntity
	case 403:
		return nil, ErrPermissionDenied
	}
	return nil, fmt.Errorf("unknown response code: %s", res.Status)
}
