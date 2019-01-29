package fmk

import (
	"crypto/tls"
	"fmt"
	"net/http"
	"time"
)

var Server FmkWebServer

func init() {
	Server = FmkWebServer{
		tlsConf:        nil,
		tlsCert:        "",
		tlsKey:         "",
		serveMux:       http.NewServeMux(),
		defaultWrapper: DefaultWrapper,
	}
}

var (
	DEFAULT_READ_TIMEOUT  = time.Duration(5) * time.Second
	DEFAULT_WRITE_TIMEOUT = time.Duration(30) * time.Second
	DEFAULT_IDLE_TIMEOUT  = time.Duration(240) * time.Second
)

type WebServer interface {
	HandleFunc(path string, handleFunc func(http.ResponseWriter, *http.Request))
	Listen(address string, port int) error

	SetReadTimeout(timeout *time.Duration)
	SetWriteTimeout(timeout *time.Duration)
	SetIdleTimeout(timeout *time.Duration)

	SetTLSConf(conf *tls.Config)
	SetTLSServing(certFile, keyFile string)

	SetDefaultWrapper(WrapperFunc)
}

type FmkWebServer struct {
	tlsConf *tls.Config
	tlsCert string
	tlsKey  string

	ReadTimeout  *time.Duration
	WriteTimeout *time.Duration
	IdleTimeout  *time.Duration

	serveMux       *http.ServeMux
	defaultWrapper WrapperFunc
}

type WrapperFunc func(func(http.ResponseWriter, *http.Request) int) func(http.ResponseWriter, *http.Request)

func DefaultWrapper(orig func(http.ResponseWriter, *http.Request) int) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, req *http.Request) {
		logstr := "%s On %s From %s: %d"
		code := orig(w, req)
		Log.Info.Println(
			fmt.Sprintf(logstr,
				req.Method,
				req.RequestURI,
				req.RemoteAddr,
				code,
			),
		)
	}
}

func (ws *FmkWebServer) HandleFunc(path string, handleFunc func(http.ResponseWriter, *http.Request) int) {
	ws.serveMux.HandleFunc(path, ws.defaultWrapper(handleFunc))
}

func (ws *FmkWebServer) SetTLSConf(conf *tls.Config) {
	ws.tlsConf = conf
}

func (ws *FmkWebServer) SetTLSServing(certPath, keyPath string) {
	ws.tlsCert = certPath
	ws.tlsKey = keyPath
}

func (ws *FmkWebServer) SetDefaultWrapper(wf WrapperFunc) {
	ws.defaultWrapper = wf
}

func (ws *FmkWebServer) FixTimeouts() {
	if ws.ReadTimeout == nil {
		ws.ReadTimeout = &DEFAULT_READ_TIMEOUT
	}
	if ws.WriteTimeout == nil {
		ws.WriteTimeout = &DEFAULT_WRITE_TIMEOUT
	}
	if ws.IdleTimeout == nil {
		ws.IdleTimeout = &DEFAULT_IDLE_TIMEOUT
	}
}

func (ws *FmkWebServer) Listen(address string, port int) error {
	listenAddr := fmt.Sprintf("%s:%d", address, port)
	Log.Info.Printf("Listening on %s\n", listenAddr)
	ws.FixTimeouts()
	server := &http.Server{
		Addr:         listenAddr,
		TLSConfig:    ws.tlsConf,
		Handler:      ws.serveMux,
		ReadTimeout:  *ws.ReadTimeout,
		WriteTimeout: *ws.WriteTimeout,
		IdleTimeout:  *ws.IdleTimeout,
	}
	var err error
	if ws.tlsCert == "" {
		err = server.ListenAndServe()
	} else {
		err = server.ListenAndServeTLS(ws.tlsCert, ws.tlsKey)
	}
	Log.Error.Println(err)
	return err
}
