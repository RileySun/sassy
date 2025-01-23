package admin

import(
	"log"
	"embed"
	"net/http"
	"html/template"
	"github.com/julienschmidt/httprouter"
)


//go:embed html/*
var HTMLFiles embed.FS

func (a *Admin) LoadLogin(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	tmpl, parseErr := template.ParseFS(HTMLFiles, "html/login.html")
	if parseErr != nil {
		log.Println("Template->LoadPage->Parse: ", parseErr)
	}
	
	//Show Popup
	tmpl.Execute(w, false)
}

func (a *Admin) LoadHome(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	authErr := a.CheckSession(r)
	if authErr != nil {
		http.Redirect(w, r, "/login", http.StatusFound)	
		return
	}

	tmpl, parseErr := template.ParseFS(HTMLFiles, "html/home.html")
	if parseErr != nil {
		log.Println("Template->LoadHome->Parse: ", parseErr)
	}
	
	tmpl.Execute(w, nil)
}