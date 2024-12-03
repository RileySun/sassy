package api

import(

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

func (a *API) AddToQuota(userID string, quotaType string) {
	
}