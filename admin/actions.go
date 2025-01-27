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
	
	//Add to header for waiting page and errors
	newCookie := &http.Cookie{
		Name:"Action",
		Value:actionType + "-" + action,
		Path:"/",
		MaxAge:-1,
	}
	http.SetCookie(w, newCookie)
	
	ok := true
	var errorMessage = ""
	switch actionType {
		case "API":
			if a.ApiAction != nil {
				a.ApiAction(action)
			} else {
				ok = false
				errorMessage = "Template->DoAction->API->" + action + ": Action for API does not exist."
			}
		case "Auth":
			if a.AuthAction != nil {
				a.AuthAction(action)
			} else {
				ok = false
				errorMessage = "Template->DoAction->Auth->" + action + ": Action for Auth does not exist."
			}
		case "Admin":
			if a.AdminAction != nil {
				a.AdminAction(action)
			} else {
				ok = false
				errorMessage = "Template->DoAction->Admin->" + action + ": Action for Admin does not exist."
			}
	}
	
	if ok {
		http.Redirect(w, r, "/waiting", http.StatusFound)	
	} else {
		log.Println(errors.New(errorMessage))
		http.Redirect(w, r, "/error", http.StatusFound)	
	}
}