package main

import(
	"os"
	"testing"
	"net/http"
)

var server *Server

//Main
func TestMain(m *testing.M) {	
	//This one
	server = NewServer()
	exitCode := m.Run()
	os.Exit(exitCode)
}

func TestCheckAuthentication(t *testing.T) {
	testUser, testPass := "tester", "TeStEd"
	
	//Create test User
	userErr := server.Auth.NewUser(testUser, testPass, false)
	if userErr != nil {
		t.Error(userErr.Error())
		t.Fail()
	}
	
	//Create test Token
	testToken, tokenErr := server.Auth.NewToken(testUser, testPass)
	if tokenErr != nil {
		t.Error(tokenErr.Error())
		t.Fail()
	}
	
	//Create test request
	req, reqErr := http.NewRequest(http.MethodGet, "localhost:8080", nil)
	if reqErr != nil {
		t.Error(reqErr.Error())
		t.Fail()
	}
	
	//Test Authorization
	req.Header["Authorization"] = []string{"OAuth "+ testToken.Access()}
	_, authErr := server.CheckAuthentication(req)
	if authErr != nil {
		t.Error(authErr.Error())
		t.Fail()
	}
	
	//Test Invalid Authorization
	req.Header["Authorization"] = []string{"OAuth 823893298dh993828hd82h9s89892832"}
	_, invalidAuthErr := server.CheckAuthentication(req)
	if invalidAuthErr == nil {
		t.Error("Authorization Should Be Invalid")
		t.Fail()
	}
	
	//Cleanup (delete, test user and token)
	deleteTokenErr := server.Auth.DeleteToken(testToken.Refresh())
	deleteUserErr := server.Auth.DeleteUser(testUser, testPass)
	if deleteTokenErr != nil || deleteUserErr != nil {
		t.Log(deleteTokenErr, deleteUserErr)
	}
}