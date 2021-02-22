package main

import (
	"./cgwf"
	"math"
	"net/http"
	"strconv"
)

func main() {
	router := cgwf.NewRouter()
	router.AppendSlash = false
	router.Route("GET", "/", home)
	router.Route("GET", "^/panic", error)
	router.Route("GET", "/square/:num/", square)
	cgwf.Run(router, "8000")
}

func home(w http.ResponseWriter, r *http.Request){
	w.Write([]byte("home"))
}

func square(w http.ResponseWriter, r *http.Request){
	num, _ := strconv.ParseFloat(cgwf.GetParam(r, "num"), 32)
	pow := strconv.FormatFloat(math.Pow(num, 2), 'f', -1, 32)
	w.Write([]byte(pow))
}

func error(w http.ResponseWriter, r *http.Request){
	panic("something bad happened!")
}