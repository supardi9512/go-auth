package controllers

import (
	"errors"
	"go-auth/config"
	"go-auth/entities"
	"go-auth/libraries"
	"go-auth/models"
	"html/template"
	"net/http"

	"golang.org/x/crypto/bcrypt"
)

type UserInput struct {
	Username string `validate:"required"`
	Password string `validate:"required"`
}

var userModel = models.NewUserModel()
var validation = libraries.NewValidation()

func Index(w http.ResponseWriter, r *http.Request) {

	session, _ := config.Store.Get(r, config.SESSION_ID)

	if len(session.Values) == 0 && session.Values["loggedIn"] != true { // check session
		http.Redirect(w, r, "/login", http.StatusSeeOther)
	} else {
		data := map[string]interface{}{
			"name": session.Values["name"],
		}

		temp, _ := template.ParseFiles("views/index.html")
		temp.Execute(w, data)
	}
}

func Login(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		session, _ := config.Store.Get(r, config.SESSION_ID)

		if len(session.Values) > 0 && session.Values["loggedIn"] == true {
			http.Redirect(w, r, "/", http.StatusSeeOther)
		} else {
			temp, _ := template.ParseFiles("views/login.html")
			temp.Execute(w, nil)
		}
	} else if r.Method == http.MethodPost {

		// get user input from login page

		r.ParseForm()

		UserInput := &UserInput{
			Username: r.Form.Get("username"),
			Password: r.Form.Get("password"),
		}

		errorMessages := validation.Struct(UserInput)

		if errorMessages != nil {
			data := map[string]interface{}{
				"validation": errorMessages,
			}

			temp, _ := template.ParseFiles("views/login.html")
			temp.Execute(w, data)
		} else {

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
		session, _ := config.Store.Get(r, config.SESSION_ID)

		if len(session.Values) > 0 && session.Values["loggedIn"] == true {
			http.Redirect(w, r, "/", http.StatusSeeOther)
		} else {
			temp, _ := template.ParseFiles("views/register.html")
			temp.Execute(w, nil)
		}
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

		errorMessages := validation.Struct(user)

		if errorMessages != nil {
			data := map[string]interface{}{
				"validation": errorMessages,
				"user":       user,
			}

			temp, _ := template.ParseFiles("views/register.html")
			temp.Execute(w, data)
		} else {
			// hash password

			hashPassword, _ := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
			user.Password = string(hashPassword)

			// insert to database

			userModel.Create(user)

			data := map[string]interface{}{
				"message": "Registration successful",
			}

			temp, _ := template.ParseFiles("views/register.html")
			temp.Execute(w, data)
		}
	}
}
