package session_test

import (
	"fmt"
	"github.com/codegangsta/negroni"
	"github.com/gorilla/context"
	"github.com/johansenj/goETS"
	"net/http"
	"net/http/httptest"
	"strings"
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

	})

	n.UseHandler(mux)

}

func Test_ClearSession(t *testing.T) {

}

func Test_ExpiredSession(t *testing.T) {

}

func Test_MissingKey(t *testing.T) {

}

func Test_KeySizeIncorrect(t *testing.T) {

}
