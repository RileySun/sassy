package auth

import(
	"log"
	"time"
	"errors"
	"strings"
	"context"
	"net/http"
	
	"github.com/julienschmidt/httprouter"
)

type AuthServer struct {
	Server *http.Server
	Auth *Auth
	
	Status string
	
	router *httprouter.Router
}

//Create
func NewAuthServer() *AuthServer {
	server := &AuthServer{
		Auth:NewAuth(),
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
	
	//Status
	s.router.GET("/status", s.CheckStatus)
	
	//Action
	s.router.POST("/action/:type", s.Action)
}

func (s *AuthServer) Shutdown() {
	shutdownCtx, shutdownRelease := context.WithTimeout(context.Background(), 1*time.Second)
	defer shutdownRelease()
	err := s.Server.Shutdown(shutdownCtx)
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
func (s *AuthServer) Action(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	_, writeErr := w.Write([]byte("OK"))
	if writeErr != nil {
		log.Println(writeErr)
	}
	
	action := ps.ByName("type")
	if action == "Shutdown" {
		s.Status = "Shutdown"
		s.Shutdown()
	} else {
		s.Status = "Restart"
		s.Restart()
	}
}

//Routes
//Authentication (How do I decouple this without causing another request check...)
func (s *AuthServer) CheckAuthentication(r *http.Request) (int, error) {
	Oauth := r.Header["Authorization"]
	
	//dont split before u check it exists
	if len(Oauth) == 0 {
		return 0, errors.New("Invalid Authorization Header")
	}
	split :=  strings.Split(Oauth[0], " ")
	
	//Is Auth there?
	if len(split) == 1 {
		return 0, errors.New("Invalid Authorization Header")
	}
	accessToken := split[1]
	
	//Check Auth
	authErr := s.Auth.CheckToken(accessToken)
	if authErr != nil {
		return 0, authErr
	}
	
	//GetUserID
	userID := s.Auth.GetUserIdFromToken(accessToken)
	if userID == 0 {
		return 0, errors.New("Authorization Failed, Please Contact Administrator")
	}
	
	return userID, nil
} //checks authorization and returns userID

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

func (s *AuthServer) CheckStatus(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	w.Write([]byte(s.GetStatus()))
}