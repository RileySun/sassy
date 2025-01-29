package admin

import(
	"api"
	
	"log"
	//"time"
	"context"
	"net/http"
	
	"github.com/julienschmidt/httprouter"
)

type Admin struct {
	Server *http.Server
	API *api.API
	router *httprouter.Router
	Sessions map[string]*Session
	
	ActionServer string
	Status string
	
	ApiAction func(string)
	AuthAction func(string)
}

func NewAdmin() *Admin {
	admin := &Admin{
		API:api.NewAPI(),
		router:httprouter.New(),
		Sessions: make(map[string]*Session),
		ActionServer:"", Status:"None",
	}
	
	admin.LoadRoutes()

	return admin
}

func (a *Admin) LaunchServer() {
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
	a.router.GET("/check", a.CheckStatus)
}

//Manage Server
func (a *Admin) Shutdown() {
	log.Println("Server is shutting down...")

	if err := a.Server.Shutdown(context.Background()); err != nil {
		log.Println("Admin->Shutdown: ", err)
	}
}

func (a *Admin) Restart() {
	a.Shutdown()
	a.LaunchServer()
}