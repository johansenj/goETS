package session_test

import (
	"fmt"
	"github.com/codegangsta/negroni"
	"github.com/gorilla/context"
	"github.com/johansenj/goETS"
	"net/http"
	"net/http/httptest"
	//"strings"
	"testing"
)

var validKey []byte = []byte("n+D+LpWrHpjzhe4HyPdALAbwrB4vk1WV")
var invalidKey []byte = []byte("secret")
var blankKey []byte = []byte("")

func Test_Session(t *testing.T) {
	var sOpt = new(session.Options)

	sOpt.MaxAge = 10000
	sOpt.CryptKey = validKey

	n := negroni.Classic()

	n.Use(session.NewSession(sOpt))

	mux := http.NewServeMux()

	mux.HandleFunc("/set", func(w http.ResponseWriter, req *http.Request) {
		context.Set(req, session.CONTEXT_KEY, "1")
	})

	mux.HandleFunc("/testSession", func(w http.ResponseWriter, req *http.Request) {
		sesStr := context.Get(req, session.CONTEXT_KEY).(string)
		if sesStr != "1" {
			t.Error("session did not comeback correctly")
		}
	})

	n.UseHandler(mux)

	res := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/set", nil)
	n.ServeHTTP(res, req)

	res2 := httptest.NewRecorder()
	req2, _ := http.NewRequest("GET", "/show", nil)
	req2.Header.Set("Session", res.Header().Get("Session"))
	n.ServeHTTP(res2, req2)
}

func Test_ClearSession(t *testing.T) {

}

func Test_ExpiredSession(t *testing.T) {

}

func Test_MissingKey(t *testing.T) {

}

func Test_KeySizeIncorrect(t *testing.T) {

}
