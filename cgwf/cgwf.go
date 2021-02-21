package cgwf

import (
	"fmt"
	"log"
	"net/http"
)

func Run(router *Router, port string){
	fmt.Println("CGWF is listing on port:", port)
	log.Fatal(http.ListenAndServe(":"+port, router))
}
