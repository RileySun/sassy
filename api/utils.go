package api

import(
	"os"

	"log"
	"bytes"
	"embed"
	"io/ioutil"
	"image"
	"image/png"
	"encoding/json"
)


type Credentials struct {
	User string `json:"user"`
	Pass string `json:"pass"`
	Host string `json:"host"`
	Port string `json:"port"`
	Database string `json:"database"`
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

func LoadKeys() map[string]string {
	jsonData := getFile("assets/keys.json")
	
	var keys map[string]string
	jsonErr := json.Unmarshal(jsonData, &keys)
	if jsonErr != nil {
		log.Fatal(jsonErr)
	}
	
	return keys
}

//TEMPORARY DONT JUMP DOWN MY THROAT
func saveKeys(newKeys map[string]string) {
	json, err := json.MarshalIndent(newKeys, "", "	")
	if err != nil {
		log.Fatal("saveData - ", err)
	}
	path := "./keys.json"
	os.WriteFile(path, json, 0755);
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