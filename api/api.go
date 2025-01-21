package api

import(
	"errors"
	"slices"
	"strconv"
)

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
func (a *API) GetUserBy(column string, identifier any) (*User, error) {
	//Validate column names
	columnErr := a.checkColumnName(column)
	if columnErr != nil {
		return nil, columnErr
	}
	
	//Get Data
	user := &User{}
	row := a.db.db.QueryRow("SELECT `id`, `name`, `trial`, `get`, `add`, `update`, `delete` FROM Users WHERE `"+column+"` = ?;", identifier)
	scanErr := row.Scan(&user.ID, &user.Name, &user.Trial, &user.Get, &user.Add, &user.Update, &user.Delete); 
	if scanErr != nil {
		return nil, scanErr
	}
	
	return user, nil
}

func (a *API) GetAllUsers() ([]*User, error) {		
	//Get Data
	rows, rowErr := a.db.db.Query("SELECT `id`, `name`, `trial`, `get`, `add`, `update`, `delete` FROM Users")
	if rowErr != nil {
		return nil, rowErr
	}
	defer rows.Close()
	
	//Extract Data to struct object
	var users []*User
	for rows.Next() {
		user := &User{}
		scanErr := rows.Scan(&user.ID, &user.Name, &user.Trial, &user.Get, &user.Add, &user.Update, &user.Delete); 
		if scanErr != nil {
			return nil, scanErr
		}
		users = append(users, user)
    }
	
	//If errors return nil, if not return struct object
	if completeErr := rows.Err(); completeErr != nil {
		return nil, completeErr
	} else {
		return users, nil
	}
} 

//Model
func (a *API) GetModelBy(column string, identifier any) (*Model, error) {
	//Validate column names
	columnErr := a.checkColumnName(column)
	if columnErr != nil {
		return nil, columnErr
	}
	
	//Get Data
	model := &Model{}
	row := a.db.db.QueryRow("SELECT * FROM Models WHERE `"+column+"` = ?;", identifier)
	scanErr := row.Scan(&model.ID, &model.Name, &model.Desc);
	if scanErr != nil {
		return nil, scanErr
	}
	
	return model, nil
} //Singular

func (a *API) AddModel(name string, desc string) error {
	statement, prepErr := a.db.db.Prepare("INSERT INTO Models (`name`, `desc`) VALUES (?, ?)")
    if prepErr != nil {
        return prepErr
    }
    _, execErr := statement.Exec(name, desc)
    return execErr
}

func (a *API) UpdateModel(id int, name string, desc string) error {
	//Is this a real model
	oldModel, modelErr := a.GetModelBy("id", id)
	if modelErr != nil {
		return errors.New("Invalid Model ID")
	}
	
	//Missing or blank data means we do not change the field (can think of a more clever way to do this, but idk might be pitfalls)
	var finalName, finalDesc string
	if name == "" {
		finalName = oldModel.Name
	} else {
		finalName = name
	}
	if desc == "" {
		finalDesc = oldModel.Desc
	} else {
		finalDesc = desc
	}

	queryString := "UPDATE Models SET `name`=?, `desc`=? WHERE `id` = ?;"	
	statement, prepErr := a.db.db.Prepare(queryString)
	if prepErr != nil {
       return prepErr
    }
	 _, execErr := statement.Exec(finalName, finalDesc, id)
    return execErr
}

func (a *API) DeleteModel(id int) error {
	statement, prepErr := a.db.db.Prepare("DELETE FROM Models WHERE `id` = ?;")
    if prepErr != nil {
        return prepErr
    }
    _, execErr := statement.Exec(id)
    return execErr
}

//Images
func (a *API) GetImagesBy(column string, identifier any) ([]*Image, error) {
	//Validate column names
	columnErr := a.checkColumnName(column)
	if columnErr != nil {
		return nil, columnErr
	}
	
	//Get Data
	rows, rowErr := a.db.db.Query("SELECT * FROM Images WHERE `"+column+"` = ?;", identifier)
	if rowErr != nil {
		return nil, rowErr
	}
	defer rows.Close()
	
	//Extract Data to struct object
	var images []*Image
	for rows.Next() {
		image := &Image{}
		scanErr := rows.Scan(&image.ID, &image.Model_ID, &image.Path, &image.Desc); 
		if scanErr != nil {
			return nil, scanErr
		}
		images = append(images, image)
    }
	
	//If errors return nil, if not return struct object
	if completeErr := rows.Err(); completeErr != nil {
		return nil, completeErr
	} else {
		return images, nil
	}
} //Plural

func (a *API) GetImageBy(column string, identifier any) (*Image, error) {
	//Validate column names
	columnErr := a.checkColumnName(column)
	if columnErr != nil {
		return nil, columnErr
	}
	
	//Get Data
	image := &Image{}
	row := a.db.db.QueryRow("SELECT * FROM Images WHERE `"+column+"` = ?;", identifier)
	scanErr := row.Scan(&image.ID, &image.Model_ID, &image.Path, &image.Desc);
	if scanErr != nil {
		return nil, scanErr
	}
	
	return image, nil
} //Singular

func (a *API) AddImage(model_id int, path string, desc string) error {
	statement, prepErr := a.db.db.Prepare("INSERT INTO Images (`model_id`, `path`, `desc`) VALUES (?, ?, ?)")
    if prepErr != nil {
        return prepErr
    }
    _, execErr := statement.Exec(model_id, path, desc)
    return execErr
}

func (a *API) UpdateImage(id int, model_id int, path string, desc string) error {
	//Is this a real image
	oldImage, imageErr := a.GetImageBy("id", id)
	if imageErr != nil {
		return errors.New("Invalid Image ID")
	}
	
	//Missing or blank data means we do not change the field (can think of a more clever way to do this, but idk might be pitfalls)
	var finalModelID, finalPath, finalDesc string
	if model_id == 0 {
		finalModelID = strconv.Itoa(oldImage.Model_ID)
	} else {
		finalModelID = strconv.Itoa(model_id)
	}
	if path == "" {
		finalPath = oldImage.Path
	} else {
		finalPath = path
	}
	if desc == "" {
		finalDesc = oldImage.Desc
	} else {
		finalDesc = desc
	}
	
	//Build Query
	queryString := "UPDATE Images SET `model_id`=?, `path`=?, `desc`=? WHERE `id` = ?;"	
	statement, prepErr := a.db.db.Prepare(queryString)
	 if prepErr != nil {
        return prepErr
    }
    
	 _, execErr := statement.Exec(finalModelID, finalPath, finalDesc, id)
    return execErr
}

func (a *API) DeleteImage(id int) error {
	//Is this a real image
	_, imageErr := a.GetImageBy("id", id)
	if imageErr != nil {
		return errors.New("Invalid Image ID")
	}
	
	statement, prepErr := a.db.db.Prepare("DELETE FROM Images WHERE `id` = ?;")
    if prepErr != nil {
        return prepErr
    }
    
    _, execErr := statement.Exec(id)
    return execErr
}

//Videos
func (a *API) GetVideosBy(column string, identifier any) ([]*Video, error) {
	//Validate column names
	columnErr := a.checkColumnName(column)
	if columnErr != nil {
		return nil, columnErr
	}
	
	//Get Data
	rows, rowErr := a.db.db.Query("SELECT * FROM Videos WHERE `"+column+"` = ?;", identifier)
	if rowErr != nil {
		return nil, rowErr
	}
	defer rows.Close()
	
	//Extract Data to struct object
	var videos []*Video
	for rows.Next() {
		video := &Video{}
		scanErr := rows.Scan(&video.ID, &video.Model_ID, &video.Path, &video.Desc); 
		if scanErr != nil {
			return nil, scanErr
		}
		videos = append(videos, video)
    }
	
	//If errors return nil, if not return struct object
	if completeErr := rows.Err(); completeErr != nil {
		return nil, completeErr
	} else {
		return videos, nil
	}
} //Plural

func (a *API) GetVideoBy(column string, identifier any) (*Video, error) {
	//Validate column names
	columnErr := a.checkColumnName(column)
	if columnErr != nil {
		return nil, columnErr
	}
	
	//Get Data
	video := &Video{}
	row := a.db.db.QueryRow("SELECT * FROM Videos WHERE `"+column+"` = ?;", identifier)
	scanErr := row.Scan(&video.ID, &video.Model_ID, &video.Path, &video.Desc);
	if scanErr != nil {
		return nil, scanErr
	}
	
	return video, nil
} //Singular

func (a *API) AddVideo(model_id int, path string, desc string) error {
	statement, prepErr := a.db.db.Prepare("INSERT INTO Videos (`model_id`, `path`, `desc`) VALUES (?, ?, ?)")
	if prepErr != nil {
    	return prepErr
    }
    
    _, execErr := statement.Exec(model_id, path, desc)
    return execErr
}

func (a *API) UpdateVideo(id int, model_id int, path string, desc string) error {
	//Is this a real video
	oldVideo, videoErr := a.GetVideoBy("id", id)
	if videoErr != nil {
		return errors.New("Invalid Video ID")
	}
	
	//Missing or blank data means we do not change the field (can think of a more clever way to do this, but idk might be pitfalls)
	var finalModelID, finalPath, finalDesc string
	if model_id == 0 {
		finalModelID = strconv.Itoa(oldVideo.Model_ID)
	} else {
		finalModelID = strconv.Itoa(model_id)
	}
	if path == "" {
		finalPath = oldVideo.Path
	} else {
		finalPath = path
	}
	if desc == "" {
		finalDesc = oldVideo.Desc
	} else {
		finalDesc = desc
	}
	
	queryString := "UPDATE Videos SET `model_id`=?, `path`=?, `desc`=? WHERE `id` = ?;"	
	statement, prepErr := a.db.db.Prepare(queryString)
	if prepErr != nil {
    	return prepErr
    }
    
	 _, execErr := statement.Exec(finalModelID, finalPath, finalDesc, id)
    return execErr
}

func (a *API) DeleteVideo(id int) error {
	//Is this a real video
	_, videoErr := a.GetVideoBy("id", id)
	if videoErr != nil {
		return errors.New("Invalid Video ID")
	}
	
	statement, prepErr := a.db.db.Prepare("DELETE FROM Videos WHERE `id` = ?;")
    if prepErr != nil {
    	return prepErr
    }
    
    _, execErr := statement.Exec(id)
    return execErr
}

//Util
func (a *API) checkColumnName(column string) error {
	allowed := []string{"id", "name", "trial", "get", "add", "update", "delete", "desc", "model_id", "path" }
	if slices.Contains(allowed, column) {
		return nil
	} else {
		return errors.New("Disallowed column name, NO SQL INJECTIONS!")
	}
}