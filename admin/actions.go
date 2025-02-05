package admin

import(	
	"log"
	"net/http"
	"html/template"
	"github.com/julienschmidt/httprouter"
)

func (a *Admin) LoadActions(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	authErr := a.CheckSession(r)
	if authErr != nil {
		a.Redirect = "actions"
		http.Redirect(w, r, "/login", http.StatusFound)	
		return
	}

	tmpl, parseErr := template.ParseFS(HTMLFiles, "html/actions.html")
	if parseErr != nil {
		log.Println("Template->LoadActions->Parse: ", parseErr)
	}
	
	tmpl.Execute(w, nil)
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