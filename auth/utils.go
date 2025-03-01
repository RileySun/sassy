package auth

import(
	"os"
	"log"
	"time"
	"net/http"
	
	"github.com/joho/godotenv"
)

type Credentials struct {
	//OToken string `json:"openBaoToken"`//Openbao Token
	//OHost string `json:"openBaoHost"` //Openbao Host
	User string `json:"user"`			//Database Username
	Pass string `json:"pass"`			//Database Pass
	Host string `json:"host"`			//Database Host
	Port string `json:"port"`			//Database Port
	Database string `json:"database"`	//Database Table
}

//Actions
func LoadCredentials() *Credentials {
	envErr := godotenv.Load()
	if envErr != nil {
		log.Println("Auth: Error loading .env file - ", envErr)
		log.Println("This may be caused by running in docker")
	}
	
	creds := &Credentials{
		User:os.Getenv("AUTH_DB_USER"),
		Pass:os.Getenv("AUTH_DB_PASS"),
		Host:os.Getenv("AUTH_DB_HOST"),
		Port:os.Getenv("AUTH_DB_PORT"),
		Database:os.Getenv("AUTH_DB_DATABASE"),
	}
	
	return creds
}

func GetServerURL(server string) string {
	envErr := godotenv.Load()
	if envErr != nil {
		log.Println("Auth: Error loading .env file - ", envErr)
		log.Println("This may be caused by running in docker")
	}
	
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
			//log.Printf("Httpserver: ListenAndServe() error: %s", err)
		} else {
			log.Printf("Httpserver: ListenAndServe() closing...")
		}
	}()

	// returning reference so caller can call Shutdown()
	return srv
}