package cgr

import (
	"net/http"
	"path"
)

func GetVars(r *http.Request) params {
	ctx := r.Context()
	p := ctx.Value("params").(map[string]string)
	return p
}

// Eliminates . and .. elements
func cleanPath(p string) string {
	if p == "" {
		return "/"
	}
	if p[0] != '/' {
		p = "/" + p
	}
	np := path.Clean(p)

	if p[len(p)-1] == '/' && np != "/" {
		np += "/"
	}

	return np
}
