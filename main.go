package main

import(
	"api"
	"auth"
	"admin"
)

func main() {
	apiServer := api.NewApiServer()
	authServer := auth.NewAuthServer()
	adminServer := admin.NewAdmin()
	
	apiServer.LaunchServer()
	authServer.LaunchServer()
	adminServer.LaunchServer()
	
	for {
	
	}
}