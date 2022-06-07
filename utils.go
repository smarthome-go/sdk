package sdk

// A generic return value for indicating the result of a request
type GenericResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Error   string `json:"error"`
	Time    string `json:"time"`
}
