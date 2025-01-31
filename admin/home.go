package admin

import(
	"log"
	"embed"
	"net/http"
	"html/template"
	"github.com/julienschmidt/httprouter"
)

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