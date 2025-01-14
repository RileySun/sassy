package api

import(
	"testing"
)

func TestIsUnderQuota(t *testing.T) {
	unlimited, unlimitedErr := api.GetUserBy("name", "unlimited")
	if unlimitedErr != nil {
		t.Error(unlimitedErr.Error())
	}
	
	limited, limitedErr := api.GetUserBy("name", "limited")
	if limitedErr != nil {
		t.Error(limitedErr.Error())
	}
	
	//Unlimited Rate User
	nonTrial := api.IsUnderQuota(unlimited.ID, "Update")
	if nonTrial != nil  {
		t.Error("User `unlimited` should be not trial limited")
	}
	
	//Rate Limited User
	underGetQuota := api.IsUnderQuota(limited.ID, "Get")
	if underGetQuota == nil {
		t.Error("User `limited` should be at max limit for Get funtion")
	}
	
	underAddQuota := api.IsUnderQuota(limited.ID, "Add")
	if underAddQuota == nil {
		t.Error("User should be at max limit for Add funtion")
	}
	
	underUpdateQuota := api.IsUnderQuota(limited.ID, "Update")
	if underUpdateQuota == nil {
		t.Error("User should be at max limit for Update funtion")
	}
	
	underDeleteQuota := api.IsUnderQuota(limited.ID, "Delete")
	if underDeleteQuota == nil {
		t.Error("User should be at max limit for Delete funtion")
	}
}