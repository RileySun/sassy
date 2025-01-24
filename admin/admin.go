package admin

import(
	"log"
	"api"
	"net/http"
	"github.com/julienschmidt/httprouter"
)

type Admin struct {
	API *api.API
	router *httprouter.Router
	Sessions map[string]*Session
}

func NewAdmin(newAPI *api.API) *Admin {
	admin := &Admin{
		API:newAPI,
		router:httprouter.New(),
		Sessions: make(map[string]*Session),
	}
	
	admin.LoadRoutes()
	admin.LaunchServer()
	
	return admin
}

func (a *Admin) LaunchServer() {
	a.router.ServeFiles("/static/*filepath", http.Dir("./admin/html/static"))
	log.Fatal(http.ListenAndServe("localhost:9090", a.router))
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
}