package main

import(
	
)

func main() {
	apiServer := NewApiServer()
	authServer := NewAuthServer()
	adminServer := NewAdminServer()
	
	//Must interconnect these before launch
	adminServer.Admin.ApiAction = apiServer.Action
	adminServer.Admin.AuthAction = authServer.Action
	adminServer.Admin.ApiStatus = apiServer.GetStatus
	adminServer.Admin.AuthStatus = authServer.GetStatus
	adminServer.Admin.DownloadReport = apiServer.API.DownloadReport

	apiServer.LaunchServer()
	authServer.LaunchServer()
	adminServer.LaunchServer()
	
	for true {
		
    }
}