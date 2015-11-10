package routes

import "net/http"

type ErrorJSON struct {
	Message string `json:"message"`
	Code    int    `json:"code"`
}

func NewErrorJSON(message string, code int) *ErrorJSON {
	return &ErrorJSON{Message: message, Code: code}
}

func NewErrorJSONExists() *ErrorJSON {
	return &ErrorJSON{Message: "resource already exists ", Code: http.StatusConflict}
}

func NewErrorJSONBadRequest() *ErrorJSON {
	return &ErrorJSON{Message: "could not parse request data ", Code: http.StatusBadRequest}
}

func NewErrorJSONUnexpectedError(message string) *ErrorJSON {
	return NewErrorJSON("unexpected error occured ", http.StatusInternalServerError)
}

func NewErrorJSONNotFound() *ErrorJSON {
	return NewErrorJSON("not found ", http.StatusNotFound)
}
