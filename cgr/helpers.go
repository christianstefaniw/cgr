package cgr

import (
	"net/http"
)

func GetVar(r *http.Request, name string) string {
	ctx := r.Context()
	params := ctx.Value("params").(map[string]string)
	return params[name]
}
