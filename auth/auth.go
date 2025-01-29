package auth

import(
	"log"
	"time"
	"errors"
	"slices"
	"strings"
	
	"github.com/google/uuid"
	
	"golang.org/x/crypto/bcrypt"
	
	//"fmt"
	//"context"
	//openbao "github.com/openbao/openbao/api/v2"
)

//Struct
type Auth struct {
	db *Database
	//baoClient *openbao.Client
}

type Token struct {
	access string
	refresh string
	expires time.Time
}

//Create
func NewAuth() (*Auth) {
	auth := &Auth{
		db:NewDB(),
	}

	return auth
}

//Token Functions
func (a *Auth) NewToken(name string, password string) (*Token, error) {
	loginErr := a.CheckLogin(name, password)
	if loginErr != nil {
		return nil, loginErr
	}

	rawAccess := uuid.New().String()
	rawRefresh := uuid.New().String()
	newToken := &Token{
		access: strings.Replace(rawAccess, "-", "", -1),
		refresh: strings.Replace(rawRefresh, "-", "", -1),
		expires: time.Now().Add(time.Second * time.Duration(86400)),
	}
	
	storeErr := a.StoreToken(name, newToken)
	if storeErr != nil {
		return nil, storeErr
	}
	
	return newToken, nil
}

func (a *Auth) StoreToken(name string, newToken *Token) error {	
	statement, prepErr := a.db.db.Prepare("INSERT INTO Tokens (`name`, `access_token`, `refresh_token`, `expires_in`) VALUES (?, ?, ?, ?)")
	if prepErr != nil {
		return prepErr
	}
	
	_, stateErr := statement.Exec(name, newToken.access, newToken.refresh, newToken.expires)
	return stateErr
}

func (a *Auth) CheckToken(accessToken string) error {
	//Get Expiration data
	var expires time.Time
	row := a.db.db.QueryRow("SELECT `expires_in` FROM Tokens WHERE `access_token` = ?;", accessToken)
	scanErr := row.Scan(&expires)
	if scanErr != nil {
		return errors.New("Invalid access token")
	}
	
	//Parse Expiration]
	isValid := time.Now().Before(expires)
	if isValid {
		return nil
	} else {
		return errors.New("Access token has expired, please generate a new token.")
	}
} //Selected by access token

func (a *Auth) GetToken(column string, identifier any) (*Token, error) {
	//Validate column names
	columnErr := a.checkColumnName(column)
	if columnErr != nil {
		return nil, columnErr
	}
	
	//Get Data
	token := &Token{}
	row := a.db.db.QueryRow("SELECT `access_token`, `refresh_token`, `expires_in` FROM Tokens WHERE `"+column+"` = ?;", identifier)
	scanErr := row.Scan(&token.access, &token.refresh, &token.expires)
	if scanErr != nil {
		return nil, scanErr
	}
	return token, nil
}

func (a *Auth) GenerateToken(refreshToken string) (string, error) {	
	//Create New Access Token and Expiration
	rawAccess := uuid.New().String()
	access := strings.Replace(rawAccess, "-", "", -1)
	expires := time.Now().Add(time.Second * time.Duration(86400))
	
	//Store new token
	statement, prepErr := a.db.db.Prepare("Update Tokens  SET `access_token` = ?, `expires_in` = ? WHERE `refresh_token` = ?")
	if prepErr != nil {
		return "", prepErr
	}
	
	_, stateErr := statement.Exec(access, expires, refreshToken)
	if stateErr != nil {
		return "", nil
	}
	
	return access, nil
} //Generate new access token

func (a *Auth) DeleteToken(refreshToken string) error {
	//Now Delete
	db := NewDB()
	defer db.db.Close()//Close this one
	
	statement, prepErr := a.db.db.Prepare("DELETE FROM Tokens WHERE `refresh_token` = ?")
	if prepErr != nil {
		return prepErr
	}
	
	_, stateErr := statement.Exec(refreshToken)
	return stateErr
} //Selected by refresh token

func (a *Auth) GetUserIdFromToken(accessToken string) int {
	name := ""
	row := a.db.db.QueryRow("SELECT `name` FROM Tokens WHERE `access_token` = ?;", accessToken)
	scanErr := row.Scan(&name)
	if scanErr != nil {
		log.Println("Error getting user ID from refresh token. NAME")
	}
	
	userID := 0
	row = a.db.db.QueryRow("SELECT `id` FROM Users WHERE `name` = ?;", name)
	scanErr = row.Scan(&userID)
	if scanErr != nil {
		log.Println("Error getting user ID from refresh token. ID")
	}
	
	return userID
}

func (t *Token) Access() string {
	return t.access
} //Returns access token

func (t *Token) Refresh() string {
	return t.refresh
} //returns RefreshToken

//User Functions
func (a *Auth) NewUser(name string, password string, trial bool) error {
	isTaken := a.IsUserTaken(name)
	if isTaken != nil {
		return isTaken
	}

	//Bcrypt pass
	pass, passErr := a.HashPassword(password)
	if passErr != nil {
		return passErr
	}
	
	statement, prepErr := a.db.db.Prepare("INSERT INTO Users (`name`, `pass`, `trial`) VALUES (?, ?, ?)")
	if prepErr != nil {
		return prepErr
	}
	
	_, stateErr := statement.Exec(name, pass, trial)
	if stateErr != nil {
		return stateErr
	}
	
	return stateErr
}

func (a *Auth) IsUserTaken(name string) error {
	//Get Data
	rows, rowErr := a.db.db.Query("SELECT `name` FROM Users")
	if rowErr != nil {
		return rowErr
	}
	defer rows.Close()
	
	//Extract Names
	for rows.Next() {
		var taken string
		scanErr := rows.Scan(&taken); 
		if scanErr != nil {
			return scanErr
		}
		
		//Check if username is taken
		if taken == name {
			return errors.New("Username is taken, please enter a different username")
		}
    }
    
    return nil
}

func (a *Auth) DeleteUser(name string, password string) error {
	//Now Delete
	db := NewDB()
	defer db.db.Close()//Close this one
	
	statement, prepErr := a.db.db.Prepare("DELETE FROM Users WHERE `name` = ?")
	if prepErr != nil {
		return prepErr
	}
	
	_, stateErr := statement.Exec(name)
	return stateErr
}

func (a *Auth) CheckLogin(name string, password string) error {
	//Get Data
	user := struct {
		Name string 
		Pass string
	}{}
	row := a.db.db.QueryRow("SELECT `name`, `pass` FROM Users WHERE `name` = ?;", name)
	scanErr := row.Scan(&user.Name, &user.Pass)
	if scanErr != nil {
		return errors.New("Invalid Login Information")
	}
	
	//Check for correct password
	if a.CheckPassword(user.Pass, password) {
		return nil
	} else {
		return errors.New("Invalid Login Information")
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

//Util
func (a *Auth) checkColumnName(column string) error {
	allowed := []string{"id", "name", "refresh_token", "expires_in"}
	if slices.Contains(allowed, column) {
		return nil
	} else {
		return errors.New("Disallowed column name, NO SQL INJECTIONS!")
	}
}

/* Not sure what we are using this for yet.
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
*/