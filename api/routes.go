package api

import(
	"log"
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
	 quotaOK := r.API.IsUnderQuota(userID, "Get")
	
	if quotaOK == nil {
		model, err := r.API.GetModelBy("id", modelID)
		if err != nil {
			return []byte("Error Model/Get, please contact administrator ")
		}
		
		quotaErr := r.API.AddToQuota(userID, "Get")
		if quotaErr != nil {
			log.Println("Quota error for user '" + strconv.Itoa(userID) + "' using method Get (Model)")
		}
		
		return model.JSON()
	} else {
		return []byte("User quota limit reached: Get")
	}
} //API/Model/Get

func (r *Routes) AddModel(name string, desc string, userID int) []byte {
	 quotaOK := r.API.IsUnderQuota(userID, "Add")
	
	if quotaOK == nil {
		err := r.API.AddModel(name, desc)
		if err != nil {
			return []byte("Error Model/Add, please contact administrator ")
		}
		
		quotaErr := r.API.AddToQuota(userID, "Add")
		if quotaErr != nil {
			log.Println("Quota error for user '" + strconv.Itoa(userID) + "' using method Add (Model)")
		}
		
		return []byte("Model `" + name + "` added successfully")
	} else {
		return []byte("User quota limit reached: Add")
	}
} //API/Model/Add

func (r *Routes) UpdateModel(modelID string, name string, desc string, userID int) []byte {
	//Validate modelID
	modelInt, convErr := strconv.Atoi(modelID)
	if convErr != nil {
		return []byte("Model ID must be a number")
	}
	
	quotaOK := r.API.IsUnderQuota(userID, "Update")
	
	if quotaOK == nil {
		err := r.API.UpdateModel(modelInt, name, desc)
		if err != nil {
			return []byte("Error Model/Update, please contact administrator ")
		}
		
		quotaErr := r.API.AddToQuota(userID, "Update")
		if quotaErr != nil {
			log.Println("Quota error for user '" + strconv.Itoa(userID) + "' using method Update (Model)")
		}
		
		return []byte("Model " + modelID + " updated successfully")
	} else {
		return []byte("User quota limit reached: Update")
	}
} //API/Model/Update

func (r *Routes) DeleteModel(modelID int, userID int) []byte {
	 quotaOK := r.API.IsUnderQuota(userID, "Delete")
	
	if quotaOK == nil {
		err := r.API.DeleteModel(modelID)
		if err != nil {
			return []byte("Error Model/Delete, please contact administrator ")
		}

		quotaErr := r.API.AddToQuota(userID, "Delete")
		if quotaErr != nil {
			log.Println("Quota error for user '" + strconv.Itoa(userID) + "' using method Delete (Model)")
		}
		
		return []byte("Model " + strconv.Itoa(modelID) + " deleted successfully")
	} else {
		return []byte("User quota limit reached: Delete")
	}
} //API/Model/Delete

//Images
func (r *Routes) GetImages(modelID int, userID int) []byte {
	 quotaOK := r.API.IsUnderQuota(userID, "Get")
	
	if quotaOK == nil {
		images, err := r.API.GetImagesBy("model_id", modelID)
		if err != nil {
			return []byte("Error Image/Get, please contact administrator ")
		}
		
		quotaErr := r.API.AddToQuota(userID, "Get")
		if quotaErr != nil {
			log.Println("Quota error for user '" + strconv.Itoa(userID) + "' using method Get (Images)")
		}
		
		var output []byte
		for _, i := range images {
			output = append(output, i.JSON()...)
		}
		return output
	} else {
		return []byte("User quota limit reached: Get")
	}
} //API/Image/Get

func (r *Routes) AddImage(modelID string, path string, desc string, userID int) []byte {
	//Validate modelID
	modelInt, convErr := strconv.Atoi(modelID)
	if convErr != nil {
		return []byte("Model ID must be a number")
	}
	
	quotaOK := r.API.IsUnderQuota(userID, "Add")
	
	if quotaOK == nil {
		err := r.API.AddImage(modelInt, path, desc)
		if err != nil {
			return []byte("Error Image/Add, please contact administrator ")
		}
		
		quotaErr := r.API.AddToQuota(userID, "Add")
		if quotaErr != nil {
			log.Println("Quota error for user '" + strconv.Itoa(userID) + "' using method Add (Images)")
		}
		
		return []byte("Image for model " + modelID + " added successfully")
	} else {
		return []byte("User quota limit reached: Add")
	}
} //API/Image/Add

func (r *Routes) UpdateImage(imageID string, modelID string, path string, desc string, userID int) []byte {
	//Validate imageID
	imageInt, convErr := strconv.Atoi(imageID)
	if convErr != nil {
		return []byte("Image ID must be a number")
	}
	
	//Validate modelID
	modelInt, convErr := strconv.Atoi(modelID)
	if convErr != nil {
		modelInt = 0
	}

	 quotaOK := r.API.IsUnderQuota(userID, "Update")
	
	if quotaOK == nil {
		err := r.API.UpdateImage(imageInt, modelInt, path, desc)
		if err != nil {
			if err.Error() == "Invalid Image ID" {
				return []byte("Invalid Image ID")
			} else {
				return []byte("Error Image/Update, please contact administrator ")
			}
		}
		
		quotaErr := r.API.AddToQuota(userID, "Update")
		if quotaErr != nil {
			log.Println("Quota error for user '" + strconv.Itoa(userID) + "' using method Update (Images)")
		}
		
		return []byte("Image " + imageID + " updated successfully")
	} else {
		return []byte("User quota limit reached: Update")
	}
} //API/Image/Update

func (r *Routes) DeleteImage(imageID int, userID int) []byte {
	 quotaOK := r.API.IsUnderQuota(userID, "Delete")
	
	if quotaOK == nil {
		err := r.API.DeleteImage(imageID)
		if err != nil {
			if err.Error() == "Invalid Image ID" {
				return []byte("Invalid Image ID")
			} else {
				return []byte("Error Image/Delete, please contact administrator ")
			}
		}
		
		quotaErr := r.API.AddToQuota(userID, "Delete")
		if quotaErr != nil {
			log.Println("Quota error for user '" + strconv.Itoa(userID) + "' using method Delete (Images)")
		}
		
		return []byte("Image " + strconv.Itoa(imageID) + " deleted successfully")
	} else {
		return []byte("User quota limit reached: Delete")
	}
} //API/Image/Delete

//Videos
func (r *Routes) GetVideos(modelID int, userID int) []byte {
	 quotaOK := r.API.IsUnderQuota(userID, "Get")
	
	if quotaOK == nil {
		videos, err := r.API.GetVideosBy("model_id", modelID)
		if err != nil {
			return []byte("Error Video/Get, please contact administrator ")
		}
		
		quotaErr := r.API.AddToQuota(userID, "Get")
		if quotaErr != nil {
			log.Println("Quota error for user '" + strconv.Itoa(userID) + "' using method Get (Videos)")
		}
		
		var output []byte
		for _, v := range videos {
			output = append(output, v.JSON()...)
		}
		return output
	} else {
		return []byte("User quota limit reached: Get")
	}
} //API/Video/Get

func (r *Routes) AddVideo(modelID string, path string, desc string, userID int) []byte {
	modelInt, convErr := strconv.Atoi(modelID)
	if convErr != nil {
		return []byte("Model ID must be a number")
	}
	
	quotaOK := r.API.IsUnderQuota(userID, "Add")
	
	if quotaOK == nil {
		err := r.API.AddVideo(modelInt, path, desc)
		if err != nil {
			return []byte("Error Video/Add, please contact administrator ")
		}
		
		quotaErr := r.API.AddToQuota(userID, "Add")
		if quotaErr != nil {
			log.Println("Quota error for user '" + strconv.Itoa(userID) + "' using method Add (Videos)")
		}
		
		return []byte("Video for model " + modelID + " added successfully")
	} else {
		return []byte("User quota limit reached: Add")
	}
} //API/Video/Add

func (r *Routes) UpdateVideo(videoID string, modelID string, path string, desc string, userID int) []byte {
	//Validate imageID
	videoInt, convErr := strconv.Atoi(videoID)
	if convErr != nil {
		return []byte("Video ID must be a number")
	}
	
	//Validate modelID
	modelInt, convErr := strconv.Atoi(modelID)
	if convErr != nil {
		modelInt = 0
	}
	
	quotaOK := r.API.IsUnderQuota(userID, "Update")
	
	if quotaOK == nil {
		err := r.API.UpdateVideo(videoInt, modelInt, path, desc)
		if err != nil {
			if err.Error() == "Invalid Video ID" {
				return []byte("Invalid Video ID")
			} else {
				return []byte("Error Video/Delete, please contact administrator ")
			}
		}
		
		quotaErr := r.API.AddToQuota(userID, "Update")
		if quotaErr != nil {
			log.Println("Quota error for user '" + strconv.Itoa(userID) + "' using method Update (Videos)")
		}
		
		return []byte("Video " + videoID + " updated successfully")
	} else {
		return []byte("User quota limit reached: Update")
	}
} //API/Video/Update

func (r *Routes) DeleteVideo(videoID int, userID int) []byte {
	 quotaOK := r.API.IsUnderQuota(userID, "Delete")
	
	if quotaOK == nil {
		err := r.API.DeleteVideo(videoID)
		if err != nil {
			if err.Error() == "Invalid Video ID" {
				return []byte("Invalid Video ID")
			} else {
				return []byte("Error Video/Delete, please contact administrator ")
			}
		}
		
		quotaErr := r.API.AddToQuota(userID, "Delete")
		if quotaErr != nil {
			log.Println("Quota error for user '" + strconv.Itoa(userID) + "' using method Delete (Videos)")
		}
		
		return []byte("Video " + strconv.Itoa(videoID) + " deleted successfully")
	} else {
		return []byte("User quota limit reached: Delete")
	}
} //API/Video/Delete