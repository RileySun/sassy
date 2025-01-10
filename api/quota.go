package api

import(
	"log"
	"strconv"
)

const TRIAL_GET_LIMIT, TRIAL_ADD_LIMIT, TRIAL_UPDATE_LIMIT, TRIAL_DELETE_LIMIT = 10, 10, 10, 10

func (a *API) IsUnderQuota(userID int, actionType string) bool {
	user := a.GetUserBy("`id`", userID)
	
	ok := true
	if user.Trial != true {
		return ok
	} //Rate Unlimited
	
	//Rate Limited
	switch actionType {
		case "Get":
			ok = user.Get < TRIAL_GET_LIMIT
		case "Add":
			ok = user.Add < TRIAL_ADD_LIMIT
		case "Update":
			ok = user.Update < TRIAL_UPDATE_LIMIT
		case "Delete":
			ok = user.Delete < TRIAL_DELETE_LIMIT
	}
	
	return ok
}

func (a *API) AddToQuota(userID int, quotaType string) {
	user := a.GetUserBy("`id`", userID)
	
	switch quotaType {
		case "Get":
			user.Get += 1
		case "Add":
			user.Add += 1
		case "Update":
			user.Update += 1
		case "Delete":
			user.Delete += 1
	}
	
	g, ad, u, d := strconv.Itoa(user.Get), strconv.Itoa(user.Add), strconv.Itoa(user.Update), strconv.Itoa(user.Delete) 
	
	a.UpdateUserQuotas(userID, g, ad, u, d)
}

func (a *API) UpdateUserQuotas(id int, getQuota string, addQuota string, updateQuota string, deleteQuota string) error {
	queryString := "UPDATE Users SET `get`=?, `add`=?, `update`=?, `delete`=? WHERE `id` = ?;"	
	statement, err := a.db.db.Prepare(queryString)
	 if err != nil {
        log.Println(err, "- API:UpdateUser")
    }
	 _, err = statement.Exec(id, getQuota, addQuota, updateQuota, deleteQuota)
    return err
}