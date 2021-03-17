package main

import (
	"fmt"
	"math"
	"net/http"
	"strconv"
	"time"

	"github.com/ChristianStefaniw/cgr"
)

func loggerMiddleware() {
	fmt.Println("Logger middleware executing...")
	t := time.Now()
	fmt.Println(t)
}

func testMiddleware() {
	fmt.Println("test middleware executing...")
}

func main() {
	r := cgr.NewRouter()
	squareConf := cgr.NewRouteConf()

	// Configuration will be passed to all routes
	r.SkipClean(true)
	r.AppendSlash(true)

	// Configuration will be passed to the route it is assigned to
	squareConf.AppendSlash(false)
	squareConf.SkipClean(false)

	logger := cgr.NewMiddleware(loggerMiddleware)
	test := cgr.NewMiddleware(testMiddleware)

	r.Route("/:msg").Method("GET").Handler(echo).Assign(logger).Assign(test).Insert()
	r.Route("/").Method("GET").Handler(home).Insert()
	r.Route("/").Method("POST").Handler(homePost).Insert()
	r.Route("/../../clean").Method("PUT").Handler(showPath).SkipClean(false).Insert()

	r.Route("/square/:num/").SetConf(squareConf).Method("GET").Handler(square).Insert()

	r.Route("/routes").Method("GET").Handler(
		func(writer http.ResponseWriter, request *http.Request) {
			for _, route := range r.ViewRouteTree() {
				writer.Write([]byte(route))
			}
		},
	).Insert()

	helloRoute := r.Route("/hello/:name/").Handler(hello).Method("GET")

	// Configure route after declaration
	helloRoute.AppendSlash(false)
	helloRoute.Insert()

	r.Run("8000")
}

func homePost(w http.ResponseWriter, _ *http.Request) {
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
	pow := fmt.Sprint(math.Pow(num, 2))
	_, err := w.Write([]byte(pow))
	if err != nil {
		panic("error")
	}
}

func showPath(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte(r.URL.Path))
}

func echo(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte(cgr.GetParams(r)["msg"]))
}
