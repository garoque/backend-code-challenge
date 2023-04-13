package custom_err

import "net/http"

const (
	INTERNAL_ERROR = "An unexpeted error occurred"
)

var (
	ERROR_INTERNAL    = New(http.StatusInternalServerError, "Internal server error")
	ERROR_BAD_REQUEST = New(http.StatusInternalServerError, "Invalid request")
)
