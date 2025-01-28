package main

import(
	
	"admin"
)

func main() {
	newAdmin := admin.NewAdmin()
	
	newAdmin.ApiAction = func(actionType string) {
		if actionType == "Shutdown" {
			newAdmin.Shutdown()
		} else {
			newAdmin.Restart()
		}
	}
	
	newAdmin.LaunchServer()
	
	for true {
		
    }
}