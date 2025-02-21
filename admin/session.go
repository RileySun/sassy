package admin

import(
	"os"
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
	expiration := time.Now().Add(1200 * time.Second)
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
	expired := s.expires.Before(time.Now())
	return expired
}

func (a *Admin) removeSession(token string) {
	for _, s := range a.Sessions {
		if s == a.Sessions[token] {
			delete(a.Sessions, token);
			return
		}
	}
}

//LoadLogin
func (a *Admin) LoadLogin(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	tmpl, parseErr := template.ParseFS(HTMLFiles, "html/login.html")
	if parseErr != nil {
		log.Println("Template->LoadPage->Parse: ", parseErr)
	}
	
	//Show Popup
	tmpl.Execute(w, false)
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
		a.removeSession(token)
		return errors.New("Expired Session Token")
	} else {
		//Add extra time due to success check
		a.Sessions[token].expires = time.Now().Add(1200 * time.Second)
	}
	
	return nil
}

//Execute Login
func (a *Admin) Login(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	r.ParseForm()
	user := r.PostFormValue("user")
	pass := r.PostFormValue("pass")
	
	//Fake it for now, we gonna hook an auth db to this
	adminUser := os.Getenv("ADMIN_USER")
	adminPass := os.Getenv("ADMIN_PASS")
	log.Println(adminUser, adminPass)
	if user == adminUser && pass == adminPass {
		a.NewSession(w, user)
		
		http.Redirect(w, r, "/" + a.Redirect, http.StatusFound)		
		a.Redirect = ""
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

//Logout
func (a *Admin) Logout(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	//If session, delete session
	c, cookieErr := r.Cookie("session_token")
	if cookieErr == nil {
		token := c.Value
		a.removeSession(token)
	}
	
	http.Redirect(w, r, "/login" + a.Redirect, http.StatusFound)
}