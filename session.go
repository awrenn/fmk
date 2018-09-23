package fmk

import (
	"crypto/rand"
	"errors"
	"net/http"
	"sync"
	"time"
)

const (
	SIDBaseLen int = 32
	CipherKeyLength int = 32
)

var (
	ALPHA      []byte = []byte("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
	DefaultSessionName string = "FMKSession"
	DefaultSessionKey []byte  = []byte("PleaseChangeMe;MoreRandom=Better")
	SIDInUse   error  = errors.New("SID In Use")
	SIDMissing error  = errors.New("SID Missing")

	SIDBaseString string = "Kitty"
	BadSID        error  = errors.New("Invalid SID")

	LenMismatch error = errors.New("Mismatched Length")
	KeyTooShort error = errors.New("KeyTooShort")
	SessionKeyUnset error = errors.New("Session Key Unset")
)

var Sessions SessionManager
func init() {
	Sessions = NewSessionManager(DefaultSessionName, time.Duration(100), time.Duration(100))
	Sessions.SetSessionKey(DefaultSessionKey)
}

type SessionManager interface {
	sessionInit(sid string) (Session, error)
	sessionGet(sid string) (Session, error)
	sessionDestroy(sid string) error
	sessionClean(sessionMaxLife time.Duration)
	GetSession(respWriter http.ResponseWriter, req *http.Request) (Session, error)

	SetSessionKey(sessionKey []byte)
	SetSessionName(sessionName string)
	SetSessionDomain(domainName string)
	SetSessionPath(pathName string)
	SetGCRate(gcRate time.Duration)
	SetSessionMaxLife(ml time.Duration)
}

type FmkSessionManager struct {
	SessionName    string

	SessionKey     []byte
	DomainName     string
	Path           string

	Lock           sync.Mutex
	SessionMaxLife time.Duration
	GCRate         time.Duration
	Book map[string]Session
}

func NewSessionManager(sessionName string, sml, gcrate time.Duration) *FmkSessionManager {
	sm := &FmkSessionManager{
		Book: make(map[string]Session),
		SessionKey:     make([]byte, 32),
		SessionName:    sessionName,
		SessionMaxLife: sml,
		GCRate:         gcrate,
	}
	go sm.GC()
	return sm
}

func (sm *FmkSessionManager) SetSessionKey(sessionKey []byte) {
	copy(sm.SessionKey, sessionKey)
}

func (sm *FmkSessionManager) SetSessionName(sessionName string) {
	sm.SessionName = sessionName
}

func (sm *FmkSessionManager) SetSessionDomain(domainName string) {
	sm.DomainName = domainName
}

func (sm *FmkSessionManager) SetSessionPath(path string) {
	sm.Path = path
}

func (sm *FmkSessionManager) SetGCRate(gcRate time.Duration) {
	sm.GCRate = gcRate
}

func (sm *FmkSessionManager) SetSessionMaxLife(ml time.Duration) {
	sm.SessionMaxLife = ml
}

func (sm *FmkSessionManager) sessionInit(sid string) (Session, error) {
	domain := sm.DomainName
	path := sm.Path
	sessionName := sm.SessionName

	newSess := NewSession(sid, domain, path, sessionName)
	_, used := sm.Book[sid]
	if used {
		return nil, SIDInUse
	}
	sm.Book[sid] = newSess
	return newSess, nil
}

func (sm *FmkSessionManager) GetSession(writer http.ResponseWriter, req *http.Request) (Session, error) {
	if string(sm.SessionKey) == string(DefaultSessionKey) {
		return nil, SessionKeyUnset
	}
        cookie, err := req.Cookie(sm.SessionName)
        sm.Lock.Lock()
        defer sm.Lock.Unlock()
        var sess Session

        // Cookie is either unset or broken
        if err != nil {
                for {
                        sid := sm.generateSID()
                        sess, err = sm.sessionInit(sid)
                        if err != SIDInUse {
                                break
                        }
                }
                http.SetCookie(writer, sess.GetCookie())
                return sess, nil
        }
        sid := cookie.Value
        err = sm.Validate(sid)
        if err != nil {
                return nil, err
        }
        sess, err = sm.sessionGet(sid)
        if err != nil {
                sess, _ = sm.sessionInit(sid)
        }
        return sess, nil
}


func (sm *FmkSessionManager) generateSID() string {
	base := make([]byte, SIDBaseLen)
	_, err := rand.Read(base)
	if err != nil {
		Log.Error.Println(err)
		return ""
	}
	n := copy(base[len(base)-len(SIDBaseString):], SIDBaseString)
	if n == len(base) {
		Log.Warning.Println("Basestring is longer than total SID Length")
	}
	if err != nil {
		Log.Error.Println(err)
		return ""
	}
	sid := sm.Sign(string(base))
	// sid := sm.Sign(base+"+"+req.RemoteAddr)
	return string(sid)
}

func (sm *FmkSessionManager) sessionGet(sid string) (Session, error) {
	sess, ok := sm.Book[sid]
	if !ok {
		return nil, SIDMissing
	}
	sess.Update()
	return sess, nil
}

func (sm *FmkSessionManager) sessionDestroy(sid string) error {
	delete(sm.Book, sid)
	return nil
}

func (sm *FmkSessionManager) sessionClean(sessionMaxLife time.Duration) {
	for sid, sess := range sm.Book {
		lu := sess.GetLastUpdate()
		if time.Since(lu) > sessionMaxLife {
			sm.sessionDestroy(sid)
		}

	}
}

func (sm *FmkSessionManager) Sign(base string) string {
	sid, err := encrypt([]byte(base), []byte(sm.SessionKey))
	if err != nil {
		return ""
	}
	return string(armor(sid))
}

func (sm *FmkSessionManager) Validate(sid string) error {
	dearmored, err := dearmor(sid)
	if err != nil {
		return err
	}
	base, err := decrypt(dearmored, []byte(sm.SessionKey))
	if err != nil {
		return err
	}
	if !equal(base[len(base)-len(SIDBaseString):], []byte(SIDBaseString)) {
		return BadSID
	}
	return nil

}

func (sm *FmkSessionManager) GC() {
	for {
		sm.Lock.Lock()
		sm.sessionClean(sm.SessionMaxLife)
		sm.Lock.Unlock()
		time.Sleep(sm.GCRate)
	}
}

type Session interface {
	Set(key, value string) error
	Get(key string) string
	Del(key string) error
	SessionID() string
	GetLastUpdate() time.Time
	Update()
	GetCookie() *http.Cookie
}

type FmkSession struct {
	Book        map[string]string
	SID         string
	DomainName  string
	Path        string
	SessionName string
	LastUpdate  time.Time
}

func NewSession(sid, domain, path, sessionName string) Session {
	newSess := &FmkSession{
		SID:         sid,
		Book:        make(map[string]string),
		LastUpdate:  time.Now(),
		DomainName:      domain,
		Path:        path,
		SessionName: sessionName,
	}
	return newSess
}

func (sess *FmkSession) Set(key, value string) error {
	sess.Book[key] = value
	return nil
}

func (sess *FmkSession) Get(key string) string {
	return sess.Book[key]
}

func (sess *FmkSession) Del(key string) error {
	delete(sess.Book, key)
	return nil
}

func (sess *FmkSession) SessionID() string {
	return sess.SID
}

func (sess *FmkSession) GetLastUpdate() time.Time {
	return sess.LastUpdate
}

func (sess *FmkSession) Update() {
	sess.LastUpdate = time.Now()
	return
}

func (sess *FmkSession) GetCookie() *http.Cookie {
	cookie := &http.Cookie{
		Name:   sess.SessionName,
		Value:  sess.SID,
		Path:   sess.Path,
		Domain: sess.DomainName,
		Secure: false,
	}
	return cookie
}
