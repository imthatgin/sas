package sas

import (
	"errors"
	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
	"net/http"
)

// TODO: Use custom error types rather than errors.New
var (
	ErrorResourceNoAccess    = errors.New("access denied to specified resource")
	ErrorResourceNotFound    = errors.New("resource requested does not exist")
	ErrorResourceInvalidID   = errors.New("id specified is not valid")
	ErrorResourceInvalidData = errors.New("invalid data for the requested operation")

	ErrorDatabaseIssue        = errors.New("database issue")
	ErrorFatalSetupNoBindType = errors.New("no bind type has been set for this operation")
)

func ManagedModelErrorHandler(err error, c echo.Context) {
	code := http.StatusInternalServerError
	message := "unknown"

	switch {
	case errors.Is(err, ErrorResourceNoAccess):
		code = http.StatusForbidden
		message = ErrorResourceNoAccess.Error()

	case errors.Is(err, ErrorResourceNotFound):
		code = http.StatusNotFound
		message = ErrorResourceNotFound.Error()

	case errors.Is(err, ErrorResourceInvalidID):
		code = http.StatusBadRequest
		message = ErrorResourceInvalidID.Error()

	case errors.Is(err, ErrorResourceInvalidData):
		code = http.StatusBadRequest
		message = ErrorResourceInvalidData.Error()

	case errors.Is(err, ErrorDatabaseIssue):
		code = http.StatusInternalServerError
		message = ErrorDatabaseIssue.Error()

	default:
		code = http.StatusInternalServerError
		message = "unknown"
		log.Errorf("Unhandled error type: %s", err)
	}

	log.Warn("Handled error: ", err)
	_ = c.String(code, message)
}
