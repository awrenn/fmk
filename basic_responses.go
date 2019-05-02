package fmk

import (
	"net/http"
)

var DefaultMessages = map[int][]byte{
	http.StatusOK:                  []byte("OK\n"),
	http.StatusNotFound:            []byte("Not Found\n"),
	http.StatusUnauthorized:        []byte("Unauthorized\n"),
	http.StatusNotAcceptable:       []byte("Not Acceptable\n"),
	http.StatusLengthRequired:      []byte("Length Required\n"),
	http.StatusMethodNotAllowed:    []byte("Method Not Allowed\n"),
	http.StatusInternalServerError: []byte("Internal Server Error\n"),
	http.StatusExpectationFailed:   []byte("Expectation Failed\n"),
}

func respond(body []byte, code int, respWriter http.ResponseWriter, req *http.Request) int {
	respWriter.WriteHeader(code)
	respWriter.Write(body)
	return code
}

func RespondOKWithBody(body []byte, respWriter http.ResponseWriter, req *http.Request) int {
	return respond(body, http.StatusOK, respWriter, req)
}

func RespondOK(respWriter http.ResponseWriter, req *http.Request) int {
	code := http.StatusOK
	body := DefaultMessages[code]
	return respond(body, code, respWriter, req)
}

func RespondNotFound(respWriter http.ResponseWriter, req *http.Request) int {
	code := http.StatusNotFound
	body := DefaultMessages[code]
	return respond(body, code, respWriter, req)
}

func RespondUnauthorized(respWriter http.ResponseWriter, req *http.Request) int {
	code := http.StatusUnauthorized
	body := DefaultMessages[code]
	return respond(body, code, respWriter, req)
}

func RespondNotAcceptable(respWriter http.ResponseWriter, req *http.Request) int {
	code := http.StatusNotAcceptable
	body := DefaultMessages[code]
	return respond(body, code, respWriter, req)
}

func RespondLengthRequired(respWriter http.ResponseWriter, req *http.Request) int {
	code := http.StatusLengthRequired
	body := DefaultMessages[code]
	return respond(body, code, respWriter, req)
}

func RespondMethodNotAllowed(respWriter http.ResponseWriter, req *http.Request) int {
	code := http.StatusMethodNotAllowed
	body := DefaultMessages[code]
	return respond(body, code, respWriter, req)
}

func RespondInternalServerError(respWriter http.ResponseWriter, req *http.Request) int {
	code := http.StatusInternalServerError
	body := DefaultMessages[code]
	return respond(body, code, respWriter, req)
}

func RespondExpectationFailed(respWriter http.ResponseWriter, req *http.Request) int {
	code := http.StatusExpectationFailed
	body := DefaultMessages[code]
	return respond(body, code, respWriter, req)
}
