package admin

import(
	"log"
	"fmt"
	"io/fs"
	"embed"
	"net/http"
	"html/template"
	"github.com/julienschmidt/httprouter"
)


//go:embed html/*
var HTMLFiles embed.FS

func (a *Admin) LoadHome(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	tmpl, parseErr := template.ParseFS(HTMLFiles, "html/index.html")
	if parseErr != nil {
		log.Println("Template->LoadHome->Parse: ", parseErr)
	}
	
	tmpl.Execute(w, nil)
}

func (a *Admin) LoadPage(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	page := ps.ByName("page")

	tmpl, parseErr := template.ParseFS(HTMLFiles, "html/" + page + ".html")
	if parseErr != nil {
		log.Println("Template->LoadPage->Parse: ", parseErr)
	}
	
	tmpl.Execute(w, nil)
}

func (a *Admin) LoadLogin(redirect string) {
	
}

func run() error {
	return fs.WalkDir(HTMLFiles, ".", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		fmt.Printf("path=%q, isDir=%v\n", path, d.IsDir())
		return nil
	})
}