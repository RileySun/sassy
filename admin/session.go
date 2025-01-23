package admin

import(
	"log"
	"time"
	"errors"
	"net/http"
	"html/template"
	
	"github.com/google/uuid"
	"github.com/julienschmidt/httprouter"
)

//Active user sessions (we dont need a database, there will never be more than 10)
type Session struct {
	user string
	expires time.Time
}

func (a *Admin) NewSession(w http.ResponseWriter, user string) {
	token := uuid.NewString()
	expiration := time.Now().Add(120 * time.Second)
	a.Sessions[token] = &Session{
		user: user,
		expires: expiration,
	}
	
	http.SetCookie(w, &http.Cookie{
		Name: "session_token",
		Value: token,
		Expires: expiration,
	})
}

func (s *Session) isExpired() bool {
	return s.expires.Before(time.Now())
}

//CheckLogin
func (a Admin) CheckSession(r *http.Request) error {
	c, cookieErr := r.Cookie("session_token")
	if cookieErr != nil {
		return errors.New("Invalid Session Cookie")
	}
	
	token := c.Value
	session, exists := a.Sessions[token]
	if !exists {
		return errors.New("Invalid Session Token")
	}
	
	if session.isExpired() {
		return errors.New("Expired Session Token")
	}
	
	return nil
}

//Execute Login
func (a *Admin) Login(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	r.ParseForm()
	user := r.PostFormValue("user")
	pass := r.PostFormValue("pass")
	
	//Fake it for now, we gonna hook an auth db to this
	if user == "janeD" && pass == "prW4nj7KL" {
		a.NewSession(w, user)
		http.Redirect(w, r, "/", http.StatusFound)		
	} else {
		//Reload Login Page
		tmpl, parseErr := template.ParseFS(HTMLFiles, "html/login.html")
		if parseErr != nil {
			log.Println("Session->Login->Error: ", parseErr)
		}
		
		//Show Popup
		tmpl.Execute(w, true)
	}
}