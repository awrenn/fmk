package fmk

import (
	"net/http"
	"strconv"
)

var DefaultMessages = map[int][]byte{
	http.StatusOK:                   []byte("OK\n"),
	http.StatusUnauthorized:         []byte("Unauthorized\n"),
	http.StatusNotFound:             []byte("Not Found\n"),
	http.StatusNotAcceptable:        []byte("Not Acceptable\n"),
	http.StatusLengthRequired:       []byte("Length Required\n"),
	http.StatusMethodNotAllowed:     []byte("Method Not Allowed\n"),
	http.StatusUnsupportedMediaType: []byte("Unsupported Media Type\n"),
	http.StatusInternalServerError:  []byte("Internal Server Error\n"),
}

func respond(body []byte, code int, respWriter http.ResponseWriter, req *http.Request) int {
	respWriter.Header().Set("Content-Length", strconv.FormatInt(int64(len(body)), 10))
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

func RespondUnsupportedMediaType(respWriter http.ResponseWriter, req *http.Request) int {
	code := http.StatusUnsupportedMediaType
	body := DefaultMessages[code]
	return respond(body, code, respWriter, req)
}

func RespondInternalServerError(respWriter http.ResponseWriter, req *http.Request) int {
	code := http.StatusInternalServerError
	body := DefaultMessages[code]
	return respond(body, code, respWriter, req)
}

func RespondUnauthorized(respWriter http.ResponseWriter, req *http.Request) int {
	code := http.StatusUnauthorized
	body := DefaultMessages[code]
	return respond(body, code, respWriter, req)
}
