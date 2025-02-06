package api

import(
	"log"
	"math"
	"sort"
	"time"
	"bytes"
	"strconv"
	"image"
	"image/png"
	
	"github.com/signintech/gopdf"
	"github.com/vicanso/go-charts/v2"
)

const TRANSACTION_COST = 0.001

//structs
type Report struct {
	//Internal
	API *API
	Users []*User
	
	//Metrics
	Total, Get, Add, Update, Delete int
	AvgTotal, AvgGet, AvgAdd, AvgUpdate, AvgDelete int
	TopAll, TopGet, TopAdd, TopUpdate, TopDelete *User
	RevTotal, RevGet, RevAdd, RevUpdate, RevDelete, RevAvg float64
	
	//PDF
	PDF gopdf.GoPdf
}

//Create
func (a *API) NewReport() *Report {
	report := &Report{
		API:a,
	}
	
	report.Coallate()
	//report.Generate() User's Choice to call
	
	return report
}

func (a *API) DownloadReport() ([]byte, error) {
	report := a.NewReport()
	
	bytes, genErr := report.Download()
	if genErr != nil {
		return nil, genErr
	}
	return bytes, nil
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

func (r *Report) Download() ([]byte, error) {
	createErr := r.CreatePDF()
	if createErr != nil {
		log.Println("Report->Generate->Create: ", createErr)
		return nil, createErr
	}
	
	writeErr := r.WritePDF()
	if writeErr != nil {
		log.Println("Report->Generate->Write: ", writeErr)
		return nil, writeErr
	}
	
	bytes := r.PDF.GetBytesPdf()
	
	return bytes, nil
}

//Coallate
func (r *Report) getAllUsers() ([]*User, error) {
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

func (r *Report) getTopUsers() (*User, *User, *User, *User, *User) {
	//Get Top User of All Routes
	var topAll, topGet, topAdd, topUpdate, topDelete *User
	topInt := 0
	for _, u := range r.Users {
		total := u.Get + u.Add + u.Update + u.Delete
		if total > topInt {
			topInt = total
			topAll = u
		}
	}
	
	//Clone r.Users
	usersClone := append([]*User{}, r.Users...)
	
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
	
	rawAll := float64(r.Total) * TRANSACTION_COST
	rawGet := float64(r.Get) * TRANSACTION_COST
	rawAdd := float64(r.Add) * TRANSACTION_COST
	rawUpdate := float64(r.Update) * TRANSACTION_COST
	rawDelete := float64(r.Delete) * TRANSACTION_COST
	rawAvg := rawAll/float64(userCount)
	
	revAll := math.Floor(rawAll * 10000) / 10000
	revGet := math.Floor(rawGet * 10000) / 10000
	revAdd := math.Floor(rawAdd * 10000) / 10000
	revUpdate := math.Floor(rawUpdate * 10000) / 10000
	revDelete := math.Floor(rawDelete * 10000) / 10000
	revAvg := math.Floor(rawAvg * 10000) / 10000
	
	
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
	logo := imageFromBytes(logoBytes)
	r.PDF.ImageFromWithOption(logo, gopdf.ImageFromOption{
		Format: "png",
		X:      0,
		Y:     	0,
		Rect:   nil,
	})
	
	//Header
	_ = r.PDF.SetFontSize(38)
	r.PDF.SetXY(215, 10)
	r.PDF.Cell(nil, "API Report")
	
	//Date
	t := time.Now()
	dateString := t.Format("2006-01-02")
	_ = r.PDF.SetFontSize(16)
	r.PDF.SetXY(260, 55)
	r.PDF.Cell(nil, dateString)
	
	//Explanation
	_ = r.PDF.SetFontSize(16)
	r.PDF.SetXY(200, 80)
	r.PDF.Cell(nil, "Monthly Report of API Usage")
	
	//Draw Line
	r.PDF.SetLineWidth(2)
	r.PDF.SetLineType("dashed")
	r.PDF.Line(50, 112, 545, 112)
	r.PDF.SetLineType("solid")
	
	//Total Usage
	_ = r.PDF.SetFontSize(22)
	r.PDF.SetXY(50, 150)
	r.PDF.Cell(nil, "- Totals")
	
	_ = r.PDF.SetFontSize(12)
	r.PDF.SetXY(70, 180)
	r.PDF.Cell(nil, "Total Usage of the API")
	
	r.TotalTable() //Print table for this section
	
	
	//Avg Usage
	_ = r.PDF.SetFontSize(22)
	r.PDF.SetXY(50, 250)
	r.PDF.Cell(nil, "- Averages")
	
	_ = r.PDF.SetFontSize(12)
	r.PDF.SetXY(70, 280)
	r.PDF.Cell(nil, "Average Usage of the API")
	
	r.AverageTable() //Print table for this section
	
	//Top Usage
	_ = r.PDF.SetFontSize(22)
	r.PDF.SetXY(50, 350)
	r.PDF.Cell(nil, "- Top Users")
	
	_ = r.PDF.SetFontSize(12)
	r.PDF.SetXY(70, 380)
	r.PDF.Cell(nil, "Highest usage user")
	
	r.TopUserTable() //Print table for this section
	
	//Usage Revenue
	_ = r.PDF.SetFontSize(22)
	r.PDF.SetXY(50, 475)
	r.PDF.Cell(nil, "- Revenue")
	
	_ = r.PDF.SetFontSize(12)
	r.PDF.SetXY(70, 505)
	r.PDF.Cell(nil, "$" + strconv.FormatFloat(TRANSACTION_COST, 'f', 3, 64) + " per transaction")
	
	r.RevenueTable() //Print table for this section
	
	//Draw Chart Image
	chart := r.RevenueChart()
	img := imageFromBytes(chart)
	r.PDF.ImageFromWithOption(img, gopdf.ImageFromOption{
		Format: "png",
		X:      50,
		Y:     	550,
		Rect:   nil,
	})
	
	return nil
}

func (r *Report) ExportPDF() error {
	writeErr := r.PDF.WritePdf("Report.pdf")
	
	return writeErr
}

//Tables
func (r *Report) TotalTable() {
	tableStartY := 150.0
	marginLeft := 250.0
	table := r.PDF.NewTableLayout(marginLeft, tableStartY, 25, 1)
	
	//Columns
	table.AddColumn("GET", 60, "right")
	table.AddColumn("ADD", 60, "right")
	table.AddColumn("UPDATE", 60, "right")
	table.AddColumn("DELETE", 60, "right")
	table.AddColumn("TOTAL", 60, "right")
	
	//Rows
	t, g, a, u, d := strconv.Itoa(r.Total), strconv.Itoa(r.Get), strconv.Itoa(r.Add), strconv.Itoa(r.Update), strconv.Itoa(r.Delete)
	table.AddRow([]string{g, a, u, d, t})
	
	//Styles
	table.SetTableStyle(gopdf.CellStyle{
		BorderStyle: gopdf.BorderStyle{
			Top:    true,
			Left:   true,
			Bottom: true,
			Right:  true,
			Width:  2.0,
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

func (r *Report) AverageTable() {
	tableStartY := 250.0
	marginLeft := 250.0
	table := r.PDF.NewTableLayout(marginLeft, tableStartY, 25, 1)
	
	//Columns
	table.AddColumn("GET", 60, "right")
	table.AddColumn("ADD", 60, "right")
	table.AddColumn("UPDATE", 60, "right")
	table.AddColumn("DELETE", 60, "right")
	table.AddColumn("AVG.", 60, "right")
	
	//Rows
	v, g, a, u, d := strconv.Itoa(r.AvgTotal), strconv.Itoa(r.AvgGet), strconv.Itoa(r.AvgAdd), strconv.Itoa(r.AvgUpdate), strconv.Itoa(r.AvgDelete)
	table.AddRow([]string{g, a, u, d, v})
	
	//Styles
	table.SetTableStyle(gopdf.CellStyle{
		BorderStyle: gopdf.BorderStyle{
			Top:    true,
			Left:   true,
			Bottom: true,
			Right:  true,
			Width:  2.0,
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
			Width:    1.0,
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

func (r *Report) TopUserTable() {
	tableStartY := 350.0
	marginLeft := 250.0
	table := r.PDF.NewTableLayout(marginLeft, tableStartY, 25, 2)
	
	//Columns
	table.AddColumn("ALL", 60, "right")
	table.AddColumn("GET", 60, "right")
	table.AddColumn("ADD", 60, "right")
	table.AddColumn("UPDATE", 60, "right")
	table.AddColumn("DELETE", 60, "right")
	
	//Rows
	table.AddRow([]string{r.TopAll.Name, r.TopGet.Name, r.TopAdd.Name, r.TopUpdate.Name, r.TopDelete.Name})
	topAll := r.TopAll.Get + r.TopAll.Add + r.TopAll.Update + r.TopAll.Delete
	t, g, a, u, d := strconv.Itoa(topAll), strconv.Itoa(r.TopGet.Get), strconv.Itoa(r.TopAdd.Add), strconv.Itoa(r.TopUpdate.Update), strconv.Itoa(r.TopUpdate.Delete)
	table.AddRow([]string{t, g, a, u, d})
	
	//Styles
	table.SetTableStyle(gopdf.CellStyle{
		BorderStyle: gopdf.BorderStyle{
			Top:    true,
			Left:   true,
			Bottom: true,
			Right:  true,
			Width:  2.0,
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
			Left:	  true,
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

func (r *Report) RevenueTable() {
	tableStartY := 475.0
	marginLeft := 250.0
	table := r.PDF.NewTableLayout(marginLeft, tableStartY, 25, 1)
	
	//Columns
	table.AddColumn("GET", 50, "right")
	table.AddColumn("ADD", 50, "right")
	table.AddColumn("UPDATE", 50, "right")
	table.AddColumn("DELETE", 50, "right")
	table.AddColumn("AVG", 50, "right")
	table.AddColumn("TOTAl", 50, "right")
	
	//Rows
	t, v, g := strconv.FormatFloat(r.RevTotal, 'f', 2, 64), strconv.FormatFloat(r.RevAvg, 'f', 2, 64), strconv.FormatFloat(r.RevGet, 'f', 2, 64)
	a, u, d := strconv.FormatFloat(r.RevAdd, 'f', 2, 64), strconv.FormatFloat(r.RevUpdate, 'f', 2, 64), strconv.FormatFloat(r.RevDelete, 'f', 2, 64)
	table.AddRow([]string{g, a, u, d, v, t,})
	
	//Styles
	table.SetTableStyle(gopdf.CellStyle{
		BorderStyle: gopdf.BorderStyle{
			Top:    true,
			Left:   true,
			Bottom: true,
			Right:  true,
			Width:  2.0,
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

func (r *Report) RevenueChart() []byte {
	values := []float64{r.RevGet, r.RevAdd, r.RevUpdate, r.RevDelete,}
	p, chartErr := charts.PieRender(
		values,
		charts.TitleOptionFunc(charts.TitleOption{
			Text:    "Revenue from API Usage",
			Subtext: "",
			Left:    charts.PositionCenter,
		}),
		charts.PaddingOptionFunc(charts.Box{
			Top:    30,
			Right:  30,
			Bottom: 30,
			Left:   30,
		}),
		charts.LegendOptionFunc(charts.LegendOption{
			Padding:charts.Box{Top:50, Bottom:0, Left:0, Right:0},
			Orient: charts.OrientVertical,
			Data: []string{
				"Get",
				"Add",
				"Update",
				"Delete",
			},
			Left: charts.PositionLeft,
		}),
		charts.PieSeriesShowLabel(),
	)
	if chartErr != nil {
		log.Println("Reports->Generate->RevenueChart:" , chartErr)
	}

	buffer, bufferError := p.Bytes()
	if bufferError != nil {
		log.Println("Reports->Generate->RevenueChart:" , bufferError)
	}
	
	return buffer
}

//Util
func imageFromBytes(byt []byte) image.Image {
	r := bytes.NewReader(byt)
	i, err := png.Decode(r)
	if err != nil {
		log.Fatal("Utils Byt2Img - " + err.Error())
	}
	return i
}