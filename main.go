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
	
	//Must interconnect these before launch (working to remove these)
	adminServer.DownloadReport = apiServer.API.DownloadReport
	
	apiServer.LaunchServer()
	authServer.LaunchServer()
	adminServer.LaunchServer()
	
	for {
	
	}
}