package admin

import(	
	"log"
	"net/http"
	"html/template"
	"github.com/julienschmidt/httprouter"
)

func (a *Admin) LoadStatus(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	authErr := a.CheckSession(r)
	if authErr != nil {
		http.Redirect(w, r, "/login", http.StatusFound)	
		return
	}

	tmpl, parseErr := template.ParseFS(HTMLFiles, "html/status.html")
	if parseErr != nil {
		log.Println("Template->LoadStatus->Parse: ", parseErr)
	}
	
	//Get Status Data
	templateData := struct {
    	Api, Auth, Admin, DB string
    	ApiTime, AuthTime, AdminTime, DBTime int
	}{
		Api:"OK", Auth:"Error", Admin:"Slow", DB:"OK",
		ApiTime:56, AuthTime:0, AdminTime:300, DBTime:20,
	}
	
	tmpl.Execute(w, templateData)
}