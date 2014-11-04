package session_test

import (
	"github.com/codegangsta/negroni"
	"github.com/gorilla/context"
	"github.com/johansenj/goETS"
	"net/http"
)

/*
	Demonstrates the general usage of this package with negroni
*/
func ExampleSession_negroniMiddleware() {
	// Setting up the session options
	var sOpt = new(session.Options)
	// Set the max age of the session in seconds
	sOpt.MaxAge = 30 * 60 // 30min * 60 sec/min
	// This is only a test key, the key needs to be secret.
	sOpt.CryptKey = []byte("n+D+LpWrHpjzhe4HyPdALAbwrB4vk1WV")

	n := negroni.Classic()

	// Using the session middleware in Negroni
	n.Use(session.NewSession(sOpt))

	mux := http.NewServeMux()

	mux.HandleFunc("/setSession", func(w http.ResponseWriter, req *http.Request) {
		// Setting the session on an individual request, if you do not modify the
		// session it will retain its settings for the request
		context.Set(req, session.CONTEXT_KEY, "1")
	})

	mux.HandleFunc("/getSession", func(w http.ResponseWriter, req *http.Request) {
		// Retrieving the session unique identifier
		_ = context.Get(req, session.CONTEXT_KEY).(string)

	})
}
