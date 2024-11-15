package api

import(
	//"fmt"
	"crypto/sha256"
	
	"github.com/google/uuid"
)

//Struct
type Authenticator struct {
	reverseKeyIndex map[string]string
}

//Create
func NewAuth() *Authenticator {
	auth := &Authenticator{}
	
	keyStorage := LoadKeys()
	auth.reverseKeyIndex = make(map[string]string)
	for name, key := range keyStorage {
       auth.reverseKeyIndex[key] = name
    }
	
	return auth
}

//Private

//Public
func (a *Authenticator) GenerateToken() string {
	raw := uuid.New().String()
	hash := sha256.Sum256([]byte(raw))
	key := string(hash[:])
	
	return key
}


func (a *Authenticator) IsValidKey(newKey string) (string, bool) {
	hash := sha256.Sum256([]byte(newKey))
    key := string(hash[:])
    
    name, found := a.reverseKeyIndex[key]
    
    return name, found
}