package api

import(
	"os"
	"log"
	"embed"
	
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
		log.Fatal("Error loading .env file")
	}
	
	creds := &Credentials{
		User:os.Getenv("API_DB_USER"),
		Pass:os.Getenv("API_DB_PASS"),
		Host:os.Getenv("API_DB_HOST"),
		Port:os.Getenv("API_DB_PORT"),
		Database:os.Getenv("API_DB_DATABASE"),
	}
	
	return creds
}

//Files
//go:embed assets
var AssetsFolder embed.FS

func getFile(path string) ([]byte, error) {
	fileByte, err := AssetsFolder.ReadFile(path)
	if err != nil {
		return nil, err
	}
	return fileByte, nil
}