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

	apiServer.LaunchServer()
	authServer.LaunchServer()
	adminServer.LaunchServer()
	
	for true {
		
    }
}