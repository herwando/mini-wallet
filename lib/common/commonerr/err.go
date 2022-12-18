package commonerr

import (
	"encoding/json"
	"net/http"

	"github.com/go-playground/validator"
)

//
// An ErrorMessage represents an error message format and list of error.
//
type ErrorMessage struct {
	ErrorList []*ErrorFormat `json:"error_list"`
	Code      int            `json:"code"`
}

//
// An ErrorFormat represents an error message format and code that we used.
//
type ErrorFormat struct {
	ErrorName        string `json:"error_name"`
	ErrorDescription string `json:"error_description"`
}

// OrderErrorFormat Error Format for order
type OrderErrorFormat struct {
	Code   string `json:"code"`
	Status string `json:"status"`
	Title  string `json:"title"`
	Detail string `json:"detail"`
}

//
// Common internal server error message
//
const (
	InternalServerName        = "internal_server_error"
	InternalServerDescription = "The server is unable to complete your request"
)

//
// Common error unautorized
//
const (
	UnauthorizedErrorName        = "access_denied"
	UnauthorizedErrorDescription = "Authorization failed by filter."
)

//
// DefaultBadRequest is common bad request error message
//
var DefaultBadRequest = ErrorFormat{
	ErrorName:        "bad_request",
	ErrorDescription: "Your request resulted in error",
}

//
// DefaultInputBody return bad request for bad body request
//
var DefaultInputBody = ErrorFormat{
	ErrorName:        "bad_request",
	ErrorDescription: "Your body request resulted in error",
}

//
// Default Not Found error message
//
const (
	NotFound            = "not_found"
	NotFoundDescription = "Page not found"
)

//
// NewErrorMessage will Create new error message
//
func NewErrorMessage() *ErrorMessage {
	return &ErrorMessage{}
}

//
// SetBadRequest will Set bad request
//
func (errorMessage *ErrorMessage) SetBadRequest() *ErrorMessage {
	errorMessage.Code = http.StatusBadRequest
	return errorMessage
}

//
// SetErrorValidator containts setter error from github.com/go-playground/validator
//
func (errorMessage *ErrorMessage) SetErrorValidator(err error) *ErrorMessage {
	if err != nil {

		// this check is only needed when your code could produce
		// an invalid value for validation such as interface with nil
		// value most including myself do not usually have code like this.
		if _, ok := err.(*validator.InvalidValidationError); ok {
			return errorMessage
		}

		for _, err := range err.(validator.ValidationErrors) {
			errorMessage.Append(err.Field(), err.Tag())
		}

	}
	return errorMessage
}

//
// SetNewError is function return new error message.
// It support to set code, error name, and error description
//
func SetNewError(code int, errorName, errDesc string) *ErrorMessage {
	return &ErrorMessage{
		Code: code,
		ErrorList: []*ErrorFormat{
			&ErrorFormat{
				ErrorName:        errorName,
				ErrorDescription: errDesc,
			},
		},
	}
}

//
// SetDefaultNewNotFound returns a default 404 error page not found
//
func SetDefaultNewNotFound() *ErrorMessage {
	return SetNewError(404, NotFound, NotFoundDescription)
}

//
// SetNewBadRequest is function return new error message with bad request standard code(400).
// It support to set error name and error description
//
func SetNewBadRequest(errorName, errDesc string) *ErrorMessage {
	return SetNewError(http.StatusBadRequest, errorName, errDesc)
}

//
// SetNewBadRequestByFormat is function return new error message with bad request standard code(400).
// It support to set error name and error description using error format
//
func SetNewBadRequestByFormat(ef *ErrorFormat) *ErrorMessage {
	return &ErrorMessage{
		Code: http.StatusBadRequest,
		ErrorList: []*ErrorFormat{
			ef,
		},
	}
}

//
// SetListErr is function return new error message with error code.
// It support to set error name and error description using error format
//
func SetListErr(code int, listErr []*ErrorFormat) *ErrorMessage {
	return &ErrorMessage{
		Code:      code,
		ErrorList: listErr,
	}
}

//
// SetDefaultNewBadRequest returns default bad request error with http code 400
//
func SetDefaultNewBadRequest() *ErrorMessage {
	return SetNewBadRequestByFormat(&DefaultBadRequest)
}

//
// SetDefaultErrBodyRequest return body request error
//
func SetDefaultErrBodyRequest() *ErrorMessage {
	return SetNewBadRequestByFormat(&DefaultInputBody)
}

//
// SetNewInternalError is function return new error message with internal server error standard code(500).
//
func SetNewInternalError() *ErrorMessage {
	return SetNewError(http.StatusInternalServerError, InternalServerName, InternalServerDescription)
}

//
// SetNewUnauthorizedError is function return new error message with unauthorized error code(401).
// It support to set error name and error description
//
func SetNewUnauthorizedError(errorName, errDesc string) *ErrorMessage {
	return SetNewError(http.StatusUnauthorized, errorName, errDesc)
}

//
// SetDefaultUnauthorized is function return new error message with unauthorized error code(401).
// It support to set error name and error description
//
func SetDefaultUnauthorized() *ErrorMessage {
	return SetNewUnauthorizedError(UnauthorizedErrorName, UnauthorizedErrorDescription)
}

//
// Append is function add error to existing error message.
// It support to set error name and error description.
//
func (errorMessage *ErrorMessage) Append(errorName, errDesc string) *ErrorMessage {
	errorMessage.ErrorList = append(errorMessage.ErrorList, &ErrorFormat{
		ErrorName:        errorName,
		ErrorDescription: errDesc,
	})
	return errorMessage
}

//
// AppendFormat is function add error to existing error message.
// It support to set error name and error description using error format
//
func (errorMessage *ErrorMessage) AppendFormat(ef *ErrorFormat) *ErrorMessage {
	errorMessage.ErrorList = append(errorMessage.ErrorList, ef)
	return errorMessage
}

//
// GetListError is function to get list error message.
//
func (errorMessage *ErrorMessage) GetListError() []*ErrorFormat {
	return errorMessage.ErrorList
}

//
// GetCode is function to get code.
//
func (errorMessage *ErrorMessage) GetCode() int {
	return errorMessage.Code
}

//
// Marshal will get error byte
//
func (errorMessage *ErrorMessage) Marshal() []byte {
	b, _ := json.Marshal(errorMessage)
	return b
}

//
// ToString will get string
//
func (errorMessage *ErrorMessage) ToString() string {
	return string(errorMessage.Marshal())
}

// Error to implement error interface
func (errorMessage *ErrorMessage) Error() string {
	return errorMessage.ToString()
}

//
// GetErrorDesc Index1 get error index 1
//
func (errorMessage *ErrorMessage) GetErrorDesc() string {
	if len(errorMessage.ErrorList) == 0 {
		return ""
	}
	return errorMessage.ErrorList[0].ErrorDescription
}
