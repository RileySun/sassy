package auth

import(
	"log"
	"embed"
	"encoding/json"
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
	jsonData, fileErr := getFile("assets/config.json")
	if fileErr != nil {
		log.Fatal(fileErr)
	}
	
	var creds *Credentials
	jsonErr := json.Unmarshal(jsonData, &creds)
	if jsonErr != nil {
		log.Fatal(jsonErr)
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