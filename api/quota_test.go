package api

import(
	"testing"
)

func TestIsUnderQuota(t *testing.T) {
	//These users correspond to real users in the database, one being unlimited rated
	//and the other one already being at max rate limit
	api := NewAPI()
	unlimited := api.GetUserBy("name", "unlimited")
	limited := api.GetUserBy("name", "limited")
	
	//Unlimited Rate User
	nonTrial := api.IsUnderQuota(unlimited.ID, "Update")
	if !nonTrial {
		t.Error("User should be not trial limited")
	}
	
	//Rate Limited User
	underGetQuota := api.IsUnderQuota(limited.ID, "Get")
	if underGetQuota {
		t.Error("User should be at max limit for Get funtion")
	}
	
	underAddQuota := api.IsUnderQuota(limited.ID, "Add")
	if underAddQuota {
		t.Error("User should be at max limit for Add funtion")
	}
	
	underUpdateQuota := api.IsUnderQuota(limited.ID, "Update")
	if underUpdateQuota {
		t.Error("User should be at max limit for Update funtion")
	}
	
	underDeleteQuota := api.IsUnderQuota(limited.ID, "Delete")
	if underDeleteQuota {
		t.Error("User should be at max limit for Delete funtion")
	}
}