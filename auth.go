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
	
	router *httprouter.Router
}

func NewAuthServer() *AuthServer {
	server := &AuthServer{
		Auth:auth.NewAuth(),
		router:httprouter.New(),
	}
	server.LoadRoutes()
	
	return server
}

//Launch
func (s *AuthServer) LaunchServer() {
	s.Server = startHTTPServer(s.router, "8080")
}

func (s *AuthServer) LoadRoutes() {
	//Auth
	s.router.POST("/token/generate", s.GenerateAccessToken)
	
}

func (s *AuthServer) Shutdown() {
	log.Println("API Server is shutting down...")

	if err := s.Server.Shutdown(context.Background()); err != nil {
		log.Println("Auth->Shutdown: ", err)
	}
}

func (s *AuthServer) Action(actionType string) {
	if actionType == "Shutdown" {
		s.Shutdown()
	} else {
		s.Shutdown()
		s.LaunchServer()
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