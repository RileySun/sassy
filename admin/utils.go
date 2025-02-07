package admin

import (
	"os"
	"log"
	"net/http"
)

func GetServerURL(server string) string {
	var url string
	switch server {
		case "API":
			url = os.Getenv("API_URL")
		case "Auth":
			url = os.Getenv("AUTH_URL")
		case "Admin":
			url = os.Getenv("ADMIN_URL")
	}
	return url
}

func startHTTPServer(r http.Handler) *http.Server {
	srv := &http.Server{
		Handler: r,
		Addr:    ":9090",
		// Good practice: enforce timeouts for servers you create!
		//WriteTimeout: 15 * time.Second,
		//ReadTimeout:  15 * time.Second,
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