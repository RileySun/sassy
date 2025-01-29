package admin

import(	
	"log"
	"errors"
	"net/http"
	"html/template"
	"github.com/julienschmidt/httprouter"
)

func (a *Admin) LoadActions(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	/*
	authErr := a.CheckSession(r)
	if authErr != nil {
		http.Redirect(w, r, "/login", http.StatusFound)	
		return
	}
	*/

	tmpl, parseErr := template.ParseFS(HTMLFiles, "html/actions.html")
	if parseErr != nil {
		log.Println("Template->LoadActions->Parse: ", parseErr)
	}
	
	tmpl.Execute(w, nil)
}

func (a *Admin) DoAction(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	r.ParseForm()
	actionType := r.PostFormValue("Type")
	action := r.PostFormValue("Action")
	actionString := actionType + "_" + action
	
	//Check function exists before we send to right page
	ok := true
	var errorMessage = ""
	switch actionType {
		case "API":
			if a.ApiAction == nil {
				ok = false
				errorMessage = "Template->DoAction->API->" + action + ": Action for API does not exist."
			}
		case "Auth":
			if a.AuthAction == nil {
				ok = false
				errorMessage = "Template->DoAction->Auth->" + action + ": Action for Auth does not exist."
			}
	}
	
	if ok {
		http.Redirect(w, r, "/waiting/" + actionString, http.StatusFound)	
	} else {
		log.Println(errors.New(errorMessage))
		http.Redirect(w, r, "/error/" + actionString, http.StatusFound)	
	}
}