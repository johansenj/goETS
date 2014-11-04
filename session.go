/*
  goETS is a session management middleware that does not require a database call to
  check the session and implements the Encrypted Token Pattern helping prevent CSRF.
  More information about the Encypted Token Pattern can be found at:
  (https://www.owasp.org/index.php/Cross-Site_Request_Forgery_(CSRF)_Prevention_Cheat_Sheet#Encrypted_Token_Pattern).
*/
package session

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"github.com/gorilla/context"
	"io"
	//	"log"
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

/*
   NewSession is used in the creation of the Negroni middleware
*/
func NewSession(opt *Options) *Session {
	if opt == nil || opt.CryptKey == nil || len(opt.CryptKey) != KeySize {
		panic(fmt.Sprintln("Missing key or key is incorrect size"))
	}
	return &Session{opt}
}

/*
   ServeHTTP is a http server handeler for the middleware which handles the
   session data and stores the session id in the context.
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
	w.Header().Add("Session", s.packHeader(string(context.Get(req, CONTEXT_KEY).(string))))

	// context cleanup
	context.Clear(req)
}

/*
   packHeader combines the session id along with when the session is to expires.
   Returns the encrypted header
*/
func (s *Session) packHeader(sessionId string) string {
	if sessionId == "" {
		return ""
	}

	sessionDuration := time.Duration(s.config.MaxAge) * time.Second
	sessionExpire := time.Now().Add(sessionDuration)

	encodedExpire, err := sessionExpire.GobEncode()
	if err != nil {
		panic(err)
	}

	header := sessionId + ";" + string(encodedExpire)

	encryptedHeader, ok := s.encryptSessionData([]byte(header))
	str := base64.StdEncoding.EncodeToString(encryptedHeader)

	if !ok {
		return ""
	}
	return str
}

/*
   unpackHeader parses the encrypted header checks to make sure the session
   expiration has not passed.
   Returns the session id
*/
func (s *Session) unpackHeader(encryptedHeader string) string {
	var sessionExpire time.Time

	data, err := base64.StdEncoding.DecodeString(encryptedHeader)
	header, ok := s.decryptSessionData(data)
	if !ok {
		return ""
	}

	rawHeader := string(header)

	splitHeader := strings.SplitN(rawHeader, ";", 2)

	err = sessionExpire.GobDecode([]byte(splitHeader[1]))
	if err != nil {
		return ""
	}

	if time.Now().Before(sessionExpire) {
		return splitHeader[0] //session id
	} else {
		return ""
	}
}

/*
   encryptSessionData encrypts the session header.
   Returns the encrypted session header.
*/
func (s *Session) encryptSessionData(session_header []byte) ([]byte, bool) {

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

	rawCrypt := gcm.Seal(nil, iv, session_header, nil)

	return append(iv, rawCrypt...), true
}

/*
   decryptSessionData decrypts the session header.
   Returns the plaintext session header.
*/
func (s *Session) decryptSessionData(crypted_session_id []byte) ([]byte, bool) {
	//var plainText []byte

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

	cryptText := crypted_session_id[nonceSize:]
	iv := crypted_session_id[:nonceSize]

	cleartext, err := gcm.Open(nil, iv, cryptText, nil)
	if err != nil {
		panic(err)
	}
	return cleartext, true
}

/*
   randBytes gets random bytes to for use in IV
*/
func randBytes(size int) []byte {
	p := make([]byte, size)
	_, err := io.ReadFull(rand.Reader, p)
	if err != nil {
		p = nil
	}
	return p
}
