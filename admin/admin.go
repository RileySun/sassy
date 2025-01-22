package admin

import(
	"log"
	"net/http"
	"github.com/julienschmidt/httprouter"
)

type Admin struct {
	router *httprouter.Router
}

func NewAdmin() *Admin {
	admin := &Admin{
		router:httprouter.New(),
	}
	
	admin.LoadRoutes()
	admin.LaunchServer()
	
	return admin
}

func (a *Admin) LaunchServer() {
	//a.router.ServeFiles()
	log.Fatal(http.ListenAndServe("localhost:9090", a.router))
}

func (a *Admin) LoadRoutes() {
	//Auth
	a.router.GET("/", a.LoadHome)
	a.router.GET("/:page", a.LoadPage)
}