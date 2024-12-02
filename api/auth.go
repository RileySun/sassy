package api

import(
	"log"
	"fmt"
	"errors"
	"context"
	
	openbao "github.com/openbao/openbao/api/v2"
)

//Struct
type Auth struct {
	baoClient *openbao.Client
}

//Create
func NewAuth() *Auth {
	auth := &Auth{}
	
	//Credentials
	creds := LoadCredentials()	
	
	//Openbao
	config := openbao.DefaultConfig()
	config.Address = creds.OHost
	var err error
	auth.baoClient, err = openbao.NewClient(config)
	if err != nil {
	    log.Fatalf("OpenBao Client Init Error: %v", err)
	}
	auth.baoClient.SetToken(creds.OToken)
	
	return auth
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
	
	value, ok := secret.Data["password"].(string)
	if !ok {
		err = errors.New(fmt.Sprintf("Type Assert Error: %T %#v", secret.Data["password"], secret.Data["password"]))
		return "", err
	}
	
	return value, nil
}