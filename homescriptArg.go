package sdk

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
)

/* Base structs */
type HmsArgInputType string

// Datatypes which a Homescript argument can use
// Type conversion is handled by the target Homescript
// These types act as a hint for the user and
var (
	String  HmsArgInputType = "string"
	Number  HmsArgInputType = "number"
	Boolean HmsArgInputType = "boolean"
)

type HmsArgDisplay string

var (
	TypeDefault    HmsArgDisplay = "type_default"    // Uses a normal input field matching the specified data type
	StringSwitches HmsArgDisplay = "string_switches" // Shows a list of switches from which the user can select one as a string
	BooleanYesNo   HmsArgDisplay = "boolean_yes_no"  // Uses `yes` and `no` as substitutes for true and false
	BooleanOnOff   HmsArgDisplay = "boolean_on_off"  // Uses `on` and `off` as substitutes for true and false
	NumberHour     HmsArgDisplay = "number_hour"     // Displays a hour picker (0 <= h <= 24)
	NumberMinute   HmsArgDisplay = "number_minute"   // Displays a minute picker (0 <= m <= 60)
)

// Represents a Homescript with its arguments
type HomescriptWithArguments struct {
	Data      Homescript      `json:"data"`
	Arguments []HomescriptArg `json:"arguments"`
}

type HomescriptArg struct {
	Id   uint              `json:"id"`   // The Id is automatically generated
	Data HomescriptArgData `json:"data"` // The main data of the argument
}

type HomescriptArgData struct {
	ArgKey       string          `json:"argKey"`       // The unique key of the argument
	HomescriptId string          `json:"homescriptId"` // The Homescript to which the argument belongs to
	Prompt       string          `json:"prompt"`       // What the user will be prompted
	InputType    HmsArgInputType `json:"inputType"`    // Which data type is expected
	Display      HmsArgDisplay   `json:"display"`      // How the prompt will look like
}

/* Specific API responses*/
// Is returned when a new Hms argument was created successfully
type AddedHomescriptArgResponse struct {
	NewId    uint            `json:"id"`
	Response GenericResponse `json:"response"`
}

/* Functions */

// Returns a slice of Homescripts arguments which belong to a given Homescript
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
func (c *Connection) ListHomescriptArgsOfHmsId(homescriptId string) ([]HomescriptArg, error) {
	if !c.ready {
		return nil, ErrNotInitialized
	}
	req, err := c.prepareRequest(fmt.Sprintf("/api/homescript/args/of/%s", url.PathEscape(homescriptId)), Get, nil)
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
		resBody, err := ioutil.ReadAll(res.Body)
		if err != nil {
			return nil, ErrReadResponseBody
		}
		var parsedBody []HomescriptArg
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

// Returns a slice of Homescripts with their arguments
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
func (c *Connection) ListHomescriptWithArgs() ([]HomescriptWithArguments, error) {
	if !c.ready {
		return nil, ErrNotInitialized
	}
	req, err := c.prepareRequest("/api/homescript/list/personal/complete", Get, nil)
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
		resBody, err := ioutil.ReadAll(res.Body)
		if err != nil {
			return nil, ErrReadResponseBody
		}
		var parsedBody []HomescriptWithArguments
		if err := json.Unmarshal(resBody, &parsedBody); err != nil {
			fmt.Println(err.Error())
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
