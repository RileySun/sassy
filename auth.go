package main

import(
	"auth"
	
	"log"
	"context"
	"net/http"
	
	"github.com/julienschmidt/httprouter"
)

type AuthServer struct {
	Server *http.Server
	Auth *auth.Auth
	
	Status string //Allows other programs to check if still running
	
	router *httprouter.Router
}

//Create
func NewAuthServer() *AuthServer {
	server := &AuthServer{
		Auth:auth.NewAuth(),
		router:httprouter.New(),
		Status:"Init",
	}
	server.LoadRoutes()
	
	return server
}

//Management
func (s *AuthServer) LaunchServer() {
	s.Status = "Running"
	s.Server = startHTTPServer(s.router, "8080")
}

func (s *AuthServer) LoadRoutes() {
	//Auth
	s.router.POST("/token/generate", s.GenerateAccessToken)
	
}

func (s *AuthServer) Shutdown() {
	err := s.Server.Shutdown(context.Background())
	if err != nil {
		log.Println("Auth->Shutdown: ", err)
	}
}

func (s *AuthServer) Restart() {
	err := s.Server.Shutdown(context.Background()); 
	if err != nil {
		log.Println("Auth->Restart: ", err)
	} else {
		s.LaunchServer()
	}
}

func (s *AuthServer) GetStatus() string {
	return s.Status
}

//Actions
func (s *AuthServer) Action(actionType string) {
	if actionType == "Shutdown" {
		s.Status = "Shutdown"
		s.Shutdown()
	} else {
		s.Status = "Restarting"
		s.Restart()
	}
}

//Routes
func (s *AuthServer) GenerateAccessToken(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	r.ParseForm()
	refreshToken := r.PostFormValue("refresh_token")
	if refreshToken == "" {
		w.Write([]byte("Invalid Refresh Token"))
		return
	}
	
	accessToken, tokenErr := s.Auth.GenerateToken(refreshToken)
	if tokenErr != nil {
		w.Write([]byte(tokenErr.Error()))
	}
	w.Write([]byte(accessToken))
}