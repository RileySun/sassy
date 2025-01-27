package admin

import(	
	"log"
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
	
	cookie, err := w.Cookie("Action")
	if err != nil {
		log.Println("Template->LoadError->Cookie:", err)
	}
	log.Println(cookie)

	tmpl, parseErr := template.ParseFS(HTMLFiles, "html/error.html")
	if parseErr != nil {
		log.Println("Template->LoadError->Parse: ", parseErr)
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