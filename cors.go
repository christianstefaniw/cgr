package cgr

import (
	"net/http"
)

func CorsHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(200)
}
