package api

import (
	"log"
	"encoding/json"
)

//Structs
type User struct {
	ID int
	Name string
	Trial bool
	Get, Add, Update, Delete int //Limits
}

type Model struct {
	ID int `json:"id"`
	Name string `json:"name"`
	Desc string `json:"desc"`
}

type Image struct {
	ID int `json:"id"`
	Model_ID int `json:"model_id"`
	Path string `json:"path"`
	Desc string `json:"desc"`
}

type Video struct {
	ID int `json:"id"`
	Model_ID int `json:"model_id"`
	Path string `json:"path"`
	Desc string `json:"desc"`
}

//Json
func (m *Model) JSON() []byte {
	json, err := json.MarshalIndent(m, "", "	")
	if err != nil {
		log.Fatal("Model JSON - ", err)
	}
	
	return json
}

func (i *Image) JSON() []byte {
	json, err := json.MarshalIndent(i, "", "	")
	if err != nil {
		log.Fatal("Image JSON - ", err)
	}
	
	return json
}

func (v *Video) JSON() []byte {
	json, err := json.MarshalIndent(v, "", "	")
	if err != nil {
		log.Fatal("Video JSON - ", err)
	}
	
	return json
}