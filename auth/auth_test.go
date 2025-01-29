package auth

import(
	"testing"
)

var auth *Auth
var token *Token

//Main
func AuthTestInit() {
	auth = NewAuth()
}

//User (Delete when finished)
func TestNewUser(t *testing.T) {
	createErr := auth.NewUser("Tester", "TeStEd", false)
	if createErr != nil {
		t.Error(createErr.Error())
		t.Fail()
	}
}

func TestIsUserTaken(t *testing.T) {
	takenErr := auth.IsUserTaken("Tester")
	if takenErr == nil {
		t.Error("Username should be taken")
	}
}

func TestCheckLogin(t *testing.T) {
	loginErr := auth.CheckLogin("Tester", "TeStEd")
	if loginErr != nil {
		t.Error(loginErr.Error())
	}
}

//Token
func TestNewToken(t *testing.T) {
	newToken, tokenErr := auth.NewToken("Tester", "TeStEd")
	if tokenErr != nil {
		t.Error(tokenErr.Error())
		t.Fail()
	}
	token = newToken
}

func TestCheckToken(t *testing.T) {
	isValid := auth.CheckToken(token.access)
	if isValid != nil {
		t.Error(isValid.Error())
	}
}

func TestGetToken(t *testing.T) {
	newToken, tokenErr := auth.GetToken("refresh_token", token.refresh)
	if tokenErr != nil {
		t.Error(tokenErr.Error())
	}
	
	if newToken.refresh != token.refresh || newToken.access != token.access {
		t.Error("Tokens do not match")
	}
}

func TestGenerateToken(t *testing.T) {
	newAccess, genErr := auth.GenerateToken(token.refresh)
	if genErr != nil {
		t.Error(genErr.Error())
	}
	
	//Refresh test copy of token
	var tokenErr error
	token, tokenErr = auth.GetToken("refresh_token", token.refresh)
	if tokenErr != nil {
		t.Error(tokenErr.Error())
	}
	
	//Does it match
	if newAccess != token.access {
		t.Error("Access token has not changed")
	}
}

func TestDeleteToken(t *testing.T) {
	deleteErr := auth.DeleteToken(token.refresh)
	if deleteErr != nil {
		t.Error(deleteErr.Error())
	}
}

//Delete User
func TestDeleteUser(t *testing.T) {
	deleteErr := auth.DeleteUser("Tester", "TeStEd")
	if deleteErr != nil {
		t.Error(deleteErr.Error())
	}
}

//Column Name Sanitization
func TestAuthCheckColumnName(t *testing.T) {
	shouldPass := api.checkColumnName("id")
	if shouldPass != nil {
		t.Error(shouldPass.Error())
	}
	
	shouldFail := api.checkColumnName("test")
	if shouldFail == nil {
		t.Error("Should not be a valid column name")
	}
}

/* See auth.go secrets section
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
*/