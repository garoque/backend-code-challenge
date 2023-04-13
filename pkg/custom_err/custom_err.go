package custom_err

import "fmt"

type Error struct {
	HTTPCode int    `mapstructure:"code" json:"http_code,omitempty"`
	Message  string `mapstructure:"message" json:"message,omitempty"`
}

func (e *Error) Error() string {
	return fmt.Sprintf("code: %v - message: %v", e.HTTPCode, e.Message)
}

func New(httpCode int, message string) error {
	return &Error{
		HTTPCode: httpCode,
		Message:  message,
	}
}
