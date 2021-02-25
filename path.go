package cgr

import "path"

/*
Clean the path

Default is false
*/
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

