package admin

import(	
	"log"
	"strings"
	"net/http"
	"html/template"
	"github.com/julienschmidt/httprouter"
)

func (a *Admin) LoadWaiting(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	authErr := a.CheckSession(r)
	if authErr != nil {
		a.Redirect = "actions"
		http.Redirect(w, r, "/login", http.StatusFound)	
		return
	}

	tmpl, parseErr := template.ParseFS(HTMLFiles, "html/waiting.html")
	if parseErr != nil {
		log.Println("Template->LoadWaiting->Parse: ", parseErr)
	}
	
	//Get Action Data
	actionString := ps.ByName("action")
	splits := strings.Split(actionString, "_")
	var server, action string
	if len(splits) < 2 {
		log.Println("Template->LoadWaiting->ActionString: Action Parse Error", actionString, splits)
		server, action = "Error", "Error"
		http.Redirect(w, r, "/", http.StatusFound)	
	} else {
		server, action = splits[0], splits[1]
	}
	
	//Template Data
	templateData := struct {
    	Server, Action string
	}{
		server, action,
	}
	
	tmpl.Execute(w, templateData)
	
	if server != "Admin" {
		actionErr := a.ServerAction(server, action)
		if actionErr != nil {
			log.Println("Admin->LoadWaiting->ServerAction: ", actionErr)
		}
	} else {
		if action == "Shutdown" {
			a.Status = "Down"
		} else {
			a.Status = "Restarting"
		}
	}
}

func (a *Admin) LoadError(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	authErr := a.CheckSession(r)
	if authErr != nil {
		a.Redirect = "actions"
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

func (a *Admin) GetServerStatus(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	server := ps.ByName("server")
	
	status, _ := a.CheckServerStatus(server)

	_, writeErr := w.Write([]byte(status))
	if writeErr != nil {
		log.Println(writeErr)
	}
}