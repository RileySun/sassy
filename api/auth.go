package api

import(
	"fmt"
	"errors"
	"context"
	
	"github.com/google/uuid"
	
	openbao "github.com/openbao/openbao/api/v2"
)

//Struct
type Auth struct {
	baoClient *openbao.Client
}

//Create
func NewAuth() (*Auth, error) {
	auth := &Auth{}
	
	//Credentials
	creds := LoadCredentials()	
	
	//Openbao
	config := openbao.DefaultConfig()
	config.Address = creds.OHost
	
	var err error
	auth.baoClient, err = openbao.NewClient(config)
	if err != nil {
	    return nil, err
	}
	
	auth.baoClient.SetToken(creds.OToken)
	
	return auth, nil
}

func NewToken() string {
	id := uuid.New()
	return id.String()
}


//Secret Storage (using openbau)
func (a *Auth) AddSecret(secretData map[string]interface{}, secretPassword string) error {
	_, err := a.baoClient.KVv2("secret").Put(context.Background(), secretPassword, secretData)
	return err
}

func (a *Auth) GetSecret(secretPassword string) (string, error) {
	secret, err := a.baoClient.KVv2("secret").Get(context.Background(), secretPassword)
	if err != nil {
    	return "", err
	}
	
	value, ok := secret.Data["Key"].(string)
	if !ok {
		err = errors.New(fmt.Sprintf("Type Assert Error: %T %#v", secret.Data["Key"], secret.Data["Key"]))
		return "", err
	}
	
	return value, nil
}