package apierror

// APIError ...
// e, ok := err.(*apierror.APIError)
type APIError struct {
	HTTPStatus int    `json:"-"`
	Code       int    `json:"code,omitempty"`
	Message    string `json:"error"`
	Err        error  `json:"-"`
}

func (e *APIError) Error() string {
	return e.Message
}

// NewError ...
func NewError(status int, code int, message string, err error) *APIError {
	return &APIError{
		HTTPStatus: status,
		Code:       code,
		Message:    message,
		Err:        err,
	}
}
