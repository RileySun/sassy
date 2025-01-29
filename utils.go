package main

import (
	"log"
	"time"
	"net/http"
)

func startHTTPServer(r http.Handler, port string) *http.Server {
	srv := &http.Server{
		Handler: r,
		Addr: ":" + port,
		WriteTimeout: 15 * time.Second,
		ReadTimeout: 15 * time.Second,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil {
			// cannot panic, because this probably is an intentional close
			log.Printf("Httpserver: ListenAndServe() error: %s", err)
		} else {
			log.Printf("Httpserver: ListenAndServe() closing...")
		}
	}()

	// returning reference so caller can call Shutdown()
	return srv
}