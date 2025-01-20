package api

import(
	"errors"
	"strconv"
)

const TRIAL_GET_LIMIT, TRIAL_ADD_LIMIT, TRIAL_UPDATE_LIMIT, TRIAL_DELETE_LIMIT = 10, 10, 10, 10

func (a *API) IsUnderQuota(userID int, actionType string) error {
	user, getErr := a.GetUserBy("id", userID)
	if getErr != nil {
		return getErr
	}
	
	//Rate Unlimited
	if user.Trial != true {
		return nil
	}
	
	//Rate Limited
	switch actionType {
		case "Get":
			if user.Get > TRIAL_GET_LIMIT {
				return errors.New("User is rate limited for usage - Get")
			}
		case "Add":
			if user.Add > TRIAL_ADD_LIMIT {
				return errors.New("User is rate limited for usage - Add")
			}
		case "Update":
			if user.Update > TRIAL_UPDATE_LIMIT {
				return errors.New("User is rate limited for usage - Update")
			}
		case "Delete":
			if user.Delete > TRIAL_DELETE_LIMIT {
				return errors.New("User is rate limited for usage - Delete")
			}
		default:
			return errors.New("Invalid Method: " + actionType)
	}
	
	return nil
}

func (a *API) AddToQuota(userID int, quotaType string) error {
	user, getErr := a.GetUserBy("id", userID)
	if getErr != nil {
		return getErr
	}
	
	switch quotaType {
		case "Get":
			user.Get += 1
		case "Add":
			user.Add += 1
		case "Update":
			user.Update += 1
		case "Delete":
			user.Delete += 1
		default:
			return errors.New("Invalid Quota Method:" + quotaType)
	}
	
	g, ad, u, d := strconv.Itoa(user.Get), strconv.Itoa(user.Add), strconv.Itoa(user.Update), strconv.Itoa(user.Delete) 
	
	return a.UpdateUserQuotas(userID, g, ad, u, d)
}

func (a *API) UpdateUserQuotas(id int, getQuota string, addQuota string, updateQuota string, deleteQuota string) error {
	queryString := "UPDATE Users SET `get`= ?, `add`= ?, `update`= ?, `delete`= ? WHERE `id` = ?;"	
	statement, stmtErr := a.db.db.Prepare(queryString)
	 if stmtErr != nil {
        return stmtErr
    }
	 _, execErr := statement.Exec(getQuota, addQuota, updateQuota, deleteQuota, id)
    return execErr
}