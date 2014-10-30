/*
// session is used to keep track of the users session and the session id can be
// retrieved and set via the context package using the CONTEXT_KEY
//
// the session is encrypted and thus verified via unencryption	 use of GCM cipher
// used for authenticity check
*/
package session

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"fmt"
	"github.com/gorilla/context"
	"io"
	"net/http"
	"strings"
	"time"
)

const CONTEXT_KEY string = "session_id"

// KeySize is size of AES-256-GCM keys in bytes.
const KeySize = 32

const nonceSize = 24

type Options struct {
	MaxAge   int64
	CryptKey []byte
}

type Session struct {
	config *Options
}

// Middleware is a struct that has a ServeHTTP method
func NewSession(opt *Options) *Session {
	if opt == nil || opt.CryptKey == nil || len(opt.CryptKey) != KeySize {
		panic(fmt.Sprintln("Missing key or key is incorrect size"))
	}
	return &Session{opt}
}

/*
// ServeHTTP returns a http server handeler for the middleware which handles the
// session data and stores the session id in the context.
// Returns the middleware handler after session setup.
*/
func (s *Session) ServeHTTP(w http.ResponseWriter, req *http.Request, next http.HandlerFunc) {
	//parse session string from header
	rawSession := req.Header.Get("Session")

	if rawSession == "" {
		context.Set(req, CONTEXT_KEY, "")
	} else {
		context.Set(req, CONTEXT_KEY, s.unpackHeader(rawSession))
	}

	// Call the next middleware handler
	next(w, req)

	// set session string to header
	w.Header().Add("Session", s.packHeader(string(context.Get(req, "session_id").(string))))

	// context cleanup
	context.Clear(req)
}

/*
// packHeader combines the session id along with when the session is to expires.
// Returns the encrypted header
*/
func (s *Session) packHeader(sessionId string) string {
	sessionDuration := time.Duration(s.config.MaxAge) * time.Second
	sessionExpire := time.Now().Add(sessionDuration)

	encodedExpire, err := sessionExpire.GobEncode()
	if err != nil {
		panic(err)
	}

	header := sessionId + ";" + string(encodedExpire)

	encryptedHeader, ok := s.encryptSessionData([]byte(header))

	if !ok {
		return ""
	}
	return string(encryptedHeader)
}

/*
// unpackHeader parses the encrypted header checks to make sure the session
// expiration has not passed.
// Returns the session id
*/
func (s *Session) unpackHeader(encryptedHeader string) string {
	var sessionExpire time.Time

	header, ok := s.decryptSessionData([]byte(encryptedHeader))
	if !ok {
		return ""
	}

	splitHeader := strings.SplitN(string(header), ";", 1)

	sessionId := splitHeader[0]

	err := sessionExpire.GobDecode([]byte(splitHeader[1]))
	if err != nil {
		return ""
	}

	if time.Now().Before(sessionExpire) {
		return sessionId
	} else {
		return ""
	}
}

/*
// encryptSessionData encrypts the session header.
// Returns the encrypted session header.
*/
func (s *Session) encryptSessionData(session_header []byte) ([]byte, bool) {

	var rawCrypt []byte

	c, err := aes.NewCipher(s.config.CryptKey)
	if err != nil {
		return nil, false
	}

	gcm, err := cipher.NewGCM(c)
	if err != nil {
		return nil, false
	}

	iv := randBytes(gcm.NonceSize())
	if iv == nil {
		return nil, false
	}

	gcm.Seal(rawCrypt, iv, session_header, nil)

	return append(iv, rawCrypt...), true
}

/*
// decryptSessionData decrypts the session header.
// Returns the plaintext session header.
*/
func (s *Session) decryptSessionData(crypted_session_id []byte) ([]byte, bool) {
	var plainText []byte

	c, err := aes.NewCipher(s.config.CryptKey)
	if err != nil {
		return nil, false
	}

	gcm, err := cipher.NewGCM(c)
	if err != nil {
		return nil, false
	}

	nonceSize := gcm.NonceSize()
	if len(crypted_session_id) < nonceSize {
		return nil, false
	}

	gcm.Open(plainText, crypted_session_id[:nonceSize], crypted_session_id[nonceSize:], nil)
	return plainText, true
}

/*
// randBytes gets random bytes to for use in IV
*/
func randBytes(size int) []byte {
	p := make([]byte, size)
	_, err := io.ReadFull(rand.Reader, p)
	if err != nil {
		p = nil
	}
	return p
}
