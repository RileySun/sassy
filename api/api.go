package api

import(
	"fmt"
	"log"
)

//The actual API isnt as important as using the api, so we're gonna go with a basic one. 
//Routing to the api will be fake.

//Object
type API struct {
	db *Database
}

//Create
func NewAPI() *API {
	return &API{
		db:NewDB(),
	}
}


//Model
func (a *API) GetModelBy(column string, identifier any) *Model {
	//Get Data
	rows, err := a.db.db.Query("SELECT * FROM Models WHERE `?` = ?;", column, identifier)
	if err != nil {
		log.Println(err, "- API:GetModelBy - Rows")
	}
	defer rows.Close()
	
	//Extract Data to struct object
	model := &Model{}
	for rows.Next() {
		err = rows.Scan(&model.ID, &model.Name, &model.Desc); 
		if err != nil {
			log.Println(fmt.Errorf("Model Where %s = %s", column, identifier))
		}
    }
	
	//If errors return nil, if not return struct object
	if err = rows.Err(); err != nil {
		log.Println(err, "- API:GetModelBy - Completion")
		return nil
	} else {
		return model
	}
}

func (a *API) AddModel(name string, desc string) error {
	statement, err := a.db.db.Prepare("INSERT INTO Models (`name`, `desc`) VALUES (?, ?)")
    if err != nil {
        log.Println(err, "- API:AddModel")
    }
    _, err = statement.Exec(name, desc)
    return err
}

func (a *API) UpdateModel(id int, name string, desc string) error {
	queryString := "UPDATE Models SET `name`=?, `desc`=? WHERE `id` = ?;"	
	statement, err := a.db.db.Prepare(queryString)
	 if err != nil {
        log.Println(err, "- API:UpdateModel")
    }
	 _, err = statement.Exec(name, desc, id)
    return err
}

func (a *API) DeleteModel(id int) error {
	statement, err := a.db.db.Prepare("DELETE FROM Models WHERE `id` = ?;")
    if err != nil {
        log.Println(err, "- API:DeleteModel")
    }
    _, err = statement.Exec(id)
    return err
}

//Images
func (a *API) GetImagesBy(column string, identifier any) []*Image {
	//Get Data
	rows, err := a.db.db.Query("SELECT * FROM Images WHERE ? = ?;", column, identifier)
	if err != nil {
		log.Println(err, "- API:GetImagesBy - Rows")
	}
	defer rows.Close()
	
	//Extract Data to struct object
	var images []*Image
	for rows.Next() {
		image := &Image{}
		err = rows.Scan(&image.ID, &image.Model_ID, &image.Path, &image.Desc); 
		if err != nil {
			log.Println(fmt.Errorf("Images Where %s = %s", column, identifier))
		}
		images = append(images, image)
    }
	
	//If errors return nil, if not return struct object
	if err = rows.Err(); err != nil {
		log.Println(err, "- API:GetImagesBy - Completion")
		return nil
	} else {
		return images
	}
}

func (a *API) AddImage(model_id int, path string, desc string) error {
	statement, err := a.db.db.Prepare("INSERT INTO Images (`model_id`, `path`, `desc`) VALUES (?, ?, ?)")
    if err != nil {
        log.Println(err, "- API:AddImage")
    }
    _, err = statement.Exec(model_id, path, desc)
    return err
}

func (a *API) UpdateImage(id int, model_id int, path string, desc string) error {
	queryString := "UPDATE Images SET `model_id`=?, `path`=?, `desc`=? WHERE `id` = ?;"	
	statement, err := a.db.db.Prepare(queryString)
	 if err != nil {
        log.Println(err, "- API:UpdateImage")
    }
	 _, err = statement.Exec(model_id, path, desc, id)
    return err
}

//Videos
func (a *API) GetVideosBy(column string, identifier any) []*Video {
	//Get Data
	rows, err := a.db.db.Query("SELECT * FROM Videos WHERE ? = ?;", column, identifier)
	if err != nil {
		log.Println(err, "- API:GetVideosBy - Rows")
	}
	defer rows.Close()
	
	//Extract Data to struct object
	var videos []*Video
	for rows.Next() {
		video := &Video{}
		err = rows.Scan(&video.ID, &video.Model_ID, &video.Path, &video.Desc); 
		if err != nil {
			log.Println(fmt.Errorf("Vidoes Where %s = %s", column, identifier))
		}
		videos = append(videos, video)
    }
	
	//If errors return nil, if not return struct object
	if err = rows.Err(); err != nil {
		log.Println(err, "- API:GetVideosBy - Completion")
		return nil
	} else {
		return videos
	}
}

func (a *API) AddVideo(model_id int, path string, desc string) error {
	statement, err := a.db.db.Prepare("INSERT INTO Videos (`model_id`, `path`, `desc`) VALUES (?, ?, ?)")
    if err != nil {
        log.Println(err, "- API:AddVideo")
    }
    _, err = statement.Exec(model_id, path, desc)
    return err
}

func (a *API) UpdateVideo(id int, model_id int, path string, desc string) error {
	queryString := "UPDATE Videos SET `model_id`=?, `path`=?, `desc`=? WHERE `id` = ?;"	
	statement, err := a.db.db.Prepare(queryString)
	 if err != nil {
        log.Println(err, "- API:UpdateVideo")
    }
	 _, err = statement.Exec(model_id, path, desc, id)
    return err
}