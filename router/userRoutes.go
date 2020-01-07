package router

import (
	"avtoSell/model"
	"bytes"
	"crypto/sha256"
	"encoding/json"
	"github.com/gorilla/mux"
	"html/template"
	"log"
	"net/http"
	"strconv"
)

func signUp(writer http.ResponseWriter, request *http.Request) {
	session, err := sessionStore.Get(request, userSession)
	if err != nil {
		log.Printf("Error in getting session: %s\n", err)
		return
	}
	if isAuthUser(session) {
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

		if ok, err := user.Registration(connection); !ok {
			log.Printf("Error in registration user: %s\n", err)
		} else {
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
	if err != nil {
		log.Printf("Error in getting session: %s\n", err)
		return
	}
	if isAuthUser(session) {
		http.Redirect(writer, request, "/cabinet", 302)
		return
	}

	var user model.User
	var errorMessage string
	if request.Method == http.MethodPost {
		user.Email = request.FormValue("email")
		passwordHash := sha256.Sum256([]byte(request.FormValue("password")))
		user.PasswordHash = passwordHash[:]
		if user.CheckUser(connection) {
			session, err := sessionStore.Get(request, userSession)
			if err != nil {
				log.Printf("error in getting user sessions: %s\n", err)
				return
			}
			err = user.GetByEmail(connection, user.Email)
			if err != nil {
				log.Printf("Error in getting user by Email: %s\n", err)
				return
			}
			session.Values["user"] = user
			session.Save(request, writer)
			http.Redirect(writer, request, "/cabinet", 302)
			return
		} else {
			errorMessage = "Данные указаны неверно!"
		}
	}

	tmpl := template.Must(template.ParseFiles("templates/users/signIn.html"))
	err = tmpl.Execute(writer, map[string]interface{}{
		"user":         user,
		"errorMessage": errorMessage,
	})
	if err != nil {
		log.Printf("Error in user routes signIn: %s\n", err)
	}
}

func cabinet(writer http.ResponseWriter, request *http.Request) {
	session, err := sessionStore.Get(request, userSession)
	if err != nil {
		log.Printf("Error with session in user routes cabinet: %s\n", err)
		return
	}
	if !isAuthUser(session) {
		http.Redirect(writer, request, "/login", 302)
		return
	}
	var user *model.User
	user = session.Values["user"].(*model.User)
	user.GetById(connection, user.Id)

	var newOrderList model.OrderList
	if err := newOrderList.GetUserNewOrders(connection, user.Id); err != nil {
		log.Printf("%s\n", err)
		return
	}

	var cancelOrderList model.OrderList
	if err := cancelOrderList.GetUserCancelOrders(connection, user.Id); err != nil {
		log.Printf("%s\n", err)
		return
	}

	var checkOrderList model.OrderList
	if err := checkOrderList.GetUserCheckOrders(connection, user.Id); err != nil {
		log.Printf("%s\n", err)
		return
	}

	tmpl := template.Must(template.ParseFiles("templates/users/cabinet.html"))
	err = tmpl.Execute(writer, map[string]interface{}{
		"user":         user,
		"newOrders":    newOrderList,
		"checkOrders":  checkOrderList,
		"cancelOrders": cancelOrderList,
	})
	if err != nil {
		log.Printf("Error in user routes cabinet: %s\n", err)
	}
	return
}

func exitUser(writer http.ResponseWriter, request *http.Request) {
	session, err := sessionStore.Get(request, userSession)
	if err != nil {
		log.Printf("Error in user routes exitUser with getting session: %s\n", err)
		return
	}
	if !isAuthUser(session) {
		http.Redirect(writer, request, "/login", 302)
		return
	}

	session.Options.MaxAge = -1
	session.Save(request, writer)
	http.Redirect(writer, request, "/login", 302)
	return
}

func changePassword(writer http.ResponseWriter, request *http.Request) {
	session, err := sessionStore.Get(request, userSession)
	if err != nil {
		log.Printf("Error with session in user routes cabinet: %s\n", err)
		return
	}
	if !isAuthUser(session) {
		http.Redirect(writer, request, "/login", 302)
		return
	}
	var user *model.User
	user = session.Values["user"].(*model.User)
	user.GetById(connection, user.Id)
	writer.WriteHeader(200)
	responseMap := make(map[string]interface{})
	responseMap["ok"] = false
	oldPasswordHash := model.GeneratePasswordHash(request.FormValue("old_password"))
	if bytes.Equal(oldPasswordHash[:], user.PasswordHash) {
		responseMap["ok"] = true
		newPasswordHash := model.GeneratePasswordHash(request.FormValue("new_password"))
		user.PasswordHash = newPasswordHash[:]
		user.EditPassword(connection)
		session.Values["user"] = user
		_ = session.Save(request, writer)
	}
	responseJson, _ := json.Marshal(responseMap)
	writer.Write(responseJson)
	return
}

func changePhone(writer http.ResponseWriter, request *http.Request) {
	session, err := sessionStore.Get(request, userSession)
	if err != nil {
		log.Printf("Error with session in user routes cabinet: %s\n", err)
		return
	}
	if !isAuthUser(session) {
		http.Redirect(writer, request, "/login", 302)
		return
	}
	responseMap := make(map[string]interface{})
	responseMap["ok"] = false
	var user *model.User
	user = session.Values["user"].(*model.User)
	user.GetById(connection, user.Id)
	writer.WriteHeader(200)
	user.Phone = request.FormValue("new_phone")
	user.EditPhone(connection)
	session.Values["user"] = user
	_ = session.Save(request, writer)
	responseMap["ok"] = true
	responseJson, _ := json.Marshal(responseMap)
	writer.Write(responseJson)
	return
}

func changeFio(writer http.ResponseWriter, request *http.Request) {
	session, err := sessionStore.Get(request, userSession)
	if err != nil {
		log.Printf("Error with session in user routes cabinet: %s\n", err)
		return
	}
	if !isAuthUser(session) {
		http.Redirect(writer, request, "/login", 302)
		return
	}
	responseMap := make(map[string]interface{})
	responseMap["ok"] = false
	var user *model.User
	user = session.Values["user"].(*model.User)
	user.GetById(connection, user.Id)
	writer.WriteHeader(200)
	user.LastName = request.FormValue("new_last_name")
	user.FirstName = request.FormValue("new_first_name")
	user.EditFio(connection)
	session.Values["user"] = user
	_ = session.Save(request, writer)
	responseMap["ok"] = true
	responseJson, _ := json.Marshal(responseMap)
	writer.Write(responseJson)
	return
}

func changeEmail(writer http.ResponseWriter, request *http.Request) {
	session, err := sessionStore.Get(request, userSession)
	if err != nil {
		log.Printf("Error with session in user routes cabinet: %s\n", err)
		return
	}
	if !isAuthUser(session) {
		http.Redirect(writer, request, "/login", 302)
		return
	}
	responseMap := make(map[string]interface{})
	responseMap["ok"] = false
	var user *model.User
	user = session.Values["user"].(*model.User)
	user.GetById(connection, user.Id)
	writer.WriteHeader(200)
	user.Email = request.FormValue("new_email")
	user.EditEmail(connection)
	session.Values["user"] = user
	_ = session.Save(request, writer)
	responseMap["ok"] = true
	responseJson, _ := json.Marshal(responseMap)
	writer.Write(responseJson)
	return
}

func cancelOrder(writer http.ResponseWriter, request *http.Request) {
	session, err := sessionStore.Get(request, userSession)
	if err != nil {
		log.Printf("Error with session in user routes cabinet: %s\n", err)
		return
	}
	if !isAuthUser(session) {
		http.Redirect(writer, request, "/login", 302)
		return
	}

	var user *model.User
	user = session.Values["user"].(*model.User)
	user.GetById(connection, user.Id)

	var order model.Order

	order.Id, _ = strconv.Atoi(mux.Vars(request)["id"])
	log.Printf("userId: %d;\norderId: %d\n", user.Id, order.Id)
	log.Printf("%s\n", order.Cancel(connection, user.Id))
	http.Redirect(writer, request, "/cabinet", 302)
	return
}
