package cgwf

import (
	"fmt"
	"log"
	"net/http"
)

func Run(router *Router, port string){
	for _, warning := range router.warnings{
		fmt.Println(warning)
	}

	fmt.Println("CGWF is listing on port:", port)
	log.Fatal(http.ListenAndServe(":"+port, router))
}
