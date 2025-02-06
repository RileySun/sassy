package admin

import(	
	"log"
	"net/http"
	"html/template"
	"github.com/julienschmidt/httprouter"
)

func (a *Admin) LoadActions(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	//Check Auth
	authErr := a.CheckSession(r)
	if authErr != nil {
		a.Redirect = "actions"
		http.Redirect(w, r, "/login", http.StatusFound)	
		return
	}
	
	//Get Template
	tmpl, parseErr := template.ParseFS(HTMLFiles, "html/actions.html")
	if parseErr != nil {
		log.Println("Template->LoadActions->Parse: ", parseErr)
	}
	
	//Server Status
	var apiStatus, authStatus string
	rawAPI, _ := a.CheckServerStatus("API")
	rawAuth, _ := a.CheckServerStatus("Auth")
	if rawAPI != "Running" {
		apiStatus = rawAPI
	}
	if rawAuth != "Running" {
		authStatus = rawAuth
	}
	
	//Template Data
	templateData := struct {
    	API, Auth string
	}{
		apiStatus, authStatus,
	}
	
	tmpl.Execute(w, templateData)
}

func (a *Admin) DoAction(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	r.ParseForm()
	server := r.PostFormValue("Server")
	action := r.PostFormValue("Action")
	actionString := server + "_" + action
	
	//Check function exists before we send to right page
	//ok := true
	http.Redirect(w, r, "/waiting/" + actionString, http.StatusFound)	
	
	/*
	
	
	
	if ok {
		http.Redirect(w, r, "/waiting/" + actionString, http.StatusFound)	
	} else {
		http.Redirect(w, r, "/error/" + actionString, http.StatusFound)	
	}
	*/
}