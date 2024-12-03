package api

import(
	"testing"
	"reflect"
)

func TestSecretManagerConnection(t *testing.T) {
	_, err := NewAuth()
	if err != nil {
		t.Error("Cannot connect to secret manager")
	}
}

func TestNewToken(t *testing.T) {
	token := NewToken()
	
	if reflect.TypeOf(token) != reflect.TypeOf("string") {
		t.Errorf("Token is not string. %T, wanted string", token)
	}
	
	if len(token) != 36 {
		t.Errorf("Token character length. %v, wanted 36", len(token))
	}
}

func TestAddGetSecret(t *testing.T) {
	auth, err := NewAuth()
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