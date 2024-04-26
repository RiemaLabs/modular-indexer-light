package httputl

import (
	"net/http"
	"time"
)

// Client is a simple wrapper with timeout, as core client of other services.
var Client = &http.Client{Timeout: time.Minute}
