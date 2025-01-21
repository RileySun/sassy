package main

import(
	"api"
	
	"log"
	"sort"
	"time"
	"strconv"
	
	"github.com/signintech/gopdf"
)

const TRANSACTION_COST = 0.001

//structs
type Report struct {
	//Internal
	API *api.API
	Users []*api.User
	
	//Metrics
	Total, Get, Add, Update, Delete int
	AvgTotal, AvgGet, AvgAdd, AvgUpdate, AvgDelete int
	TopAll, TopGet, TopAdd, TopUpdate, TopDelete *api.User
	RevTotal, RevGet, RevAdd, RevUpdate, RevDelete, RevAvg float64
	
	//PDF
	PDF gopdf.GoPdf
}

//Create
func NewReport() *Report {
	report := &Report{
		API:api.NewAPI(),
	}
	
	report.Coallate()
	report.Generate()
	
	return report
}

//Main
func (r *Report) Coallate() {
	//get User info
	var userErr error
	r.Users, userErr = r.getAllUsers()
	if userErr != nil {
		log.Println("Report->DB Error: ", userErr)
	}
	
	//Get total usage of all routes and per route
	r.Total, r.Get, r.Add, r.Update, r.Delete = r.getTotalUsage()
	
	//Get average usage of all routes and per route
	r.AvgTotal, r.AvgGet, r.AvgAdd, r.AvgUpdate, r.AvgDelete = r.getAverageUsage()
	
	//Get top account(s?)
	r.TopAll, r.TopGet, r.TopAdd, r.TopUpdate, r.TopDelete = r.getTopUsers()
	
	//Get revenue
	r.RevTotal, r.RevGet, r.RevAdd, r.RevUpdate, r.RevDelete, r.RevAvg = r.getRevenue()
}

func (r *Report) Generate() {
	createErr := r.CreatePDF()
	if createErr != nil {
		log.Println("Report->Generate->Create: ", createErr)
		return
	}
	
	writeErr := r.WritePDF()
	if writeErr != nil {
		log.Println("Report->Generate->Write: ", writeErr)
		return
	}
	
	exportErr := r.ExportPDF()
	if exportErr != nil {
		log.Println("Report->Generate->Export: ", exportErr)
		return
	}
}

//Coallate
func (r *Report) getAllUsers() ([]*api.User, error) {
	//get all usage data
	users, userErr := r.API.GetAllUsers()
	if userErr != nil {
		return nil, userErr
	}
	return users, nil
}

func (r *Report) getTotalUsage() (int, int, int, int, int) {
	//Add it up
	get, add, update, delete := 0, 0, 0, 0
	for _, u := range r.Users {
		get += u.Get
		add += u.Add
		update += u.Update
		delete += u.Delete
	}
	
	all := get + add + update + delete
	return all, get, add, update, delete
}

func (r *Report) getAverageUsage() (int, int, int, int, int) {
	userCount := len(r.Users)
	
	avgAll := r.Total/userCount
	avgGet := r.Get/userCount
	avgAdd := r.Add/userCount
	avgUpdate := r.Update/userCount
	avgDelete := r.Delete/userCount
	
	return avgAll, avgGet, avgAdd, avgUpdate, avgDelete
}

func (r *Report) getTopUsers() (*api.User, *api.User, *api.User, *api.User, *api.User) {
	//Get Top User of All Routes
	var topAll, topGet, topAdd, topUpdate, topDelete *api.User
	topInt := 0
	for _, u := range r.Users {
		total := u.Get + u.Add + u.Update + u.Delete
		if total > topInt {
			topInt = total
			topAll = u
		}
	}
	
	//Clone r.Users
	usersClone := append([]*api.User{}, r.Users...)
	
	//Get Top Get by sorting
	sort.Slice(usersClone, func(i, j int) bool {
		return usersClone[i].Get < usersClone[j].Get
	})
	topGet = usersClone[0]
	
	//Get Top Add by sorting
	sort.Slice(usersClone, func(i, j int) bool {
		return usersClone[i].Add < usersClone[j].Add
	})
	topAdd = usersClone[0]
	
	//Get Top Update by sorting
	sort.Slice(usersClone, func(i, j int) bool {
		return usersClone[i].Update < usersClone[j].Update
	})
	topUpdate = usersClone[0]
	
	//Get Top Delete by sorting
	sort.Slice(usersClone, func(i, j int) bool {
		return usersClone[i].Delete < usersClone[j].Delete
	})
	topDelete = usersClone[0]
	
	return topAll, topGet, topAdd, topUpdate, topDelete
}

func (r *Report) getRevenue() (float64, float64, float64, float64, float64, float64) {
	userCount := len(r.Users)
	
	revAll := float64(r.Total) * TRANSACTION_COST
	revGet := float64(r.Get) * TRANSACTION_COST
	revAdd := float64(r.Add) * TRANSACTION_COST
	revUpdate := float64(r.Update) * TRANSACTION_COST
	revDelete := float64(r.Delete) * TRANSACTION_COST
	revAvg := revAll/float64(userCount)
	
	return revAll, revGet, revAdd, revUpdate, revDelete, revAvg
}

//Generate
func (r *Report) CreatePDF() error {
	//Create onbject
	r.PDF = gopdf.GoPdf{}
	r.PDF.Start(gopdf.Config{ PageSize: *gopdf.PageSizeA4 })
	
	//Add First Page
	r.PDF.AddPage()
	
	//Create Font
	fontErr := r.PDF.AddTTFFont("roboto", "./RAW/Font.ttf")
	if fontErr != nil {
		return fontErr
	}
	
	//
	setFontErr := r.PDF.SetFont("roboto", "", 12)
	if setFontErr != nil {
		log.Println(setFontErr)
	}
	
	return nil
}

func (r *Report) WritePDF() error {
	//Logo
	//pdf.Image("../imgs/gopher.jpg", 200, 50, nil)
	
	//Header
	_ = r.PDF.SetFontSize(38)
	r.PDF.SetXY(220, 10)
	r.PDF.Cell(nil, "API Report")
	
	//Date
	t := time.Now()
	dateString := t.Format("2006-01-02")
	_ = r.PDF.SetFontSize(16)
	r.PDF.SetXY(260, 50)
	r.PDF.Cell(nil, dateString)
	
	//Total Usage
	_ = r.PDF.SetFontSize(22)
	r.PDF.SetXY(10, 100)
	r.PDF.Cell(nil, "- Totals")
	
	_ = r.PDF.SetFontSize(12)
	r.PDF.SetXY(30, 130)
	r.PDF.Cell(nil, "Total Usage of the API")
	
	r.TotalTable() //Print table for this section
	
	
	return nil
}

func (r *Report) ExportPDF() error {
	writeErr := r.PDF.WritePdf("Report.pdf")
	
	return writeErr
}

//Tables
func (r *Report) TotalTable() {
	tableStartY := 150.0
	marginLeft := 30.0
	table := r.PDF.NewTableLayout(marginLeft, tableStartY, 25, 1)
	
	//Columns
	table.AddColumn("TOTAL", 50, "right")
	table.AddColumn("GET", 50, "right")
	table.AddColumn("ADD", 50, "right")
	table.AddColumn("UPDATE", 50, "right")
	table.AddColumn("DELETE", 50, "right")
	
	//Rows
	t, g, a, u, d := strconv.Itoa(r.Total), strconv.Itoa(r.Get), strconv.Itoa(r.Add), strconv.Itoa(r.Update), strconv.Itoa(r.Delete)
	table.AddRow([]string{t, g, a, u, d})
	
	//Styles
	table.SetTableStyle(gopdf.CellStyle{
		BorderStyle: gopdf.BorderStyle{
			Top:    true,
			Left:   true,
			Bottom: true,
			Right:  true,
			Width:  1.0,
		},
		FillColor: gopdf.RGBColor{R: 255, G: 255, B: 255},
		TextColor: gopdf.RGBColor{R: 0, G: 0, B: 0},
		FontSize:  10,
	})
	
	table.SetHeaderStyle(gopdf.CellStyle{
		BorderStyle: gopdf.BorderStyle{
			Top:      true,
			Left:     true,
			Bottom:   true,
			Right:    true,
			Width:    2.0,
			RGBColor: gopdf.RGBColor{R: 0, G: 0, B: 0},
		},
		FillColor: gopdf.RGBColor{R: 200, G: 200, B: 200},
		TextColor: gopdf.RGBColor{R: 0, G: 0, B: 0},
		Font:      "Roboto",
		FontSize:  12,
	})
	
	table.SetCellStyle(gopdf.CellStyle{
		BorderStyle: gopdf.BorderStyle{
			Right:    true,
			Bottom:   true,
			Width:    0.5,
			RGBColor: gopdf.RGBColor{R: 0, G: 0, B: 0},
		},
		FillColor: gopdf.RGBColor{R: 255, G: 255, B: 255},
		TextColor: gopdf.RGBColor{R: 0, G: 0, B: 0},
		Font:      "Roboto",
		FontSize:  10,
	})
	
	//Draw
	table.DrawTable()
}


//Get total usage of api DONE
//Get average usage of api DONE
//Get top account usage of api DONE
//Create cost of using api DONE
//Chart revenue from using api
//Create pdf file of api usage and revenue