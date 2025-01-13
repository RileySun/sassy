package main

import(
	"log"
	"fmt"
	"net/http"
	"github.com/julienschmidt/httprouter"
)


//Main
func LaunchServer() {
    router := httprouter.New()
    router.GET("/", Index)
    router.GET("/hello/:name", Hello)

    log.Fatal(http.ListenAndServe("localhost:8080", router))
}

//Authentication


//Routes
func Index(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
    fmt.Fprint(w, "Welcome!\n")
}


func Hello(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
    fmt.Fprintf(w, "hello, %s!\n", ps.ByName("name"))
}