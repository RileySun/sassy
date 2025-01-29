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
	actionString := ps.ByName("action")
	splits := strings.Split(actionString, "_")
	var action, actionType string
	if len(splits) < 2 {
		log.Println("Template->LoadError->ActionString: Action Parse Error", actionString, splits)
		action, actionType = "Error", "Error"
		http.Redirect(w, r, "/", http.StatusFound)	
	} else {
		action, actionType = splits[0], splits[1]
	}
	
	//Template Data
	templateData := struct {
    	Action, ActionType string
	}{
		action, actionType,
	}
	
	tmpl.Execute(w, templateData)
	
	a.ApiAction("Shutdown")
	
	switch action {
		case "Api":
			if actionType == "Shutdown" {
				a.ApiAction("Shutdown")
				a.Status = "Down"
			} else {
				a.ApiAction("Restart")
				a.Status = "Restarting"
			}	
		case "Auth":
			if actionType == "Shutdown" {
				a.AuthAction("Shutdown")
				a.Status = "Down"
			} else {
				a.AuthAction("Restart")
				a.Status = "Restarting"
			}
		case "Admin":
			if actionType == "Shutdown" {
				a.Shutdown()
				a.Status = "Down"
			} else {
				a.Restart()
				a.Status = "Restarting"
			}
	}
}

func (a *Admin) LoadError(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	authErr := a.CheckSession(r)
	if authErr != nil {
		http.Redirect(w, r, "/login", http.StatusFound)	
		return
	}
	
	actionString := ps.ByName("action")
	splits := strings.Split(actionString, "_")
	var action, actionType string
	if len(splits) < 2 {
		log.Println("Template->LoadError->ActionString: Action Parse Error", actionString, splits)
		action, actionType = "Error", "Error"
		http.Redirect(w, r, "/", http.StatusFound)	
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

func (a *Admin) CheckStatus(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	_, writeErr := w.Write([]byte(a.Status))
	if writeErr != nil {
		log.Println(writeErr)
	}
}