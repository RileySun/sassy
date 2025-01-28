package admin

import(	
	"log"
	"strings"
	"net/http"
	"html/template"
	"github.com/julienschmidt/httprouter"
)

func (a *Admin) LoadWaiting(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	/*
	authErr := a.CheckSession(r)
	if authErr != nil {
		http.Redirect(w, r, "/login", http.StatusFound)	
		return
	}
	*/

	tmpl, parseErr := template.ParseFS(HTMLFiles, "html/waiting.html")
	if parseErr != nil {
		log.Println("Template->LoadWaiting->Parse: ", parseErr)
	}
	
	//Get Action Data
	actionType := r.Header.Get("Action")
	action := r.Header.Get("Action-Type")
	
	//Template Data
	templateData := struct {
    	Action, ActionType string
	}{
		action, actionType,
	}
	
	tmpl.Execute(w, templateData)
}

func (a *Admin) LoadError(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	/*
	authErr := a.CheckSession(r)
	if authErr != nil {
		http.Redirect(w, r, "/login", http.StatusFound)	
		return
	}
	*/
	actionString := ps.ByName("action")
	splits := strings.Split(actionString, "_")
	var action, actionType string
	if len(splits) < 2 {
		log.Println("Template->LoadError->ActionString: Action Parse Error", actionString, splits)
		action, actionType = "Error", "Error"
	} else {
		action, actionType = splits[0], splits[1]
	}
	

	tmpl, parseErr := template.ParseFS(HTMLFiles, "html/waiting.html")
	if parseErr != nil {
		log.Println("Template->LoadError->Parse: ", parseErr)
	}
	
	//Template Data
	templateData := struct {
    	Action, ActionType string
	}{
		action, actionType,
	}
	
	tmpl.Execute(w, templateData)
}