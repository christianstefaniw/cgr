package cgr

// route configurations
type RouteConf struct {
	appendSlash     bool
	skipClean       bool
	handlePreflight bool
}

// example.com/path is treated the same as example.com/path/
// Default is true
func (route *Route) AppendSlash(value bool) *Route {
	route.appendSlash = value
	return route
}

// Whether or not to handle cors preflight request
func (route *Route) HandlePreflight(value bool) *Route {
	route.handlePreflight = value
	return route
}

//If set to false:
// 1. Replace multiple slashes with a single slash.
// 2. Eliminate each . path name element (the current directory).
// 3. Eliminate each inner .. path name element (the parent directory)
//    along with the non-.. element that precedes it.
// 4. Eliminate .. elements that begin a rooted path:
//    that is, replace "/.." by "/" at the beginning of a path
//
//Default value is false
func (route *Route) SkipClean(value bool) *Route {
	route.skipClean = value
	return route
}

//If set to false:
// 1. Replace multiple slashes with a single slash.
// 2. Eliminate each . path name element (the current directory).
// 3. Eliminate each inner .. path name element (the parent directory)
//    along with the non-.. element that precedes it.
// 4. Eliminate .. elements that begin a rooted path:
//    that is, replace "/.." by "/" at the beginning of a path
//
//Default value is false

func (conf *RouteConf) SkipClean(value bool) *RouteConf {
	conf.skipClean = value
	return conf
}

// example.com/path is treated the same as example.com/path/
// Default is true

func (conf *RouteConf) AppendSlash(value bool) *RouteConf {
	conf.appendSlash = value
	return conf
}

// Whether or not to handle cors preflight request
func (conf *RouteConf) HandlePreflight(value bool) *RouteConf {
	conf.handlePreflight = value
	return conf
}

// Set custom configurations for a route
func (route *Route) SetConf(conf *RouteConf) *Route {
	route.RouteConf = *conf
	return route
}

func (conf *RouteConf) setDefaultRouteConf() {
	conf.appendSlash = true
	conf.skipClean = false
	conf.handlePreflight = false
}
