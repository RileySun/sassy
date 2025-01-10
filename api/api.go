package api

import(
	"fmt"
	"log"
	"errors"
)

//The actual API isnt as important as using the api, so we're gonna go with a basic one. 
//Routing to the api will be fake.
//column parameters in GET functions must include backticks because golang mysql driver is weird

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

//User
func (a *API) GetUserBy(column string, identifier any) *User {
	//Get Data
	rows, err := a.db.db.Query("SELECT `id`, `name`, `trial`, `get`, `add`, `update`, `delete` FROM Users WHERE ? = ?;", column, identifier)
	
	if err != nil {
		log.Println(err, "- API:GetUserBy - Rows")
	}
	defer rows.Close()
	
	//Extract Data to struct object
	user := &User{}
	for rows.Next() {
		err = rows.Scan(&user.ID, &user.Name, &user.Trial, &user.Get, &user.Add, &user.Update, &user.Delete); 
		if err != nil {
			log.Println(fmt.Errorf("User Where %s = %s", column, identifier))
		}
    }
	
	//If errors return nil, if not return struct object
	if err = rows.Err(); err != nil {
		log.Println(err, "- API:GetUserBy - Completion")
		return nil
	} else {
		return user
	}
}

//Model
func (a *API) GetModelBy(column string, identifier any) (*Model, error) {
	//Get Data
	rows, err := a.db.db.Query("SELECT * FROM Models WHERE ? = ?;", column, identifier)
	if err != nil {
		log.Println(err, "- API:GetModelBy - Rows")
		return nil, errors.New("Error Model/Get")
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
		return nil, errors.New("Error Model/Get")
	} else {
		return model, nil
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
func (a *API) GetImagesBy(column string, identifier any) ([]*Image, error) {
	//Get Data
	rows, err := a.db.db.Query("SELECT * FROM Images WHERE ? = ?;", column, identifier)
	if err != nil {
		log.Println(err, "- API:GetImagesBy - Rows")
		return nil, errors.New("Error Images/Get")
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
		return nil, errors.New("Error Images/Get")
	} else {
		return images, nil
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

func (a *API) DeleteImage(id int) error {
	statement, err := a.db.db.Prepare("DELETE FROM Images WHERE `id` = ?;")
    if err != nil {
        log.Println(err, "- API:DeleteImage")
    }
    _, err = statement.Exec(id)
    return err
}

//Videos
func (a *API) GetVideosBy(column string, identifier any) ([]*Video, error) {
	//Get Data
	rows, err := a.db.db.Query("SELECT * FROM Videos WHERE ? = ?;", column, identifier)
	if err != nil {
		log.Println(err, "- API:GetVideosBy - Rows")
		return nil, errors.New("Error Videos/Get")
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
		return nil, errors.New("Error Videos/Get")
	} else {
		return videos, nil
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

func (a *API) DeleteVideo(id int) error {
	statement, err := a.db.db.Prepare("DELETE FROM Videos WHERE `id` = ?;")
    if err != nil {
        log.Println(err, "- API:DeleteVideo")
    }
    _, err = statement.Exec(id)
    return err
}