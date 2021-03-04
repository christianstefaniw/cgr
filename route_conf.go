package cgr


// route configurations
type routeConf struct {
	appendSlash bool
	skipClean   bool
}


/*
	example.com/path is treated the same as example.com/path/

	Default is true
*/
func (conf *routeConf) AppendSlash(value bool) *routeConf {
	conf.appendSlash = value
	return conf
}

/*
	example.com/path is treated the same as example.com/path/

	Default is true
*/
func (route *route) AppendSlash(value bool) *route {
	route.appendSlash = value
	return route
}



/*
If set to false:
 1. Replace multiple slashes with a single slash.
 2. Eliminate each . path name element (the current directory).
 3. Eliminate each inner .. path name element (the parent directory)
    along with the non-.. element that precedes it.
 4. Eliminate .. elements that begin a rooted path:
    that is, replace "/.." by "/" at the beginning of a path

Default value is false
*/
func (route *route) SkipClean(value bool) *route {
	route.skipClean = value
	return route
}
/*
If set to false:
 1. Replace multiple slashes with a single slash.
 2. Eliminate each . path name element (the current directory).
 3. Eliminate each inner .. path name element (the parent directory)
    along with the non-.. element that precedes it.
 4. Eliminate .. elements that begin a rooted path:
    that is, replace "/.." by "/" at the beginning of a path

Default value is false
*/
func (conf *routeConf) SkipClean(value bool) *routeConf {
	conf.skipClean = value
	return conf
}


// Set custom configurations for a route
func (route *route) SetConf(conf *routeConf) *route {
	route.routeConf = *conf
	return route
}


func (conf *routeConf) setDefaultRouteConf() {
	conf.appendSlash = true
	conf.skipClean = false
}