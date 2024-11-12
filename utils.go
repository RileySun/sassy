package main

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
	User string `json:"user"`
	Pass string `json:"pass"`
	Host string `json:"host"`
	Port string `json:"port"`
	Database string `json:"database"`
}

//Actions
func LoadCredentials() *Credentials {
	jsonData := getFile("RAW/config.json")
	
	var creds *Credentials
	jsonErr := json.Unmarshal(jsonData, &creds)
	if jsonErr != nil {
		log.Fatal(jsonErr)
	}
	
	return creds
}

//Files
//go:embed RAW
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