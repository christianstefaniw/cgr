# My Golang Router

Instillation  
```golang
go get github.com/ChristianStefaniw/cgr
```

Example

```golang
package main

import (
	"github.com/ChristianStefaniw/cgr"
	"math"
	"net/http"
	"strconv"
)

func main() {
	r := cgr.NewRouter()
	squareConf := cgr.NewRouteConf()

	// Configuration will be passed to all routes
	r.SkipClean(true)
	r.AppendSlash(true)

	// Configuration will be passed to the route it is assigned to
	squareConf.AppendSlash(false)
	squareConf.SkipClean(false)
	
	r.Route("/:msg").Method("GET").Handler(echo).Insert(r)
	r.Route("/").Method("GET").Handler(home).Insert(r)
	r.Route("/").Method("POST").Handler(homePost).Insert(r)
	r.Route("/../../clean").Method("PUT").Handler(showPath).SkipClean(false).Insert(r)

	r.Route("/square/:num/").SetConf(squareConf).Method("GET").Handler(square).Insert(r)

	r.Route("/routes").Method("GET").Handler(
		func(writer http.ResponseWriter, request *http.Request) {
			for _, route := range r.ViewRouteTree() {
				writer.Write([]byte(route))
			}
		},
	).Insert(r)

	helloRoute := r.Route("/hello/:name/").Handler(hello).Method("GET")

	// Configure route after declaration
	helloRoute.AppendSlash(false)
	helloRoute.Insert(r)

	cgr.Run("8000", r)
}

func homePost(w http.ResponseWriter, _ *http.Request){
	w.Write([]byte("post"))
}

func home(w http.ResponseWriter, _ *http.Request) {
	_, err := w.Write([]byte("home"))
	if err != nil {
		panic("error")
	}
}

func hello(w http.ResponseWriter, r *http.Request) {
	name := cgr.GetParams(r)["name"]
	_, err := w.Write([]byte("Hello " + name))
	if err != nil {
		panic("error")
	}
}

func square(w http.ResponseWriter, r *http.Request) {
	num, _ := strconv.ParseFloat(cgr.GetParams(r)["num"], 32)
	pow := strconv.FormatFloat(math.Pow(num, 2), 'f', -1, 32)
	_, err := w.Write([]byte(pow))
	if err != nil {
		panic("error")
	}
}

func showPath(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte(r.URL.Path))
}

func echo(w http.ResponseWriter, r *http.Request){
	w.Write([]byte(cgr.GetParams(r)["msg"]))
}
```
