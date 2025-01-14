package api

import(
	"errors"
	"slices"
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
}

func (a *API) AddModel(name string, desc string) error {
	statement, prepErr := a.db.db.Prepare("INSERT INTO Models (`name`, `desc`) VALUES (?, ?)")
    if prepErr != nil {
        return prepErr
    }
    _, execErr := statement.Exec(name, desc)
    return execErr
}

func (a *API) UpdateModel(id int, name string, desc string) error {
	queryString := "UPDATE Models SET `name`=?, `desc`=? WHERE `id` = ?;"	
	statement, prepErr := a.db.db.Prepare(queryString)
	 if prepErr != nil {
       return prepErr
    }
	 _, execErr := statement.Exec(name, desc, id)
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
}

func (a *API) AddImage(model_id int, path string, desc string) error {
	statement, prepErr := a.db.db.Prepare("INSERT INTO Images (`model_id`, `path`, `desc`) VALUES (?, ?, ?)")
    if prepErr != nil {
        return prepErr
    }
    _, execErr := statement.Exec(model_id, path, desc)
    return execErr
}

func (a *API) UpdateImage(id int, model_id int, path string, desc string) error {
	queryString := "UPDATE Images SET `model_id`=?, `path`=?, `desc`=? WHERE `id` = ?;"	
	statement, prepErr := a.db.db.Prepare(queryString)
	 if prepErr != nil {
        return prepErr
    }
    
	 _, execErr := statement.Exec(model_id, path, desc, id)
    return execErr
}

func (a *API) DeleteImage(id int) error {
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
}

func (a *API) AddVideo(model_id int, path string, desc string) error {
	statement, prepErr := a.db.db.Prepare("INSERT INTO Videos (`model_id`, `path`, `desc`) VALUES (?, ?, ?)")
	if prepErr != nil {
    	return prepErr
    }
    
    _, execErr := statement.Exec(model_id, path, desc)
    return execErr
}

func (a *API) UpdateVideo(id int, model_id int, path string, desc string) error {
	queryString := "UPDATE Videos SET `model_id`=?, `path`=?, `desc`=? WHERE `id` = ?;"	
	statement, prepErr := a.db.db.Prepare(queryString)
	if prepErr != nil {
    	return prepErr
    }
    
	 _, execErr := statement.Exec(model_id, path, desc, id)
    return execErr
}

func (a *API) DeleteVideo(id int) error {
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