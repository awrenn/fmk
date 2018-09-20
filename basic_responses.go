package fmk

import (
	"net/http"
)

var DefaultMessages = map[int][]byte{
	http.StatusOK:                  []byte("OK\n"),
	http.StatusNotFound:            []byte("Not Found\n"),
	http.StatusNotAcceptable:       []byte("Not Acceptable\n"),
	http.StatusLengthRequired:      []byte("Length Required\n"),
	http.StatusMethodNotAllowed:    []byte("Method Not Allowed\n"),
	http.StatusInternalServerError: []byte("Internal Server Error\n"),
}

type Responses struct{}

var BasicResponses Responses

func (r Responses) respond(body []byte, code int, respWriter http.ResponseWriter, req *http.Request) int {
	respWriter.WriteHeader(code)
	respWriter.Write(body)
	return code
}

func (r Responses) RespondOKWithBody(body []byte, respWriter http.ResponseWriter, req *http.Request) int {
	return r.respond(body, http.StatusOK, respWriter, req)
}

func (r Responses) RespondOK(respWriter http.ResponseWriter, req *http.Request) int {
	code := http.StatusOK
	body := DefaultMessages[code]
	return r.respond(body, code, respWriter, req)
}

func (r Responses) RespondNotFound(respWriter http.ResponseWriter, req *http.Request) int {
	code := http.StatusNotFound
	body := DefaultMessages[code]
	return r.respond(body, code, respWriter, req)
}

func (r Responses) RespondNotAcceptable(respWriter http.ResponseWriter, req *http.Request) int {
	code := http.StatusNotAcceptable
	body := DefaultMessages[code]
	return r.respond(body, code, respWriter, req)
}

func (r Responses) RespondLengthRequired(respWriter http.ResponseWriter, req *http.Request) int {
	code := http.StatusLengthRequired
	body := DefaultMessages[code]
	return r.respond(body, code, respWriter, req)
}

func (r Responses) RespondMethodNotAllowed(respWriter http.ResponseWriter, req *http.Request) int {
	code := http.StatusMethodNotAllowed
	body := DefaultMessages[code]
	return r.respond(body, code, respWriter, req)
}

func (r Responses) RespondInternalServerError(respWriter http.ResponseWriter, req *http.Request) int {
	code := http.StatusInternalServerError
	body := DefaultMessages[code]
	return r.respond(body, code, respWriter, req)
}
