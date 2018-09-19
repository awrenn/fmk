package fmk

import (
	"net/http"
)

const (
	FRONT_SLASH = 47 // literally ascii/utf8 for '/'
)

var (
	OKBody                  = []byte("OK\n")
	NotFoundBody            = []byte("Not Found\n")
	NotAcceptableBody       = []byte("Not Acceptable\n")
	LengthRequiredBody      = []byte("Length Required\n")
	MethodNotAllowedBody    = []byte("Method Not Allowed\n")
	InternalServerErrorBody = []byte("Internal Server Error\n")
)

func RespondOKWithBody(body []byte, respWriter http.ResponseWriter, req *http.Request) int {
	respWriter.WriteHeader(http.StatusOK)
	respWriter.Write(body)
	return http.StatusOK
}

func RespondNotFound(respWriter http.ResponseWriter, req *http.Request) int {
	respWriter.WriteHeader(http.StatusNotFound)
	respWriter.Write(NotFoundBody)
	return http.StatusNotFound
}

func RespondNotAcceptable(respWriter http.ResponseWriter, req *http.Request) int {
	respWriter.WriteHeader(http.StatusNotAcceptable)
	respWriter.Write(NotAcceptableBody)
	return http.StatusNotAcceptable
}

func RespondLengthRequired(respWriter http.ResponseWriter, req *http.Request) int {
	respWriter.WriteHeader(http.StatusLengthRequired)
	respWriter.Write(LengthRequiredBody)
	return http.StatusNotFound
}

func RespondOK(respWriter http.ResponseWriter, req *http.Request) int {
	respWriter.WriteHeader(http.StatusOK)
	respWriter.Write(OKBody)
	return http.StatusOK
}

func RespondMethodNotAllowed(respWriter http.ResponseWriter, req *http.Request) int {
	respWriter.WriteHeader(http.StatusMethodNotAllowed)
	respWriter.Write(MethodNotAllowedBody)
	return http.StatusMethodNotAllowed
}

func RespondInternalServerError(respWriter http.ResponseWriter, req *http.Request) int {
	respWriter.WriteHeader(http.StatusInternalServerError)
	respWriter.Write(InternalServerErrorBody)
	return http.StatusInternalServerError
}

