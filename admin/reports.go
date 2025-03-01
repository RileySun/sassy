package admin

import(
	"io"
	"log"
	"api"
	"net/http"
	"html/template"
	"encoding/base64"
	
	"github.com/julienschmidt/httprouter"
)

func (a *Admin) LoadReports(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	authErr := a.CheckSession(r)
	if authErr != nil {
		a.Redirect = "reports"
		http.Redirect(w, r, "/login", http.StatusFound)	
		return
	}


	tmpl, parseErr := template.ParseFS(HTMLFiles, "html/reports.html")
	if parseErr != nil {
		log.Println("Template->LoadReports->Parse: ", parseErr)
	}
	
	//Get Reports Data
	report := a.API.NewReport()
	/*
	reportBytes, reportErr := report.Download()
	if reportErr != nil {
		log.Println("Template->LoadReports->Download")
	}
	*/
	reportBytes := a.DownloadReport()
	
	//Template Data
	templateData := struct {
    	Total, Get, Add, Update, Delete int
		AvgTotal, AvgGet, AvgAdd, AvgUpdate, AvgDelete int
		TopAll, TopGet, TopAdd, TopUpdate, TopDelete *api.User
		RevTotal, RevGet, RevAdd, RevUpdate, RevDelete float64
		Chart string
		Report []byte
	}{
		report.Total, report.Get, report.Add, report.Update, report.Delete,
		report.AvgTotal, report.AvgGet, report.AvgAdd, report.AvgUpdate, report.AvgDelete,
		report.TopAll, report.TopGet, report.TopAdd, report.TopUpdate, report.TopDelete,
		report.RevTotal, report.RevGet, report.RevAdd, report.RevUpdate, report.RevDelete,
		base64.StdEncoding.EncodeToString(report.RevenueChart()),
		reportBytes,
	}
	
	tmpl.Execute(w, templateData)
}

func (a *Admin) DownloadReport() ([]byte) {
	resp, getErr := http.Get("http://localhost:"+API_PORT+"/report/download")
	if getErr != nil {
		log.Println("Admin->DownloadReport->POST:", getErr)
	}
	
	body, bodyErr := io.ReadAll(resp.Body)
	if bodyErr != nil {
		log.Println("Admin->DownloadReport->->Body:", bodyErr)
	}
	
	return body
}