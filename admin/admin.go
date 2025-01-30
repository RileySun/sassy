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
	
	Status string
	
	//Server Actions
	ApiAction func(string)
	AuthAction func(string)
	
	//Server Statuses
	ApiStatus func() string
	AuthStatus func() string
	
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
	a.router.GET("/check/:server", a.CheckStatus)
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