package session_test

import (
	"fmt"
	"github.com/codegangsta/negroni"
	"github.com/johansenj/goETS"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

const validKey []byte = []byte("")
const invalidKey []byte = []byte("secret")
const blankKey []byte = []byte("")

func Test_Session(t *testing.T) {

}

func Test_ClearSession(t *testing.T) {

}

func Test_ExpiredSession(t *testing.T) {

}

func Test_MissingKey(t *testing.T) {

}

func Test_KeySizeIncorrect(t *testing.T) {

}
