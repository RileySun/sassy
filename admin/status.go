package admin

import(	
	"log"
	"net/http"
	"html/template"
	"github.com/julienschmidt/httprouter"
)


const SPEED_MAX = 150

func (a *Admin) LoadStatus(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	authErr := a.CheckSession(r)
	if authErr != nil {
		http.Redirect(w, r, "/login", http.StatusFound)	
		return
	}

	tmpl, parseErr := template.ParseFS(HTMLFiles, "html/status.html")
	if parseErr != nil {
		log.Println("Template->LoadStatus->Parse: ", parseErr)
	}
	
	//Images
	var apiImage, authImage, adminImage string
	
	//Get status
	var apiStatus, authStatus string
	if a.ApiStatus() != "Running" {
		apiStatus = "OFFLINE"
		apiImage = "Error"
	} else {
		apiImage = "OK"
	}
	if a.AuthStatus() != "Running" {
		authStatus = "OFFLINE"
		authImage = "Error"
	} else {
		authImage = "OK"
	}
	
	//Get Speed
	var apiSpeed, authSpeed, adminSpeed int = 56, 100, 300
	if apiSpeed > SPEED_MAX {
		apiImage = "Slow"
	}
	if authSpeed > SPEED_MAX {
		authImage = "Slow"
	}
	if adminSpeed > SPEED_MAX {
		adminImage = "Slow"
	}
	
	
	//Get Status Data
	templateData := struct {
    	Api, Auth, Admin, DB string
    	ApiTime, AuthTime, AdminTime, DBTime int
    	ApiStatus, AuthStatus string
	}{
		Api:apiImage, Auth:authImage, Admin:adminImage,
		ApiTime:apiSpeed, AuthTime:authSpeed, AdminTime:adminSpeed,
		ApiStatus:apiStatus, AuthStatus:authStatus,
	}
	
	tmpl.Execute(w, templateData)
}