package sdk

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// Switch Response from Smarthome
// Id is a unique identifier used for many actions regarding switches
type Switch struct {
	Id      string `json:"id"`
	Name    string `json:"name"`
	RoomId  string `json:"roomId"`
	PowerOn bool   `json:"powerOn"`
	Watts   uint16 `json:"watts"`
}

// Returns a list of switches to which the user has access to
/** Errors
- nil
- ErrInvalidCredentials
- ErrServiceUnavailable
- ErrReadResponseBody
- ErrConnFailed
- ErrNotInitialized
- PrepareRequest errors
*/
func (c *Connection) GetPersonalSwitches() (switches []Switch, err error) {
	if !c.ready {
		return nil, ErrNotInitialized
	}
	req, err := c.prepareRequest("/api/switch/list/personal", Get, nil)
	if err != nil {
		return nil, err
	}
	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return nil, ErrConnFailed
	}
	switch res.StatusCode {
	case 200:
		resBody, err := io.ReadAll(res.Body)
		if err != nil {
			return nil, ErrReadResponseBody
		}
		var parsedBody []Switch
		if err := json.Unmarshal(resBody, &parsedBody); err != nil {
			return nil, ErrReadResponseBody
		}
		return parsedBody, nil
	case 401:
		return nil, ErrInvalidCredentials
	case 503:
		return nil, ErrServiceUnavailable
	}
	return nil, fmt.Errorf("unknown response code: %s", res.Status)
}

// Returns a list containing all switches of the target instance
/** Errors
- nil
- ErrServiceUnavailable
- ErrReadResponseBody
- ErrConnFailed
- ErrNotInitialized
- PrepareRequest errors
*/
func (c *Connection) GetAllSwitches() (switches []Switch, err error) {
	if !c.ready {
		return nil, ErrNotInitialized
	}
	req, err := c.prepareRequest("/api/switch/list/all", Get, nil)
	if err != nil {
		return nil, err
	}
	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return nil, ErrConnFailed
	}
	switch res.StatusCode {
	case 200:
		resBody, err := io.ReadAll(res.Body)
		if err != nil {
			return nil, ErrReadResponseBody
		}
		var parsedBody []Switch
		if err := json.Unmarshal(resBody, &parsedBody); err != nil {
			return nil, ErrReadResponseBody
		}
		return parsedBody, nil
	case 503:
		return nil, ErrServiceUnavailable
	}
	return nil, fmt.Errorf("unknown response code: %s", res.Status)
}

// Sends a power request to Smarthome
// Only switch to which the user has permission to will work
/** Errors
- nil
- ErrNotInitialized
- ErrConnFailed
- ErrServiceUnavailable
- ErrInvalidSwitch
- ErrPermissionDenied
- PrepareRequest errors
- Unknown
*/
func (c *Connection) SetPower(switchId string, powerOn bool) error {
	if !c.ready {
		return ErrNotInitialized
	}
	req, err := c.prepareRequest("/api/power/set", Post, struct {
		Switch  string `json:"switch"`
		PowerOn bool   `json:"powerOn"`
	}{
		Switch:  switchId,
		PowerOn: powerOn,
	})
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
	case 503:
		return ErrServiceUnavailable
	case 422:
		return ErrInvalidSwitch
	case 403:
		return ErrPermissionDenied
	}
	return fmt.Errorf("unknown response code: %s", res.Status)
}
