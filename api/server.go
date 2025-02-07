package api

import(
	"io"
	"log"
	"fmt"
	"time"
	"errors"
	"context"
	"strconv"
	"net/http"
	
	"github.com/julienschmidt/httprouter"
)

type ApiServer struct {
	Server *http.Server
	API *API
	
	Routes *Routes
	
	Status string
	
	router *httprouter.Router
}

//Create
func NewApiServer() *ApiServer {
	server := &ApiServer{
		API:NewAPI(),
		
		router:httprouter.New(),
		Status:"Init",
	}
	server.Routes = server.API.NewRoutes()
	server.LoadRoutes()
	
	return server
}

//Manage
func (s *ApiServer) LaunchServer() {
	s.Status = "Running"
	s.Server = startHTTPServer(s.router, "7070")
}

func (s *ApiServer) LoadRoutes() {	
	//Models
	s.router.GET("/model/get/:modelID", s.GetModel)
	s.router.POST("/model/add", s.AddModel)
	s.router.POST("/model/update", s.UpdateModel)
	s.router.GET("/model/delete/:modelID", s.DeleteModel)
	
	//Images
	s.router.GET("/images/get/:modelID", s.GetImages)
	s.router.POST("/images/add", s.AddImage)
	s.router.POST("/images/update", s.UpdateImage)
	s.router.GET("/images/delete/:imageID", s.DeleteImage)
	
	//Videos
	s.router.GET("/videos/get/:modelID", s.GetVideos)
	s.router.POST("/videos/add", s.AddVideo)
	s.router.POST("/videos/update", s.UpdateVideo)
	s.router.GET("/videos/delete/:videoID", s.DeleteVideo)
	
	//Status
	s.router.GET("/status", s.CheckStatus)
	
	//Action
	s.router.POST("/action/:type", s.Action)
	
	//Error
	s.router.GET("/", Error404)
	
	//Download Report
	s.router.GET("/report/download", s.DownloadReport)
}

func (s *ApiServer) Shutdown() {
	shutdownCtx, shutdownRelease := context.WithTimeout(context.Background(), 1*time.Second)
	defer shutdownRelease()
	err := s.Server.Shutdown(shutdownCtx)
	if err != nil {
		log.Println("API->Shutdown: ", err)
	}
}

func (s *ApiServer) Restart() {
	err := s.Server.Shutdown(context.Background()); 
	if err != nil {
		log.Println("API->Restart: ", err)
	} else {
		s.LaunchServer()
	}
}

func (s *ApiServer) GetStatus() string {
	return s.Status
}

//Action
func (s *ApiServer) Action(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	_, writeErr := w.Write([]byte("OK"))
	if writeErr != nil {
		log.Println(writeErr)
	}
	
	action := ps.ByName("type")
	if action == "Shutdown" {
		s.Status = "Shutdown"
		s.Shutdown()
	} else {
		s.Status = "Restart"
		s.Restart()
	}
}

//Invalid
func Error404(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
    fmt.Fprint(w, "Invalid API call")
}

//Authentication (How do I decouple this without causing another request check...)
func (s *ApiServer) CheckAuthentication(r *http.Request) (int, error) {
	//Create Request
	client := &http.Client{}
	url := GetServerURL("Auth") + "/check"
	req, reqErr := http.NewRequest("GET", url, nil)
	if reqErr != nil {
		log.Println("API->CheckAuthentication->Req: ", reqErr)
		return -1, reqErr
	}
	req.Header = r.Header
	
	//Do Request
	resp, getErr := client.Do(req)
	if getErr != nil {
		log.Println("API->CheckAuthentication->Get: ", getErr)
		return -1, getErr
	}
	
	//Read Body
	body, bodyErr := io.ReadAll(resp.Body)
	if bodyErr != nil {
		log.Println("API->CheckAuthentication->Body: ", bodyErr)
		return -1, bodyErr
	}
	
	//Parse User ID
	userID, parseErr := strconv.Atoi(string(body[:]))
	if parseErr != nil {
		log.Println("API->CheckAuthentication->Parse", parseErr)
		return -1, parseErr
	}
	
	if userID < 1 {
		return -1, errors.New("Invalid of Expired Access Token")
	}
	
	return userID, nil
} //checks authorization and returns userID

//Models
func (s *ApiServer) GetModel(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	//Check Auth and Get UserID
	userID, authErr := s.CheckAuthentication(r)
	if authErr != nil {
		w.Write([]byte(authErr.Error()))
		return
	}

	//Get Model ID from request
	modelID, convErr := strconv.Atoi(ps.ByName("modelID"))
	if convErr != nil || modelID == 0 {
		w.Write([]byte("Model ID must be a non-zero number"))
		return
	}

	//Retrieve Model
	model := s.Routes.GetModel(modelID, userID)
	_, writeErr := w.Write(model)
	if writeErr != nil {
		log.Println(writeErr)
	}
} //API/Model/Get GET

func (s *ApiServer) AddModel(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	//Check Auth and Get UserID
	userID, authErr := s.CheckAuthentication(r)
	if authErr != nil {
		w.Write([]byte(authErr.Error()))
		return
	}

	//Get POST Data
	r.ParseForm()
	name := r.PostFormValue("name")
	desc := r.PostFormValue("desc")
	
	//Run API Function
	resp := s.Routes.AddModel(name, desc, userID)
	
	//Return Response
	_, writeErr := w.Write(resp)
	if writeErr != nil {
		log.Println(writeErr)
	}
} //API/Model/Add POST

func (s *ApiServer) UpdateModel(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	//Check Auth and Get UserID
	userID, authErr := s.CheckAuthentication(r)
	if authErr != nil {
		w.Write([]byte(authErr.Error()))
		return
	}

	//Get POST Data
	r.ParseForm()
	modelID := r.PostFormValue("modelID")
	name := r.PostFormValue("name")
	desc := r.PostFormValue("desc")
	
	//Run API Function
	resp := s.Routes.UpdateModel(modelID, name, desc, userID)
	
	//Return Response
	_, writeErr := w.Write(resp)
	if writeErr != nil {
		log.Println(writeErr)
	}
} //API/Model/Update POST

func (s *ApiServer) DeleteModel(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	//Get User
	userID := 1
	
	//Get Model ID from request
	modelID, convErr := strconv.Atoi(ps.ByName("modelID"))
	if convErr != nil || modelID == 0 {
		w.Write([]byte("Model ID must be a non-zero number"))
		return
	}
	
	//Delete Model
	model := s.Routes.DeleteModel(modelID, userID)
	_, writeErr := w.Write(model)
	if writeErr != nil {
		log.Println(writeErr)
	}
} //API/Model/Delete GET?

//Images
func (s *ApiServer) GetImages(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	//Check Auth and Get UserID
	userID, authErr := s.CheckAuthentication(r)
	if authErr != nil {
		w.Write([]byte(authErr.Error()))
		return
	}
	
	//Get Model ID from request
	modelID, convErr := strconv.Atoi(ps.ByName("modelID"))
	if convErr != nil || modelID == 0 {
		w.Write([]byte("Model ID must be a non-zero number"))
		return
	}

	//Retrieve Images
	model := s.Routes.GetImages(modelID, userID) //modelID int, userID int
	_, writeErr := w.Write(model)
	if writeErr != nil {
		log.Println(writeErr)
	}
} //API/Image/Get GET

func (s *ApiServer) AddImage(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	//Check Auth and Get UserID
	userID, authErr := s.CheckAuthentication(r)
	if authErr != nil {
		w.Write([]byte(authErr.Error()))
		return
	}

	//Get POST Data
	r.ParseForm()
	modelID := r.PostFormValue("modelID")
	path := r.PostFormValue("path")
	desc := r.PostFormValue("desc")
	
	//Run API Function
	resp := s.Routes.AddImage(modelID, path, desc, userID)
	
	//Return Response
	_, writeErr := w.Write(resp)
	if writeErr != nil {
		log.Println(writeErr)
	}
} //API/Image/Add POST

func (s *ApiServer) UpdateImage(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	//Check Auth and Get UserID
	userID, authErr := s.CheckAuthentication(r)
	if authErr != nil {
		w.Write([]byte(authErr.Error()))
		return
	}

	//Get POST Data
	r.ParseForm()
	imageID := r.PostFormValue("imageID")
	modelID := r.PostFormValue("modelID")
	path := r.PostFormValue("path")
	desc := r.PostFormValue("desc")
	
	//Run API Function
	resp := s.Routes.UpdateImage(imageID, modelID, path, desc, userID)
	
	//Return Response
	_, writeErr := w.Write(resp)
	if writeErr != nil {
		log.Println(writeErr)
	}
} //API/Image/Update POST

func (s *ApiServer) DeleteImage(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	//Check Auth and Get UserID
	userID, authErr := s.CheckAuthentication(r)
	if authErr != nil {
		w.Write([]byte(authErr.Error()))
		return
	}
	
	//Get Model ID from request
	imageID, convErr := strconv.Atoi(ps.ByName("imageID"))
	if convErr != nil || imageID == 0 {
		w.Write([]byte("Image ID must be a non-zero number"))
		return
	}
	
	//Delete Image
	image := s.Routes.DeleteImage(imageID, userID)
	_, writeErr := w.Write(image)
	if writeErr != nil {
		log.Println(writeErr)
	}
} //API/Image/Delete GET?

//Videos
func (s *ApiServer) GetVideos(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	//Check Auth and Get UserID
	userID, authErr := s.CheckAuthentication(r)
	if authErr != nil {
		w.Write([]byte(authErr.Error()))
		return
	}
	
	//Get Model ID from request
	modelID, convErr := strconv.Atoi(ps.ByName("modelID"))
	if convErr != nil || modelID == 0 {
		w.Write([]byte("Model ID must be a non-zero number"))
		return
	}

	//Retrieve Videos
	model := s.Routes.GetVideos(modelID, userID)
	_, writeErr := w.Write(model)
	if writeErr != nil {
		log.Println(writeErr)
	}
	//
} //API/Video/Get POST

func (s *ApiServer) AddVideo(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	//Check Auth and Get UserID
	userID, authErr := s.CheckAuthentication(r)
	if authErr != nil {
		w.Write([]byte(authErr.Error()))
		return
	}

	//Get POST Data
	r.ParseForm()
	modelID := r.PostFormValue("modelID")
	path := r.PostFormValue("path")
	desc := r.PostFormValue("desc")
	
	//Run API Function
	resp := s.Routes.AddVideo(modelID, path, desc, userID)
	
	//Return Response
	_, writeErr := w.Write(resp)
	if writeErr != nil {
		log.Println(writeErr)
	}
} //API/Video/Add POST

func (s *ApiServer) UpdateVideo(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	//Check Auth and Get UserID
	userID, authErr := s.CheckAuthentication(r)
	if authErr != nil {
		w.Write([]byte(authErr.Error()))
		return
	}

	//Get POST Data
	r.ParseForm()
	videoID := r.PostFormValue("videoID")
	modelID := r.PostFormValue("modelID")
	path := r.PostFormValue("path")
	desc := r.PostFormValue("desc")
	
	//Run API Function
	resp := s.Routes.UpdateVideo(videoID, modelID, path, desc, userID)
	
	//Return Response
	_, writeErr := w.Write(resp)
	if writeErr != nil {
		log.Println(writeErr)
	}
} //API/Video/Update POST

func (s *ApiServer) DeleteVideo(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	//Check Auth and Get UserID
	userID, authErr := s.CheckAuthentication(r)
	if authErr != nil {
		w.Write([]byte(authErr.Error()))
		return
	}
	
	//Get Model ID from request
	videoID, convErr := strconv.Atoi(ps.ByName("videoID"))
	if convErr != nil || videoID == 0 {
		w.Write([]byte("Video ID must be a non-zero number"))
		return
	}
	
	//Delete Video
	image := s.Routes.DeleteVideo(videoID, userID)
	_, writeErr := w.Write(image)
	if writeErr != nil {
		log.Println(writeErr)
	}
} //API/Video/Delete GET?

//Server Status
func (s *ApiServer) CheckStatus(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	_, writeErr := w.Write([]byte(s.GetStatus()))
	if writeErr != nil {
		log.Println(writeErr)
	}
}

//Download Report
func (s *ApiServer) DownloadReport(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	report := s.API.NewReport()
	reportBytes, reportErr := report.Download()
	if reportErr != nil {
		log.Println("API->Server->DownloadReport", reportErr)
	}
	
	_, writeErr := w.Write(reportBytes)
	if writeErr != nil {
		log.Println(writeErr)
	}
}