package admin

import(
	"api"
	
	"io"
	"log"
	"time"
	"embed"
	"context"
	"net/http"
	
	"github.com/julienschmidt/httprouter"
)

const API_PORT, AUTH_PORT, ADMIN_PORT = "7070", "8080", "9090"

//Create
type Admin struct {
	Server *http.Server
	API *api.API
	router *httprouter.Router
	Sessions map[string]*Session
	
	Status string
	Redirect string //if redirected to login, go elsewhere after
	
	//Server Actions
	ApiAction, AuthAction func(string)
	
	//Get PDF Report
	DownloadReport func() ([]byte, error)
}

func NewAdmin() *Admin {
	admin := &Admin{
		API:api.NewAPI(),
		router:httprouter.New(),
		Sessions: make(map[string]*Session),
		Status:"Init",
	}
	
	admin.LoadRoutes()

	return admin
}

//Manage
func (a *Admin) LaunchServer() {
	a.Status = "Running"
	a.router.ServeFiles("/static/*filepath", http.Dir("./admin/html/static"))
	a.Server = startHTTPServer(a.router)
}

func (a *Admin) LoadRoutes() {	
	//Login
	a.router.GET("/login", a.LoadLogin)
	a.router.POST("/login", a.Login)
	
	//Home
	a.router.GET("/", a.LoadHome)
	a.router.GET("/home", a.LoadHome)
	
	//Usage
	a.router.GET("/usage", a.LoadUsage)
	
	//Status
	a.router.GET("/status", a.LoadStatus)
	
	//Reports
	a.router.GET("/reports", a.LoadReports)
	
	//Actions
	a.router.GET("/actions", a.LoadActions)
	a.router.POST("/actions", a.DoAction)
	
	//Waiting
	a.router.GET("/waiting/:action", a.LoadWaiting)
	a.router.GET("/error/:action", a.LoadError)
	a.router.GET("/check/:server", a.GetServerStatus) //Waiting page tick for redirect 
	
	//Status
	a.router.GET("/checkstatus", a.CheckStatus)
}

//Manage Server
func (a *Admin) Shutdown() {
	if err := a.Server.Shutdown(context.Background()); err != nil {
		log.Println("Admin->Shutdown: ", err)
	}
}

func (a *Admin) Restart() {
	err := a.Server.Shutdown(context.Background()); 
	if err != nil {
		log.Println("Admin->Restart: ", err)
	} else {
		a.LaunchServer()
	}
}

func (s *Admin) GetStatus() string {
	return s.Status
}

//Action
func (a *Admin) Action(action string) {
	if action == "Shutdown" {
		a.Status = "Shutdown"
		a.Shutdown()
	} else {
		a.Status = "Restart"
		a.Restart()
	}
}

func (a *Admin) CheckServerStatus(server string) (string, int) {
	start := time.Now()
	
	var url = "http://localhost:"
	switch server {
		case "API":
			url += API_PORT + "/status"
		case "Auth":
			url += AUTH_PORT + "/status"
		case "Admin":
			url += ADMIN_PORT + "/checkstatus"
		default:
			return "Error - " + server, -1
	}
	
	resp, getErr := http.Get(url)
	if getErr != nil {
		//this means server offline, no need to log it
		//log.Println("Admin->CheckStatus->"+server+"->Get: ", getErr)
		return "Shutdown", -1
	}
	
	body, bodyErr := io.ReadAll(resp.Body)
	if bodyErr != nil {
		log.Println("Admin->CheckStatus->"+server+"->Body: ", bodyErr)
	}
	
	timing := time.Since(start).Milliseconds()
	if timing < 1 {
		timing = 1
	}
	
	return string(body), int(timing)
}

func (a *Admin) ServerAction(server string, action string) error {
	var url = "http://localhost:"
	switch server {
		case "API":
			url += API_PORT + "/action/" + action
		case "Auth":
			url += AUTH_PORT + "/action/" + action			
	}
	
	resp, postErr := http.Post(url, "text/plain", nil)
	if postErr != nil {
		log.Println("Admin->ServerAction->"+server+"->Get: ", postErr)
		return postErr
	}
	
	_, bodyErr := io.ReadAll(resp.Body)
	if bodyErr != nil {
		log.Println("Admin->CheckStatus->"+server+"->Body: ", bodyErr)
		return bodyErr
	}
	
	return nil
}

//Routes
func (a *Admin) CheckStatus(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	w.Write([]byte(a.GetStatus()))
}

//Embed
//go:embed html/*
var HTMLFiles embed.FS