package api

import(
	"log"
	"bytes"
	"embed"
	"io/ioutil"
	"image"
	"image/png"
	"encoding/json"
)

type Credentials struct {
	OToken string `json:"openBaoToken"` //Openbao Token
	OHost string `json:"openBaoHost"` 	//Openbao Host
	User string `json:"user"`			//Database Username
	Pass string `json:"pass"`			//Database Pass
	Host string `json:"host"`			//Database Host
	Port string `json:"port"`			//Database Port
	Database string `json:"database"`	//Database Table
}

//Actions
func LoadCredentials() *Credentials {
	jsonData := getFile("assets/config.json")
	
	var creds *Credentials
	jsonErr := json.Unmarshal(jsonData, &creds)
	if jsonErr != nil {
		log.Fatal(jsonErr)
	}
	
	return creds
}

//Files
//go:embed assets
var rawFolder embed.FS

func getFile(path string) []byte {
	imgByte, err := rawFolder.ReadFile(path)
	if err != nil {
		log.Fatalln(err)
	}
	return imgByte
}

func getDataFile(path string) []byte {
	f, _ := ioutil.ReadFile(path)
	return f
}

func imageFromBytes(byt []byte) image.Image {
	r := bytes.NewReader(byt)
	i, err := png.Decode(r)
	if err != nil {
		log.Fatal("Utils Byt2Img - " + err.Error())
	}
	return i
}