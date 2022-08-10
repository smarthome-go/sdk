package sdk

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type AutomationTimingMode string

const (
	TimingNormal  AutomationTimingMode = "normal"  // Will not change, automation will always execute based on this time
	TimingSunrise AutomationTimingMode = "sunrise" // Uses the local time for sunrise, each run of a set automation will update the actual time and regenerate a cron expression
	TimingSunset  AutomationTimingMode = "sunset"  // Same as above, just for sunset
)

// The automation struct, used in listing and data retrieval
type Automation struct {
	Id              uint                 `json:"id"`
	Name            string               `json:"name"`
	Description     string               `json:"description"`
	CronExpression  string               `json:"cronExpression"`
	CronDescription string               `json:"cronDescription"`
	HomescriptId    string               `json:"homescriptId"`
	Owner           string               `json:"owner"`
	Enabled         bool                 `json:"enabled"`
	TimingMode      AutomationTimingMode `json:"timingMode"`
}

// Returns a slice of automations
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
func (c *Connection) ListAutomations() ([]Automation, error) {
	if !c.ready {
		return nil, ErrNotInitialized
	}
	req, err := c.prepareRequest("/api/automation/list/personal", Get, nil)
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
		var parsedBody []Automation
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
