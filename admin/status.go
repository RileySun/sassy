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
		a.Redirect = "status"
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
	apiStatus, apiSpeed := a.CheckServerStatus("API")
	authStatus, authSpeed := a.CheckServerStatus("Auth")
	adminStatus, adminSpeed := a.CheckServerStatus("Admin")
	
	//Text & Images
	var apiText, authText, adminText string
	if apiStatus != "Running" {
		apiText = "OFFLINE"
		apiImage = "Error"
	} else {
		apiImage = "OK"
	}
	if authStatus != "Running" {
		authText = "OFFLINE"
		authImage = "Error"
	} else {
		authImage = "OK"
	}
	if adminStatus != "Running" {
		adminText = "OFFLINE"
		adminImage = "Error"
	} else {
		adminImage = "OK"
	}
	
	//Speed
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
    	ApiStatus, AuthStatus, AdminStatus string
	}{
		Api:apiImage, Auth:authImage, Admin:adminImage,
		ApiTime:apiSpeed, AuthTime:authSpeed, AdminTime:adminSpeed,
		ApiStatus:apiText, AuthStatus:authText, AdminStatus:adminText,
	}
	
	tmpl.Execute(w, templateData)
}