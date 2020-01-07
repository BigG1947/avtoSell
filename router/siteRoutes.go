package router

import (
	"avtoSell/model"
	"github.com/gorilla/mux"
	"html/template"
	"log"
	"net/http"
	"strconv"
)

func index(writer http.ResponseWriter, request *http.Request) {
	session, err := sessionStore.Get(request, userSession)
	if err != nil {
		log.Printf("Error with session in user routes cabinet: %s\n", err)
		return
	}
	var user *model.User

	if !isAuthUser(session) {
		user = new(model.User)
	} else {
		user = session.Values["user"].(*model.User)
		user.GetById(connection, user.Id)
	}

	var nl model.NewsList
	nl.GetLatest(connection)

	var cl model.CarList
	cl.GetLatest(connection)

	tmpl := template.Must(template.ParseFiles("templates/index.html"))
	err = tmpl.Execute(writer, map[string]interface{}{
		"news_list": nl,
		"car_list":  cl,
		"isAuth":    isAuthUser(session),
		"user":      user,
	})
	if err != nil {
		log.Printf("Error in site routes 'index': %s\n", err)
	}
	return
}

func news(writer http.ResponseWriter, request *http.Request) {
	session, err := sessionStore.Get(request, userSession)
	if err != nil {
		log.Printf("Error with session in user routes cabinet: %s\n", err)
		return
	}
	var user *model.User

	if !isAuthUser(session) {
		user = new(model.User)
	} else {
		user = session.Values["user"].(*model.User)
		user.GetById(connection, user.Id)
	}

	tmpl := template.Must(template.ParseFiles("templates/news.html"))
	err = tmpl.Execute(writer, map[string]interface{}{
		"isAuth": isAuthUser(session),
		"user":   user,
	})
	if err != nil {
		log.Printf("Error in site routes 'news': %s\n", err)
	}
	return
}

func post(writer http.ResponseWriter, request *http.Request) {
	session, err := sessionStore.Get(request, userSession)
	if err != nil {
		log.Printf("Error with session in user routes cabinet: %s\n", err)
		return
	}
	var user *model.User

	if !isAuthUser(session) {
		user = new(model.User)
	} else {
		user = session.Values["user"].(*model.User)
		user.GetById(connection, user.Id)
	}

	id, err := strconv.Atoi(mux.Vars(request)["id"])
	if err != nil {
		log.Printf("%s\n", err)
		return
	}
	var news model.News
	if err := news.Get(connection, id); err != nil {
		log.Printf("%s\n", err)
		return
	}

	text := template.HTML(news.Description)

	tmpl := template.Must(template.ParseFiles("templates/post.html"))
	err = tmpl.Execute(writer, map[string]interface{}{
		"news":   news,
		"text":   text,
		"isAuth": isAuthUser(session),
		"user":   user,
	})
	if err != nil {
		log.Printf("Error in site routes 'post': %s\n", err)
	}
	return
}

func product(writer http.ResponseWriter, request *http.Request) {
	var err error
	session, err := sessionStore.Get(request, userSession)
	if err != nil {
		log.Printf("%s\n")
		return
	}

	var c model.Car
	if c.Id, err = strconv.Atoi(mux.Vars(request)["id"]); err != nil {
		log.Printf("%s\n", err)
		return
	}
	if err = c.Get(connection, c.Id); err != nil {
		log.Printf("%s\n", err)
		return
	}
	var user *model.User
	if isAuthUser(session) {
		user = session.Values["user"].(*model.User)
	} else {
		user = new(model.User)
	}

	if request.Method == http.MethodPost {
		var order model.Order
		order.User.Id, _ = strconv.Atoi(request.FormValue("user"))
		order.Car.Id, _ = strconv.Atoi(request.FormValue("car"))
		order.Date = request.FormValue("order-date")
		order.Status = 1
		if err := order.Add(connection); err != nil {
			log.Printf("%s\n", err)
			return
		}
		http.Redirect(writer, request, "/cabinet", 302)
		return
	}

	tmpl := template.Must(template.ParseFiles("templates/product.html"))
	err = tmpl.Execute(writer, map[string]interface{}{
		"car":    c,
		"text":   template.HTML(c.Description),
		"isAuth": isAuthUser(session),
		"user":   user,
	})
	if err != nil {
		log.Printf("Error in site routes 'product': %s\n", err)
	}
	return
}

func catalog(writer http.ResponseWriter, request *http.Request) {
	session, err := sessionStore.Get(request, userSession)
	if err != nil {
		log.Printf("Error with session in user routes cabinet: %s\n", err)
		return
	}
	var user *model.User

	if !isAuthUser(session) {
		user = new(model.User)
	} else {
		user = session.Values["user"].(*model.User)
		user.GetById(connection, user.Id)
	}

	var ml model.ManufacturerList
	if err := ml.GetAll(connection); err != nil {
		log.Printf("%s\n", err)
		return
	}

	var cl model.ColorList
	if err := cl.GetAll(connection); err != nil {
		log.Printf("%s\n", err)
		return
	}

	var yl model.YearsList
	if err := yl.Get(connection); err != nil {
		log.Printf("%s\n", err)
		return
	}

	var ct model.CategoryList
	if err := ct.GetAll(connection); err != nil {
		log.Printf("%s\n", err)
		return
	}

	var cars model.CarList
	minPrice, _ := cars.GetMinPrice(connection)
	maxPrice, _ := cars.GetMaxPrice(connection)
	tmpl := template.Must(template.ParseFiles("templates/catalog.html"))
	err = tmpl.Execute(writer, map[string]interface{}{
		"manufacturer": ml,
		"colors":       cl,
		"years":        yl,
		"category":     ct,
		"minPrice":     minPrice,
		"maxPrice":     maxPrice,
		"isAuth":       isAuthUser(session),
		"user":         user,
	})
	if err != nil {
		log.Printf("Error in site routes 'catalog': %s\n", err)
	}
	return
}
