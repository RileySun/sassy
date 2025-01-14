package api

import(
	"os"
	"strconv"
	"testing"
	"math/rand"
)

var api *API
var model *Model
var testImage *Image
var testVideo *Video

//Main
func TestMain(m *testing.M) {
	//Others
	AuthTestInit()
	
	//This one
	api = NewAPI()
	exitCode := m.Run()
	os.Exit(exitCode)
}

//Model (Delete After Images & Videos)
func TestAddModel(t *testing.T) {
	modelErr := api.AddModel("Leighton Stultz", "Reality TV Contestant")
	if modelErr != nil {
		t.Error(modelErr.Error())
		t.Fail()
	}
}

func TestGetModelBy(t *testing.T) {
	newModel, getErr := api.GetModelBy("name", "Leighton Stultz")
	if getErr != nil {
		t.Error(getErr.Error())
	}
	model = newModel
}

func TestUpdateModel(t *testing.T) {
	updateErr := api.UpdateModel(model.ID, "Anna Faris", "Former Movie Actress")
	if updateErr != nil {
		t.Error(updateErr.Error())
	}
	
	newModel, getErr := api.GetModelBy("id", model.ID)
	if getErr != nil {
		t.Error(getErr.Error())
	}
	
	if newModel.Name == model.Name || newModel.Desc == model.Desc {
		t.Error("Models should not match")
	}
}

//Images
func TestAddImages(t *testing.T) {
	for i := 0; i < 5; i++ {
		path := "/Images/Test/" + strconv.Itoa(i) + ".jpg"
		desc := "Tester Image " + strconv.Itoa(i)
		addErr := api.AddImage(model.ID, path, desc)
		if addErr != nil {
			t.Error(addErr.Error())
			t.Fail()
		}
	}
}

func TestGetImagesBy(t *testing.T) {
	images, getErr := api.GetImagesBy("model_id", model.ID)
	if getErr != nil {
		t.Error(getErr.Error())
	}
	
	//Is right amount of images
	if len(images) != 5 {
		t.Error("Should be 5 test images")
	}
	
	//Sample random image
	i := rand.Intn(4)
	correctPath := "/Images/Test/" + strconv.Itoa(i) + ".jpg"
	correctDesc := "Tester Image " + strconv.Itoa(i)
	if images[i].Path != correctPath || images[i].Desc != correctDesc {
		t.Error("Random sampling of paths and descs do not align.")
	}
	
	testImage = images[i]
}

func TestUpdateImage(t *testing.T) {
	updateErr := api.UpdateImage(testImage.ID, testImage.Model_ID, "nil", "nil")
	if updateErr != nil {
		t.Error(updateErr.Error())
	}
	
	//Get Updated Data
	images, getErr := api.GetImagesBy("id", testImage.ID)
	if getErr != nil {
		t.Error(getErr.Error())
	}
	
	//Compare to old data
	if images[0].Path == testImage.Path || images[0].Desc == testImage.Desc {
		t.Error("Image data did not update")
	}
}

func TestDeleteImages(t *testing.T) {
	//Get image Id's by Model
	images, getErr := api.GetImagesBy("model_id", testImage.Model_ID)
	if getErr != nil {
		t.Error(getErr.Error())
	}
	
	//Delete Images
	for _, i := range images {
		deleteErr := api.DeleteImage(i.ID)
		if deleteErr != nil {
			t.Error(deleteErr.Error())
		}
	}
	
	//Check if deleted
	deletedImages, checkErr := api.GetImagesBy("model_id", testImage.Model_ID)
	if checkErr != nil {
		t.Error(checkErr.Error())
	}
	if len(deletedImages) > 0 {
		t.Error("Images were not deleted.")
	}
}

//Videos
func TestAddVideos(t *testing.T) {
	for i := 0; i < 5; i++ {
		path := "/Videos/Test/" + strconv.Itoa(i) + ".mp4"
		desc := "Tester Video " + strconv.Itoa(i)
		addErr := api.AddVideo(model.ID, path, desc)
		if addErr != nil {
			t.Error(addErr.Error())
			t.Fail()
		}
	}
}

func TestGetVideosBy(t *testing.T) {
	videos, getErr := api.GetVideosBy("model_id", model.ID)
	if getErr != nil {
		t.Error(getErr.Error())
	}
	
	//Is right amount of videos
	if len(videos) != 5 {
		t.Error("Should be 5 test videos")
	}
	
	//Sample random video
	i := rand.Intn(4)
	correctPath := "/Videos/Test/" + strconv.Itoa(i) + ".mp4"
	correctDesc := "Tester Video " + strconv.Itoa(i)
	if videos[i].Path != correctPath || videos[i].Desc != correctDesc {
		t.Error("Random sampling of paths and descs do not align.")
	}
	
	testVideo = videos[i]
}

func TestUpdateVideo(t *testing.T) {
	updateErr := api.UpdateVideo(testVideo.ID, testVideo.Model_ID, "nil", "nil")
	if updateErr != nil {
		t.Error(updateErr.Error())
	}
	
	//Get Updated Data
	videos, getErr := api.GetVideosBy("id", testVideo.ID)
	if getErr != nil {
		t.Error(getErr.Error())
	}
	
	//Compare to old data
	if videos[0].Path == testVideo.Path || videos[0].Desc == testVideo.Desc {
		t.Error("Video data did not update")
	}
}

func TestDeleteVideos(t *testing.T) {
	//Get video Id's by Model
	videos, getErr := api.GetVideosBy("model_id", testVideo.Model_ID)
	if getErr != nil {
		t.Error(getErr.Error())
	}
	
	//Delete Videos
	for _, i := range videos {
		deleteErr := api.DeleteVideo(i.ID)
		if deleteErr != nil {
			t.Error(deleteErr.Error())
		}
	}
	
	//Check if deleted
	deletedVideos, checkErr := api.GetVideosBy("model_id", testVideo.Model_ID)
	if checkErr != nil {
		t.Error(checkErr.Error())
	}
	if len(deletedVideos) > 0 {
		t.Error("Videos were not deleted.")
	}
}

//Delete Model
func TestDeleteModel(t *testing.T) {
	deleteErr := api.DeleteModel(model.ID)
	if deleteErr != nil {
		t.Error(deleteErr.Error())
	}
	
	_, getErr := api.GetModelBy("id", model.ID)
	if getErr == nil {
		t.Error("Model should not exist")
	}
}

//Column Name Sanitization
func TestAPICheckColumnName(t *testing.T) {
	shouldPass := api.checkColumnName("id")
	if shouldPass != nil {
		t.Error(shouldPass.Error())
	}
	
	shouldFail := api.checkColumnName("test")
	if shouldFail == nil {
		t.Error("Should not be a valid column name")
	}
}