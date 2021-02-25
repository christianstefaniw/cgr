package cgr

const (
	pathDelimiter  = '/'
	paramDelimiter = ':'
)

func (p *params) paramsToMap() map[string]string {
	paramsAsMap := make(map[string]string)
	for _, k := range *p {
		paramsAsMap[k.key] = k.value
	}
	return paramsAsMap
}

func appendSlash(path string) string{
	return path + "/"
}