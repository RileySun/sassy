package admin

import(
	"api"
	
	"log"
	"net/http"
	"html/template"
	"github.com/julienschmidt/httprouter"
)

func (a *Admin) LoadUsage(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	/*
	authErr := a.CheckSession(r)
	if authErr != nil {
		http.Redirect(w, r, "/login", http.StatusFound)	
		return
	}
	*/

	tmpl, parseErr := template.ParseFS(HTMLFiles, "html/usage.html")
	if parseErr != nil {
		log.Println("Template->LoadUsage->Parse: ", parseErr)
	}
	
	//Get User Data
	users, dbErr := a.API.GetAllUsers() 
	if dbErr != nil {
		log.Println("Template->LoadUsage->GetAllUsers: ", dbErr)
	}
	templateData := struct {
    	Data []*api.User	
	}{
		Data:users,
	}
	
	tmpl.Execute(w, templateData)
}