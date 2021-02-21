package cgwf

import (
	"net/http"
)

func GetParam(r *http.Request, name string) string {
	ctx := r.Context()
	params := ctx.Value("params").(map[string]string)
	return params[name]
}