package main

import(
	"api"
	
	"log"
	"fmt"
	"strconv"
	"net/http"
	
	"github.com/julienschmidt/httprouter"
)

//Header["Authorization"] then i guess substr it for bearer token?

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
	
	//Error
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

	//Retrieve Model
	model := s.Routes.GetModel(modelID, 1)//modelID int, userID int
	_, writeErr := w.Write(model)
	if writeErr != nil {
		log.Println(writeErr)
	}
} //API/Model/Get GET

func (s *Server) AddModel(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	//Get User
	userID := 1

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

func (s *Server) UpdateModel(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	//Get User
	userID := 1

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

func (s *Server) DeleteModel(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
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
func (s *Server) GetImages(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	//Get Model ID from request
	modelID, convErr := strconv.Atoi(ps.ByName("modelID"))
	if convErr != nil || modelID == 0 {
		w.Write([]byte("Model ID must be a non-zero number"))
		return
	}

	//Retrieve Images
	model := s.Routes.GetImages(modelID, 1) //modelID int, userID int
	_, writeErr := w.Write(model)
	if writeErr != nil {
		log.Println(writeErr)
	}
} //API/Image/Get GET

func (s *Server) AddImage(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	//Get User
	userID := 1

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

func (s *Server) UpdateImage(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	//Get User
	userID := 1

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

func (s *Server) DeleteImage(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	//Get User
	userID := 1
	
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
func (s *Server) GetVideos(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	//Get Model ID from request
	modelID, convErr := strconv.Atoi(ps.ByName("modelID"))
	if convErr != nil || modelID == 0 {
		w.Write([]byte("Model ID must be a non-zero number"))
		return
	}

	//Retrieve Videos
	model := s.Routes.GetVideos(modelID, 1) //modelID int, userID int
	_, writeErr := w.Write(model)
	if writeErr != nil {
		log.Println(writeErr)
	}
	//
} //API/Video/Get POST

func (s *Server) AddVideo(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	//Get User
	userID := 1

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

func (s *Server) UpdateVideo(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	//Get User
	userID := 1

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

func (s *Server) DeleteVideo(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	//Get User
	userID := 1
	
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