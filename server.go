package fmk

import (
	"fmt"
        "crypto/tls"
	"net/http"
)

var Server FmkWebServer

func init() {
	Server = FmkWebServer {
		tlsConf: nil,
		tlsCert: "",
		tlsKey: "",
		serveMux: http.NewServeMux(),
		defaultWrapper: DefaultWrapper,
	}
}

type WebServer interface {
        HandleFunc(path string, handleFunc func(http.ResponseWriter, *http.Request))
	Listen(address string, port int) error

	SetTLSConf(conf *tls.Config, certPath, keyPath string)

        SetDefaultWrapper(WrapperFunc)
}

type FmkWebServer struct {
	tlsConf *tls.Config
	tlsCert string
	tlsKey string

	serveMux *http.ServeMux
	defaultWrapper WrapperFunc
}

type WrapperFunc func(func(http.ResponseWriter, *http.Request)(int)) func(http.ResponseWriter, *http.Request)
func DefaultWrapper(orig func(http.ResponseWriter, *http.Request)(int)) func(http.ResponseWriter, *http.Request) {
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

func (ws *FmkWebServer) HandleFunc(path string, handleFunc func(http.ResponseWriter, *http.Request)(int)) {
	ws.serveMux.HandleFunc(path, ws.defaultWrapper(handleFunc))
}

func (ws *FmkWebServer) SetTLSConf(conf *tls.Config, certPath, keyPath string) {
	ws.tlsConf = conf
	ws.tlsCert = certPath
	ws.tlsKey  = keyPath
}

func (ws *FmkWebServer) SetDefaultWrapper(wf WrapperFunc) {
	ws.defaultWrapper = wf
}

func (ws *FmkWebServer) Listen(address string, port int) error {
        listenAddr := fmt.Sprintf("%s:%d", address, port)
	Log.Info.Printf("Listening on %s\n", listenAddr)
        server := &http.Server{
                Addr:      listenAddr,
                TLSConfig: ws.tlsConf,
		Handler: ws.serveMux,
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
