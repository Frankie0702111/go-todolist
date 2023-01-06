package responses

import "strings"

const (
	// 4xx
	NoTokenFound                           = 400001
	BearerTokenNotInProperFormat           = 400002
	TokenInvalid                           = 400003
	EmailAlreadyExists                     = 400004
	TokenDoesNotExistOrExpired             = 401001
	InvalidCredential                      = 401002
	TokenContainsAnInvalidNumberOfSegments = 401003
	FailedToLogout                         = 401004

	// 5xx
	SignatureFailed = 500001
)

var (
	messages = map[int]string{
		// 4xx
		400001: "No token found.",
		400002: "Bearer token not in proper format.",
		400003: "Token invalid.",
		400004: "Email already exists.",
		401001: "Token does not exist or expired.",
		401002: "Invalid credential.",
		401003: "Token contains an invalid number of segments.",
		401004: "Failed to logout.",

		// 5xx
		500001: "Signature failed.",
	}
)

// Create a new struct for the response data
type Response struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Errors  interface{} `json:"errors"`
	Data    interface{} `json:"data"`
}

// EmptyObj object is used when data doesnt want to be null on json
// type EmptyObject struct {
// }

// SuccessResponse returns a success response with the given data
func SuccessResponse(code int, message string, data interface{}) Response {
	return Response{
		Code:    code,
		Message: message,
		Errors:  nil,
		Data:    data,
	}
}

// ErrorResponse returns an error response with the given data
func ErrorsResponse(code int, message string, err string, data interface{}) Response {
	splittedError := strings.Split(err, "\n")
	return Response{
		Code:    code,
		Message: message,
		Errors:  splittedError,
		Data:    data,
	}
}

// ErrorResponse Returns an error response with the message given by the code
func ErrorsResponseByCode(code int, message string, errCode int, data interface{}) Response {
	return Response{
		Code:    code,
		Message: message,
		Errors:  messages[errCode],
		Data:    data,
	}
}
