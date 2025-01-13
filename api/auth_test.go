package api

import(
	"os"
	//"log"
	"testing"
)

var auth *Auth
var token *Token

func TestMain(m *testing.M) {
	auth = NewAuth()
	exitCode := m.Run()
	os.Exit(exitCode)
}

func TestNewUser(t *testing.T) {
	createErr := auth.NewUser("Tester", "TeStEd", false)
	if createErr != nil {
		t.Error(createErr.Error())
	}
}

func TestNewToken(t *testing.T) {
	newToken, tokenErr := auth.NewToken("Tester", "TeStEd")
	if tokenErr != nil {
		t.Error(tokenErr.Error())
	}
	token = newToken
}

func TestDeleteToken(t *testing.T) {
	deleteErr := auth.DeleteToken(token.refresh)
	if deleteErr != nil {
		t.Error(deleteErr.Error())
	}
}

func TestCheckLogin(t *testing.T) {
	loginErr := auth.CheckLogin("Tester", "TeStEd")
	if loginErr != nil {
		t.Error(loginErr.Error())
	}
}

func TestDeleteUser(t *testing.T) {
	deleteErr := auth.DeleteUser("Tester", "TeStEd")
	if deleteErr != nil {
		t.Error(deleteErr.Error())
	}
}

//Secrets
func TestSecretManagerConnection(t *testing.T) {	
	err := auth.NewSecretManager()
	if err != nil {
		t.Error(err.Error())
	}
}

func TestAddGetSecret(t *testing.T) {
	err := auth.NewSecretManager()
	if err != nil {
		t.Error("Cannot connect to secret manager")
		t.Fail()
	}
	
    janeData := map[string]interface{}{"Key": "73b1d07b-d4b8-4976-96cd-9b76b99e45b1",}
	janeToken := "8fb22136-0571-43e1-9c21-220ecc59821f"
	auth.AddSecret(janeData, janeToken)
	
	storedSecret, err := auth.GetSecret("8fb22136-0571-43e1-9c21-220ecc59821f")
	
	if err != nil {
		t.Error("Error Retrieving Stored Secret")
	}
	
	if storedSecret != "73b1d07b-d4b8-4976-96cd-9b76b99e45b1" {
		t.Errorf("Secrets dont match. %s, wanted 73b1d07b-d4b8-4976-96cd-9b76b99e45b1", storedSecret)
	}
}