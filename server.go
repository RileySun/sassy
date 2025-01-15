package main

import(
	"api"
	
	"log"
	"fmt"
	"strconv"
	"net/http"
	
	"github.com/julienschmidt/httprouter"
)


type Server struct {
	API *api.API
	Auth *api.Auth
	Routes *api.Routes
	
	router *httprouter.Router
}

//Create
func NewServer() *Server {
	server := &Server{
		API:api.NewAPI(),
		Auth:api.NewAuth(),
		router:httprouter.New(),
	}
	server.Routes = server.API.NewRoutes()
	
	return server
}

//Launch
func (s *Server) Launch() {
	//Models
	s.router.GET("/model/get/:modelID", s.GetModel)
	
	//Images
	s.router.GET("/images/get/:modelID", s.GetImages)
	
	//Videos
	s.router.GET("/videos/get/:modelID", s.GetVideos)
	
	//Media
	
	
	s.router.GET("/", Error404)

	//Launch
	log.Fatal(http.ListenAndServe("localhost:8080", s.router))
}

//Invalid
func Error404(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
    fmt.Fprint(w, "Invalid API call")
}

//Authentication


//Models
func (s *Server) GetModel(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	//Get Model ID from request
	modelID, convErr := strconv.Atoi(ps.ByName("modelID"))
	if convErr != nil || modelID == 0 {
		w.Write([]byte("Model ID must be a non-zero number"))
		return
	}

	//
	model := s.Routes.GetModel(modelID, 1)//modelID int, userID int
	_, writeErr := w.Write(model)
	if writeErr != nil {
		log.Println(writeErr)
	}
} //API/Model/Get GET

func (s *Server) AddModel(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	//s.Routes.AddModel(name string, desc string, userID int)
} //API/Model/Add POST

func (s *Server) UpdateModel(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	//s.Routes.UpdateModel(modelID int, name string, desc string, userID int)
} //API/Model/Update POST

func (s *Server) DeleteModel(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	//s.Routes.DeleteModel(modelID int, userID int)
} //API/Model/Delete GET?

//Images
func (s *Server) GetImages(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	//Get Model ID from request
	modelID, convErr := strconv.Atoi(ps.ByName("modelID"))
	if convErr != nil || modelID == 0 {
		w.Write([]byte("Model ID must be a non-zero number"))
		return
	}

	//
	model := s.Routes.GetImages(modelID, 1) //modelID int, userID int
	_, writeErr := w.Write(model)
	if writeErr != nil {
		log.Println(writeErr)
	}
} //API/Image/Get GET

func (s *Server) AddImage(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	//s.Routes.AddImage(modelID int, path string, desc string, userID int)
} //API/Image/Add POST

func (s *Server) UpdateImage(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	//s.Routes.UpdateImage(imageID int, modelID int, path string, desc string, userID int)
} //API/Image/Update POST

func (s *Server) DeleteImage(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	//s.Routes.DeleteImage(imageID int, userID int)
} //API/Image/Delete GET?

//Videos
func (s *Server) GetVideos(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	//Get Model ID from request
	modelID, convErr := strconv.Atoi(ps.ByName("modelID"))
	if convErr != nil || modelID == 0 {
		w.Write([]byte("Model ID must be a non-zero number"))
		return
	}

	//
	model := s.Routes.GetVideos(modelID, 1) //modelID int, userID int
	_, writeErr := w.Write(model)
	if writeErr != nil {
		log.Println(writeErr)
	}
	//
} //API/Video/Get POST

func (s *Server) AddVideo(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	 //s.Routes.AddVideo(modelID int, path string, desc string, userID int)
} //API/Video/Add POST

func (s *Server) UpdateVideo(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	 //s.Routes.UpdateVideo(imageID int, modelID int, path string, desc string, userID int)
} //API/Video/Update POST

func (s *Server) DeleteVideo(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	//s.Routes.DeleteVideo(videoID int, userID int)
} //API/Video/Delete GET?