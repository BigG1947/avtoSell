package router

import (
	"avtoSell/model"
	"crypto/sha256"
	"html/template"
	"log"
	"net/http"
)

func admin(writer http.ResponseWriter, request *http.Request) {
	session, err := sessionStore.Get(request, adminSession)
	if err != nil {
		log.Printf("Error in admin routes adminLogin with getting session: %s\n", err)
	}
	if !isAuthAdmin(session) {
		http.Redirect(writer, request, "/admin/login", 302)
		return
	}

	tmpl := template.Must(template.ParseFiles("templates/admin/index.html"))
	err = tmpl.Execute(writer, map[string]interface{}{})
	if err != nil {
		log.Printf("Error in admin routes `admin`: %s\n", err)
	}
	return
}

func adminLogin(writer http.ResponseWriter, request *http.Request) {
	session, err := sessionStore.Get(request, adminSession)
	if err != nil {
		log.Printf("Error in admin routes adminLogin with getting session: %s\n", err)
		return
	}
	if isAuthAdmin(session) {
		http.Redirect(writer, request, "/admin", 302)
		return
	}

	var admin model.User
	var errorMessage string

	if request.Method == http.MethodPost {
		admin.Login = request.FormValue("login")
		passwordHash := sha256.Sum256([]byte(request.FormValue("password")))
		admin.PasswordHash = passwordHash[:]
		if admin.CheckAdmin(connection) {
			err = admin.GetByLogin(connection, admin.Login)
			if err != nil {
				log.Printf("Error in getting user by login: %s\n", err)
				return
			}
			session.Values["admin"] = admin
			session.Save(request, writer)
			http.Redirect(writer, request, "/admin", 302)
			return
		}
		errorMessage = "Введены некоректные данные!"
	}

	tmpl := template.Must(template.ParseFiles("templates/admin/login.html"))
	err = tmpl.Execute(writer, map[string]interface{}{
		"user":         admin,
		"errorMessage": errorMessage,
	})
	if err != nil {
		log.Printf("Error in admin routes adminLogin: %s\n", err)
	}
	return
}

func adminExit(writer http.ResponseWriter, request *http.Request) {
	session, err := sessionStore.Get(request, adminSession)
	if err != nil{
		log.Printf("Error in getting session: %s\n", err)
		return
	}
	if !isAuthAdmin(session){
		http.Redirect(writer, request, "/admin/login", 302)
		return
	}

	session.Options.MaxAge = -1
	session.Save(request, writer)
	http.Redirect(writer, request, "/admin/login", 302)
	return
}