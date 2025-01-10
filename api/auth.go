package api

import(
	"fmt"
	"errors"
	"context"
	
	"github.com/google/uuid"
	
	"golang.org/x/crypto/bcrypt"
	
	openbao "github.com/openbao/openbao/api/v2"
)

//Struct
type Auth struct {
	baoClient *openbao.Client
}

//Create
func NewAuth() (*Auth) {
	auth := &Auth{}
	
	
	
	return auth
}

func NewToken() string {
	id := uuid.New()
	return id.String()
}

//User Functions
func (a *Auth) NewUser(name string, password string, trial bool) error {
	db := NewDB()
	defer db.db.Close()//Close this one
	
	//Bcrypt pass
	pass, passErr := a.HashPassword(password)
	if passErr != nil {
		return passErr
	}
	
	statement, prepErr := db.db.Prepare("INSERT INTO Users (`name`, `pass`, `key`, `trial`) VALUES (?, ?, ?, ?)")
	if prepErr != nil {
		return prepErr
	}
	
	_, stateErr := statement.Exec(name, pass, NewToken(), trial)
	return stateErr
}

func (a *Auth) DeleteUser(name string, password string) error {
	//CheckPassword, is this right user?
	if a.CheckLogin(name, password) != nil {
		return errors.New("No such user, or password incorrect")
	} 

	//Now Delete
	db := NewDB()
	defer db.db.Close()//Close this one
	
	statement, prepErr := db.db.Prepare("DELETE FROM Users WHERE `name` = ?")
	if prepErr != nil {
		return prepErr
	}
	
	_, stateErr := statement.Exec(name)
	return stateErr
} //Might Delete user if usernames are the same and index is first, maybe we check in user creation so no users with same name...

func (a *Auth) CheckLogin(name string, password string) error {
	db := NewDB()
	defer db.db.Close()//Close this one, we wont need it later.
	
	//Get User by Name
	rows, rowErr := db.db.Query("SELECT `name`, `pass` FROM Users WHERE `name` = ?;", name)
	if rowErr != nil {
		return rowErr
	}
	defer rows.Close()
	
	//Don't use API user, we dont need all that info and API doesnt need passwords
	user := struct {
		Name string 
		Pass string
	}{}
	for rows.Next() {
		scanErr := rows.Scan(&user.Name, &user.Pass); 
		if scanErr != nil {
			return scanErr
		}
    }
	
	//Completion Errors
	if completeErr := rows.Err(); completeErr != nil {
		return completeErr
	}
	
	//Check for correct password
	if a.CheckPassword(user.Pass, password) {
		return nil
	} else {
		return errors.New("Invalid Password")
	}
}

func (a *Auth) HashPassword(password string) (string, error) {
	bytes := []byte(password)
	hasedBytes, err := bcrypt.GenerateFromPassword(bytes, bcrypt.MinCost)
	return string(hasedBytes), err
}

func (a *Auth) CheckPassword(hashed string, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashed), []byte(password))
	return err == nil
}

//Secret Storage (using openbau)
func (a *Auth) NewSecretManager() error {
	//Credentials
	creds := LoadCredentials()	
	
	//Openbao
	config := openbao.DefaultConfig()
	config.Address = creds.OHost
	
	var baoErr error
	a.baoClient, baoErr = openbao.NewClient(config)
	if baoErr != nil {
	    return baoErr
	}
	
	a.baoClient.SetToken(creds.OToken)
	
	return nil
}

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