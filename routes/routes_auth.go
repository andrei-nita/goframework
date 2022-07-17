package routes

import (
	fk "github.com/andrei-nita/goframework/framework"
	"github.com/andrei-nita/goframework/framework/auth"
	"github.com/gorilla/csrf"
	"github.com/xyproto/simplebolt"
	"log"
	"net/http"
)

const (
	cookieErrLogin  = "errLogin"
	cookieErrSignup = "errSignup"
)

// GET /login
// Show the login page
func login(w http.ResponseWriter, r *http.Request) {
	data, _, next := methodCsrfFlashAuth(w, r, http.MethodGet, cookieErrLogin)
	if next {
		tmplExecute(w, r, "login", data)
	}
}

// POST /authenticate
// the user given the email and password
func authenticate(w http.ResponseWriter, r *http.Request) {
	if isAllowed := allowMethod(w, r, http.MethodPost); !isAllowed {
		return
	}

	err := r.ParseForm()

	user, err := auth.UserByEmail(r.PostFormValue("email"))
	if err != nil {
		if err == simplebolt.ErrKeyNotFound {
			fk.CookieSetFlash(cookieErrLogin, w, "Invalid email or password")
			http.Redirect(w, r, "/login", 302)
			return
		}
		log.Println(err, "problem getting the user")
		return
	}

	if user.Password == auth.HashPsw(r.PostFormValue("password")) {
		session, err := user.CreateSession()
		if err != nil {
			log.Println(err, "cannot create session")
			return
		}
		fk.CookieSetSession("_cookie", w, session.Uuid)
		http.Redirect(w, r, "/", 302)
	} else {
		fk.CookieSetFlash(cookieErrLogin, w, "Invalid email or password")
		http.Redirect(w, r, "/login", 302)
	}

}

// signup page with GET and Post methods
func signup(w http.ResponseWriter, r *http.Request) {
	// Get /signup
	if r.Method == http.MethodGet {
		data, _ := csrfAuth(r)
		tmplExecute(w, r, "signup", data)
	}

	// POST /signup
	if r.Method == http.MethodPost {
		err := r.ParseForm()
		if err != nil {
			log.Println(err, "cannot parse form")
			return
		}

		user := auth.User{
			Name:     r.PostFormValue("name"),
			Email:    r.PostFormValue("email"),
			Password: r.PostFormValue("password"),
		}

		hasErr, errorsSignup := user.Validate()
		if hasErr == true {
			err := t.Lookup("signup").Execute(w, H{
				"errName":        auth.GetMapValue(errorsSignup, "name"),
				"errEmail":       auth.GetMapValue(errorsSignup, "email"),
				"errPassword":    auth.GetMapValue(errorsSignup, "password"),
				"user":           user,
				csrf.TemplateTag: csrf.TemplateField(r),
			})
			if err != nil {
				log.Println(err, "cannot execute template")
			}
			return
		}

		// check if email already exists into db
		check, err := user.Exists()
		if err != nil {
			log.Println(err, "cannot check user")
			return
		} else if check {
			fk.CookieSetFlash(cookieErrSignup, w, "Email already exists")
			http.Redirect(w, r, "/signup", 302)
			return
		}

		// create the user
		if err = user.Create(); err != nil {
			log.Println(err, "cannot create user")
			fk.CookieSetFlash(cookieErrSignup, w, "Cannot create user")
		}

		fk.CookieSetFlash(cookieErrLogin, w, "Enter your email and password to log in")
		http.Redirect(w, r, "/login", 302)
	}

}

// POST /signup
// Create the user account
func signupAccount(w http.ResponseWriter, r *http.Request) {
	if isAllowed := allowMethod(w, r, http.MethodPost); !isAllowed {
		return
	}

	err := r.ParseForm()
	if err != nil {
		log.Println(err, "cannot parse form")
		return
	}

	user := auth.User{
		Name:     r.PostFormValue("name"),
		Email:    r.PostFormValue("email"),
		Password: r.PostFormValue("password"),
	}

	hasErr, errorsSignup := user.Validate()
	if hasErr == true {
		err := t.Lookup("signup").Execute(w, H{
			"errors":         errorsSignup,
			"user":           user,
			csrf.TemplateTag: csrf.TemplateField(r),
		})
		if err != nil {
			log.Println(err, "cannot execute template")
		}
		return
	}

	// check if email already exists into db
	check, err := user.Exists()
	if err != nil {
		log.Println(err, "cannot check user")
		return
	} else if check {
		fk.CookieSetFlash(cookieErrSignup, w, "Email already exists")
		http.Redirect(w, r, "/signup", 302)
		return
	}

	// create the user
	if err = user.Create(); err != nil {
		log.Println(err, "cannot create user")
		fk.CookieSetFlash(cookieErrSignup, w, "Cannot create user")
	}

	fk.CookieSetFlash(cookieErrLogin, w, "Enter your email and password to log in")
	http.Redirect(w, r, "/login", 302)
}

// GET /logout
// Logs the user out
func logout(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("_cookie")
	if err != http.ErrNoCookie {
		session := auth.Session{Uuid: cookie.Value}
		session.Delete()

		err := fk.CookieDeleteSession("_cookie", w, r)
		if err != nil {
			log.Println(err, "cannot delete session")
		}
	}
	http.Redirect(w, r, "/", 302)
}

func users(w http.ResponseWriter, r *http.Request) {
	data, _, next := methodAuth(w, r, http.MethodGet)
	if next {
		users, err := auth.Users()
		if err != nil {
			log.Println(err, "cannot get all users")
			return
		}
		data["users"] = users
		tmplExecute(w, r, "users", data)
	}
}
