# GO Encrypted Token Session [![GoDoc](http://godoc.org/github.com/johansenj/goETS?status.png)](http://godoc.org/github.com/johansenj/goETS)

goETS is an implementation of the Encrypted Token Paturn Session written for 
Negroni middleware.

## Usage

### Options
	- MaxAge   - is the max length of time in seconds that a session token is valid
	- CryptKey - is the secret 256 bit AES key 

## Examples

### General Example
~~~ go
    // Setting up the session options
	var sOpt = new(session.Options)
	
	// Set the max age of the session in seconds
	sOpt.MaxAge = 30 * 60 // 30min * 60 sec/min
	
	// This is only a test key, the key needs to be secret.
	sOpt.CryptKey = []byte("n+D+LpWrHpjzhe4HyPdALAbwrB4vk1WV")

	n := negroni.Classic()

	// Using the session middleware in Negroni
	n.Use(session.NewSession(sOpt))
~~~

### Setting Session
~~~ go
	context.Set(req, session.CONTEXT_KEY, "1")
~~~

### Clearing Session
~~~ go
	context.Set(req, session.CONTEXT_KEY, "")
~~~

### Retrieving session
~~~ go
	sesStr := context.Get(req, session.CONTEXT_KEY).(string)
~~~