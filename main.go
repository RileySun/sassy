package main

import(
	"fmt"

	"api"
)

func main() {
	auth := api.NewAuth()
	
	keys := api.LoadKeys()
	fmt.Println(auth.IsValidKey(keys["RILEY"]))
	auth.GenerateToken()
}
