package httpserver

import (
	"fmt"
	"net/http"
)

// Error represents an HTTP Error
type Error struct {
	// Status is the http status. This should be >= 400.
	Status int `json:"-"`
	// Code is the error code that will be printed in the json response
	Code string `json:"code"`
	// Description is the error description that will be printed in the json response
	Desc string `json:"description"`
}

// Error satisfies the error interface and returns the error details in string representation
func (e Error) Error() string {
	return fmt.Sprintf("httpserver:Error: Status:[%d],Code:[%s],Desc:[%s]", e.Status, e.Code, e.Desc)
}

// ErrInternalServer is the default err for server side failures
var ErrInternalServer = &Error{
	Status: http.StatusInternalServerError,
	Code:   "INTERNAL_SERVER_ERROR",
	Desc:   "Internal Server Error",
}
