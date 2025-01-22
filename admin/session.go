package admin

import(
	"time"
	"net/http"
	"github.com/julienschmidt/httprouter"
)

//Active user sessions (we dont need a database, there will never be more than 10)
var sessions = map[string]session{}

type session struct {
	username string
	expiry   time.Time
}

func (a *Admin) Login(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	
}