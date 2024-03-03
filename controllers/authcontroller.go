package controllers

import (
	"errors"
	"go-auth/config"
	"go-auth/entities"
	"go-auth/models"
	"html/template"
	"net/http"

	"golang.org/x/crypto/bcrypt"
)

type UserInput struct {
	Username string
	Password string
}

var userModel = models.NewUserModel()

func Index(w http.ResponseWriter, r *http.Request) {

	session, _ := config.Store.Get(r, config.SESSION_ID)

	if len(session.Values) == 0 { // check session
		http.Redirect(w, r, "/login", http.StatusSeeOther)
	} else {
		if session.Values["loggedIn"] != true { // check session
			http.Redirect(w, r, "/login", http.StatusSeeOther)
		} else {

			data := map[string]interface{}{
				"name": session.Values["name"],
			}

			temp, _ := template.ParseFiles("views/index.html")
			temp.Execute(w, data)
		}
	}
}

func Login(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		temp, _ := template.ParseFiles("views/login.html")
		temp.Execute(w, nil)
	} else if r.Method == http.MethodPost {

		// get user input from login page

		r.ParseForm()

		UserInput := &UserInput{
			Username: r.Form.Get("username"),
			Password: r.Form.Get("password"),
		}

		// check username

		var user entities.User
		userModel.Where(&user, "username", UserInput.Username)

		var message error

		if user.Username == "" {
			message = errors.New("Your username or password is incorrect!")
		} else {

			// check password

			errPassword := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(UserInput.Password))

			if errPassword != nil {
				message = errors.New("Your username or password is incorrect!")
			}
		}

		if message != nil { // if login is error

			data := map[string]interface{}{
				"error": message,
			}

			temp, _ := template.ParseFiles("views/login.html")
			temp.Execute(w, data)
		} else {
			// set session

			session, _ := config.Store.Get(r, config.SESSION_ID)

			session.Values["loggedIn"] = true
			session.Values["name"] = user.Name
			session.Values["email"] = user.Email
			session.Values["username"] = user.Username

			session.Save(r, w)

			http.Redirect(w, r, "/", http.StatusSeeOther)
		}
	}
}

func Logout(w http.ResponseWriter, r *http.Request) {
	session, _ := config.Store.Get(r, config.SESSION_ID)

	// delete session
	session.Options.MaxAge = -1
	session.Save(r, w)
	http.Redirect(w, r, "/login", http.StatusSeeOther)
}

func Register(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		temp, _ := template.ParseFiles("views/register.html")
		temp.Execute(w, nil)
	} else if r.Method == http.MethodPost {

		// get user input

		r.ParseForm()

		user := entities.User{
			Name:      r.Form.Get("name"),
			Email:     r.Form.Get("email"),
			Username:  r.Form.Get("username"),
			Password:  r.Form.Get("password"),
			Cpassword: r.Form.Get("cpassword"),
		}

		// form validation

		errorMessages := make(map[string]interface{})

		if user.Name == "" {
			errorMessages["Name"] = "Name is required"
		}

		if user.Email == "" {
			errorMessages["Email"] = "Email is required"
		}

		if user.Username == "" {
			errorMessages["Username"] = "Username is required"
		}

		if user.Password == "" {
			errorMessages["Password"] = "Password is required"
		}

		if user.Cpassword == "" {
			errorMessages["Cpassword"] = "Confirm Password is required"
		} else {
			if user.Cpassword != user.Password {
				errorMessages["Cpassword"] = "Confirm Password doesn't match"
			}
		}

		if len(errorMessages) > 0 {
			data := map[string]interface{}{
				"validation": errorMessages,
			}

			temp, _ := template.ParseFiles("views/register.html")
			temp.Execute(w, data)
		} else {

			// register process to database

		}
	}
}
