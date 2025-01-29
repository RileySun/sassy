package main

import(
	"admin"
)



type AdminServer struct {
	Admin *admin.Admin
	
	ApiAction,  AuthAction func(string)
}

func NewAdminServer() *AdminServer {
	server := &AdminServer{
		Admin:admin.NewAdmin(),
	}
	
	return server
}

//Launch
func (s *AdminServer) LaunchServer() {
	s.Admin.LaunchServer()
}