# My Golang Router

Example

```golang
package main

import (
	"./cgr"
	"math"
	"net/http"
	"strconv"
)

func main() {
	router := cgr.NewRouter()
	squareConf := cgr.NewRouteConf()

	// Configuration will be passed to all routes
	router.AppendSlash(false)

	// Configuration will be passed to the route it is assigned to
	squareConf.AppendSlash(true)
	squareConf.SkipClean(true)


	router.Route("/").Method("GET").Handler(home)
	router.Route("/square/:num").SetConf(squareConf).Method("GET").Handler(square)
	cgr.Run("8000", router)
}

func home(w http.ResponseWriter, r *http.Request){
	_, err := w.Write([]byte("home"))
	if err != nil{
		panic("error")
	}
}

func square(w http.ResponseWriter, r *http.Request){
	num, _ := strconv.ParseFloat(cgr.GetVar(r, "num"), 32)
	pow := strconv.FormatFloat(math.Pow(num, 2), 'f', -1, 32)
	_, err := w.Write([]byte(pow))
	if err != nil{
		panic("error")
	}
}
```
