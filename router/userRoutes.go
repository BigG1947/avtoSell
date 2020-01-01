package router

import (
	"avtoSell/model"
	"crypto/sha256"
	"html/template"
	"log"
	"net/http"
)

func signUp(writer http.ResponseWriter, request *http.Request) {
	session, err := sessionStore.Get(request, userSession)
	if err != nil{
		log.Printf("Error in getting session: %s\n", err)
		return
	}
	if isAuthUser(session){
		http.Redirect(writer, request, "/cabinet", 302)
		return
	}

	if request.Method == http.MethodPost {
		var user model.User
		user.Login = request.FormValue("login")
		user.FirstName = request.FormValue("first_name")
		user.LastName = request.FormValue("last_name")
		passwordHash := model.GeneratePasswordHash(request.FormValue("password"))
		user.PasswordHash = passwordHash[:]
		user.Email = request.FormValue("email")
		user.Phone = request.FormValue("phone")

		if ok, err := user.Registration(connection); !ok{
			log.Printf("Error in registration user: %s\n", err)
		}else{
			tmpl := template.Must(template.ParseFiles("templates/users/registrationSuccess.html"))
			err := tmpl.Execute(writer, nil)
			if err != nil {
				log.Printf("Error in user routes signUp: %s\n", err)
			}
			return
		}
	}

	tmpl := template.Must(template.ParseFiles("templates/users/signUp.html"))
	err = tmpl.Execute(writer, nil)
	if err != nil {
		log.Printf("Error in user routes signUp: %s\n", err)
	}
}

func signIn(writer http.ResponseWriter, request *http.Request) {
	session, err := sessionStore.Get(request, userSession)
	if err != nil{
		log.Printf("Error in getting session: %s\n", err)
		return
	}
	if isAuthUser(session){
		http.Redirect(writer, request, "/cabinet", 302)
		return
	}

	var user model.User
	var errorMessage string
	if request.Method == http.MethodPost{
		user.Email = request.FormValue("email")
		passwordHash := sha256.Sum256([]byte(request.FormValue("password")))
		user.PasswordHash = passwordHash[:]
		if user.CheckUser(connection){
			session, err := sessionStore.Get(request, userSession)
			if err != nil {
				log.Printf("error in getting user sessions: %s\n", err)
				return
			}
			err = user.GetByEmail(connection, user.Email)
			if err != nil{
				log.Printf("Error in getting user by Email: %s\n", err)
				return
			}
			session.Values["user"] = user
			session.Save(request, writer)
			http.Redirect(writer, request, "/cabinet", 302)
			return
		}else{
			errorMessage = "Данные указаны неверно!"
		}
	}

	tmpl := template.Must(template.ParseFiles("templates/users/signIn.html"))
	err = tmpl.Execute(writer, map[string]interface{}{
		"user": user,
		"errorMessage": errorMessage,
	})
	if err != nil {
		log.Printf("Error in user routes signIn: %s\n", err)
	}
}

func cabinet(writer http.ResponseWriter, request *http.Request) {
	session, err := sessionStore.Get(request, userSession)
	if err != nil{
		log.Printf("Error with session in user routes cabinet: %s\n", err)
		return
	}
	if !isAuthUser(session){
		http.Redirect(writer, request, "/login", 302)
		return
	}
	var user *model.User
	user = session.Values["user"].(*model.User)

	tmpl := template.Must(template.ParseFiles("templates/users/cabinet.html"))
	err = tmpl.Execute(writer, map[string]interface{}{"user": user})
	if err != nil{
		log.Printf("Error in user routes cabinet: %s\n", err)
	}
	return
}

func exitUser(writer http.ResponseWriter, request *http.Request){
	session, err := sessionStore.Get(request, userSession)
	if err != nil{
		log.Printf("Error in user routes exitUser with getting session: %s\n", err)
		return
	}
	if !isAuthUser(session){
		http.Redirect(writer, request, "/login", 302)
		return
	}

	session.Options.MaxAge = -1
	session.Save(request, writer)
	http.Redirect(writer, request, "/login", 302)
	return
}