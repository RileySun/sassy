package api

import(
	"strconv"
)

//Struct
type Routes struct {
	API *API
}

//Create
func (a *API) NewRoutes() *Routes {
	return &Routes{API:a}
}

//Model
func (r *Routes) GetModel(modelID int, userID int) []byte {
	 quotaOK := r.API.IsUnderQuota(userID, "get")
	
	if quotaOK {
		model, err := r.API.GetModelBy("`id`", modelID)
		if err != nil {
			return []byte("Error Model/Get, please contact administrator ")
		}
		r.API.AddToQuota(userID, "Get")
		
		return model.JSON()
	} else {
		return []byte("User quota limit reached: Get")
	}
} //API/Model/Get

func (r *Routes) AddModel(name string, desc string, userID int) []byte {
	 quotaOK := r.API.IsUnderQuota(userID, "Add")
	
	if quotaOK {
		err := r.API.AddModel(name, desc)
		if err != nil {
			return []byte("Error Model/Add, please contact administrator ")
		}
		r.API.AddToQuota(userID, "Add")
		
		return []byte("Model `" + name + "` added successfully")
	} else {
		return []byte("User quota limit reached: Add")
	}
} //API/Model/Add

func (r *Routes) UpdateModel(modelID int, name string, desc string, userID int) []byte {
	 quotaOK := r.API.IsUnderQuota(userID, "Update")
	
	if quotaOK {
		err := r.API.UpdateModel(modelID, name, desc)
		if err != nil {
			return []byte("Error Model/Update, please contact administrator ")
		}
		r.API.AddToQuota(userID, "Update")
		
		return []byte("Model `" + strconv.Itoa(modelID) + "` updated successfully")
	} else {
		return []byte("User quota limit reached: Update")
	}
} //API/Model/Update

func (r *Routes) DeleteModel(modelID int, userID int) []byte {
	 quotaOK := r.API.IsUnderQuota(userID, "Delete")
	
	if quotaOK {
		err := r.API.DeleteModel(modelID)
		if err != nil {
			return []byte("Error Model/Delete, please contact administrator ")
		}
		r.API.AddToQuota(userID, "Delete")
		
		return []byte("Model `" + strconv.Itoa(modelID) + "` deleted successfully")
	} else {
		return []byte("User quota limit reached: Delete")
	}
} //API/Model/Delete

//Images
func (r *Routes) GetImages(modelID int, userID int) []byte {
	 quotaOK := r.API.IsUnderQuota(userID, "get")
	
	if quotaOK {
		images, err := r.API.GetImagesBy("`model_id`", modelID)
		if err != nil {
			return []byte("Error Image/Get, please contact administrator ")
		}
		r.API.AddToQuota(userID, "Get")
		
		var output []byte
		for _, i := range images {
			output = append(output, i.JSON()...)
		}
		return output
	} else {
		return []byte("User quota limit reached: Get")
	}
} //API/Image/Get

func (r *Routes) AddImage(modelID int, path string, desc string, userID int) []byte {
	 quotaOK := r.API.IsUnderQuota(userID, "Add")
	
	if quotaOK {
		err := r.API.AddImage(modelID, path, desc)
		if err != nil {
			return []byte("Error Image/Add, please contact administrator ")
		}
		r.API.AddToQuota(userID, "Add")
		
		return []byte("Image for model `" + strconv.Itoa(modelID) + " added successfully")
	} else {
		return []byte("User quota limit reached: Add")
	}
} //API/Image/Add

func (r *Routes) UpdateImage(imageID int, modelID int, path string, desc string, userID int) []byte {
	 quotaOK := r.API.IsUnderQuota(userID, "Update")
	
	if quotaOK {
		err := r.API.UpdateImage(imageID, modelID, path, desc)
		if err != nil {
			return []byte("Error Images/Update, please contact administrator ")
		}
		r.API.AddToQuota(userID, "Update")
		
		return []byte("Image `" + strconv.Itoa(modelID) + "` updated successfully")
	} else {
		return []byte("User quota limit reached: Update")
	}
} //API/Image/Update

func (r *Routes) DeleteImage(imageID int, userID int) []byte {
	 quotaOK := r.API.IsUnderQuota(userID, "Delete")
	
	if quotaOK {
		err := r.API.DeleteImage(imageID)
		if err != nil {
			return []byte("Error Image/Delete, please contact administrator ")
		}
		r.API.AddToQuota(userID, "Delete")
		
		return []byte("Image `" + strconv.Itoa(imageID) + "` deleted successfully")
	} else {
		return []byte("User quota limit reached: Delete")
	}
} //API/Image/Delete

//Videos
func (r *Routes) GetVideos(modelID int, userID int) []byte {
	 quotaOK := r.API.IsUnderQuota(userID, "get")
	
	if quotaOK {
		videos, err := r.API.GetVideosBy("`model_id`", modelID)
		if err != nil {
			return []byte("Error Video/Get, please contact administrator ")
		}
		r.API.AddToQuota(userID, "Get")
		
		var output []byte
		for _, v := range videos {
			output = append(output, v.JSON()...)
		}
		return output
	} else {
		return []byte("User quota limit reached: Get")
	}
} //API/Video/Get

func (r *Routes) AddVideo(modelID int, path string, desc string, userID int) []byte {
	 quotaOK := r.API.IsUnderQuota(userID, "Add")
	
	if quotaOK {
		err := r.API.AddVideo(modelID, path, desc)
		if err != nil {
			return []byte("Error Video/Add, please contact administrator ")
		}
		r.API.AddToQuota(userID, "Add")
		
		return []byte("Video for model `" + strconv.Itoa(modelID) + " added successfully")
	} else {
		return []byte("User quota limit reached: Add")
	}
} //API/Video/Add

func (r *Routes) UpdateVideo(imageID int, modelID int, path string, desc string, userID int) []byte {
	 quotaOK := r.API.IsUnderQuota(userID, "Update")
	
	if quotaOK {
		err := r.API.UpdateVideo(imageID, modelID, path, desc)
		if err != nil {
			return []byte("Error Videos/Update, please contact administrator ")
		}
		r.API.AddToQuota(userID, "Update")
		
		return []byte("Video `" + strconv.Itoa(modelID) + "` updated successfully")
	} else {
		return []byte("User quota limit reached: Update")
	}
} //API/Video/Update

func (r *Routes) DeleteVideo(videoID int, userID int) []byte {
	 quotaOK := r.API.IsUnderQuota(userID, "Delete")
	
	if quotaOK {
		err := r.API.DeleteVideo(videoID)
		if err != nil {
			return []byte("Error Video/Delete, please contact administrator ")
		}
		r.API.AddToQuota(userID, "Delete")
		
		return []byte("Video `" + strconv.Itoa(videoID) + "` deleted successfully")
	} else {
		return []byte("User quota limit reached: Delete")
	}
} //API/Video/Delete