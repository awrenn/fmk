package fmk

import (
	"crypto/tls"
	"fmt"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

var Server FmkWebServer

const (
    cth string = "Content-Type"
)

var (
	doubleDot []rune = []rune("..")
)

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
	ServeStatic(staticDir, pathRoot string)

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

func (ws *FmkWebServer) ServeStatic(staticDir, pathRoot string) {
	serveFiles := func(w http.ResponseWriter, req *http.Request) int {
		p, err := processURI(pathRoot, req.RequestURI)
		if err != nil {
			Log.Error.Println(err)
			return http.StatusNotFound
		}
		p = filepath.Join(staticDir, p)
		f, err := os.OpenFile(p, os.O_RDONLY, 0400)
		if err != nil {
			Log.Warning.Printf("Error attempting to open file: %s\n", err.Error())
			return http.StatusNotFound
		}
		parts := strings.Split(p, ".")
		ext := parts[len(parts)-1]
		switch ext {
		case "html":
			w.Header().Add(cth, "text/html")
		case "js":
			w.Header().Add(cth, "text/javascript")
		case "css":
			w.Header().Add(cth, "text/css")
		}
		return RespondOKWithReader(f, w, req)
	}
	ws.HandleFunc(pathRoot, serveFiles)
	//ws.serveMux.Handle(pathRoot + "/", http.StripPrefix(pathRoot, fs))
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

func processURI(webRoot, requestURI string) (fpath string, err error) {
	// processURI should take a URI, and convert it to a safe filesystem path
	// First, we need to stike the webroot out of the URI
	// Second, we should make sure that the path is safe.
	//      We are going to be overly cautious:
	//      No .. allowed in path at all - even mid word
	//      Cannot start with / (duh)
	//      Must contain only printable characters
	webRootRunes := []rune(webRoot)
	requestRunes := []rune(requestURI)
	for i := range webRootRunes {
		if webRootRunes[i] != requestRunes[i] {
			return "", fmt.Errorf("Request URI start did not match webRoot")
		}
	}
	if len(webRootRunes) == len(requestRunes) {
		return "", fmt.Errorf("No file requested")
	}

	// webRoot should end in a slash
	// so requestURI should start right start with a non-slash character right away
	requestRunes = requestRunes[len(webRootRunes):]

	if requestRunes[0] == '/' {
		return "", fmt.Errorf("Requested file started with root slash")
	}

	doubleDotCount := 0
	for _, val := range requestRunes {
		if !strconv.IsPrint(val) {
			return "", fmt.Errorf("Path contains unprintable bytes")
		}
		if val == doubleDot[doubleDotCount] {
			doubleDotCount += 1
		} else {
            doubleDotCount = 0
        }
		if doubleDotCount == len(doubleDot) {
			return "", fmt.Errorf("Path contains consecutive periods")
		}
	}
	return path.Clean(string(requestRunes)), nil
}
