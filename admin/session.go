package admin

import(
	"io"
	"os"
	"log"
	"time"
	"errors"
	"strconv"
	"strings"
	"net/http"
	"html/template"
	
	URL "net/url"
	
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
	
	//First check admin root user/password
	adminUser := os.Getenv("ADMIN_USER")
	adminPass := os.Getenv("ADMIN_PASS")
	if user == adminUser && pass == adminPass {
		a.NewSession(w, user)
		
		http.Redirect(w, r, "/" + a.Redirect, http.StatusFound)		
		a.Redirect = ""
		return
	}
	
	//Then check auth server DB user/password
	loginErr := a.CheckAdminLogin(user, pass)
	if loginErr == nil {
		a.NewSession(w, user)
		
		http.Redirect(w, r, "/" + a.Redirect, http.StatusFound)		
		a.Redirect = ""
		return
	}
	
	//If none of those work Login failed, Reload Login Page
	tmpl, parseErr := template.ParseFS(HTMLFiles, "html/login.html")
	if parseErr != nil {
		log.Println("Session->Login->Error: ", parseErr)
	}
	
	//Show Popup
	tmpl.Execute(w, true)
}

func (a *Admin) CheckAdminLogin(user, pass string) error {
	//Create Request
	client := &http.Client{}
	url := GetServerURL("Auth") + "/admin"
	req, reqErr := http.NewRequest("POST", url, nil)
	if reqErr != nil {
		log.Println("Admin->CheckAdminLogin->Req: ", reqErr)
		return reqErr
	}
	
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	
	form := URL.Values{}
	form.Add("user", user)
	form.Add("pass", pass)
	req.Body = io.NopCloser(strings.NewReader(form.Encode()))
	
	//Do Request
	resp, getErr := client.Do(req)
	defer resp.Body.Close()
	if getErr != nil {
		log.Println("Admin->CheckAdminLogin->Get: ", getErr)
		return getErr
	}
	
	//Read Body
	body, bodyErr := io.ReadAll(resp.Body)
	if bodyErr != nil {
		log.Println("Admin->CheckAdminLogin->Body: ", bodyErr)
		return bodyErr
	}
	
	//Parse User ID
	answer, parseErr := strconv.Atoi(string(body[:]))
	if parseErr != nil {
		log.Println("Admin->CheckAdminLogin->Parse", parseErr)
		return parseErr
	}
	
	if answer == 1 {
		return nil
	} else {
		return errors.New("Invalid Auth DB User/Pass")
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