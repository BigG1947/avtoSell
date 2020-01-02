package router

import (
	"avtoSell/model"
	"crypto/sha256"
	"github.com/gorilla/mux"
	"html/template"
	"log"
	"net/http"
	"strconv"
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

	var news model.NewsList
	if err := news.GetAll(connection); err != nil {
		log.Printf("Error in NewsList.GetAll() method: %s\n", err)
		return
	}

	tmpl := template.Must(template.ParseFiles("templates/admin/index.html"))
	err = tmpl.Execute(writer, map[string]interface{}{
		"news": news,
	})
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
	if err != nil {
		log.Printf("Error in getting session: %s\n", err)
		return
	}
	if !isAuthAdmin(session) {
		http.Redirect(writer, request, "/admin/login", 302)
		return
	}

	session.Options.MaxAge = -1
	session.Save(request, writer)
	http.Redirect(writer, request, "/admin/login", 302)
	return
}

func adminNewsAdd(writer http.ResponseWriter, request *http.Request) {
	session, err := sessionStore.Get(request, adminSession)
	if err != nil {
		log.Printf("Error in getting session: %s\n", err)
		return
	}
	if !isAuthAdmin(session) {
		http.Redirect(writer, request, "/admin/login", 302)
		return
	}

	if request.Method == http.MethodPost {
		err = request.ParseMultipartForm(22 << 20)
		if err != nil {
			http.Redirect(writer, request, "/500", 302)
			log.Printf("%s\n", err)
			return
		}

		file, header, err := request.FormFile("images")
		if err != nil {
			http.Redirect(writer, request, "/500", 302)
			log.Printf("%s\n", err)
			return
		}
		defer file.Close()
		if (header.Header.Get("Content-Type") != "image/jpeg" && header.Header.Get("Content-Type") != "image/png") || (header.Size > 2<<20) {
			log.Printf("AddProject: Неверный тип формата файла: %s\n", header.Header.Get("Content-Type"))
			http.Redirect(writer, request, "/413", 302)
			return
		}

		images, err := UploadImages(file)
		if err != nil {
			http.Redirect(writer, request, "/500", 302)
			log.Printf("%s\n", err)
			return
		}

		title := request.FormValue("title")
		miniDesc := request.FormValue("mini-desc")
		description := request.FormValue("description")
		news := new(model.News)
		news.Title = title
		news.MiniDesc = miniDesc
		news.Description = description
		news.Image = images
		if err := news.Add(connection); err != nil {
			log.Printf("User add error: %s\n", err)
			return
		}
		http.Redirect(writer, request, "/admin", 302)
		return
	}

	tmpl := template.Must(template.ParseFiles("templates/admin/addNews.html"))
	if err = tmpl.Execute(writer, nil); err != nil {
		log.Printf("Error in admin routes News Add: %s\n", err)
		return
	}
}

func adminNewsEdit(writer http.ResponseWriter, request *http.Request) {
	session, err := sessionStore.Get(request, adminSession)
	if err != nil {
		log.Printf("Error in getting session: %s\n", err)
		return
	}
	if !isAuthAdmin(session) {
		http.Redirect(writer, request, "/admin/login", 302)
		return
	}

	var news model.News
	id, _ := strconv.Atoi(mux.Vars(request)["id"])
	if err := news.Get(connection, id); err != nil {
		log.Printf("Error in User.Get(id) method: %s\n", err)
		return
	}

	if request.Method == http.MethodPost {
		err = request.ParseMultipartForm(22 << 20)
		if err != nil {
			http.Redirect(writer, request, "/500", 302)
			log.Printf("%s\n", err)
			return
		}

		var images string

		file, header, err := request.FormFile("images")
		if err == nil {
			defer file.Close()
			if header.Header.Get("Content-Type") != "image/jpeg" && header.Header.Get("Content-Type") != "image/png" {
				log.Printf("Edit news: Неверный тип формата файла: %s\n", header.Header.Get("Content-Type"))
				http.Redirect(writer, request, "/500", 302)
				return
			}
			images, err = UploadImages(file)
			if err != nil {
				http.Redirect(writer, request, "/500", 302)
				log.Printf("%s\n", err)
				return
			}

			err = DeleteImages(news.Image)
			if err != nil {
				http.Redirect(writer, request, "/500", 302)
				log.Printf("%s\n", err)
				return
			}
		} else if err != http.ErrMissingFile {
			http.Redirect(writer, request, "/500", 302)
			log.Printf("%s\n", err)
			return
		}

		title := request.FormValue("title")
		miniDesc := request.FormValue("mini-desc")
		description := request.FormValue("description")
		news.Title = title
		news.MiniDesc = miniDesc
		news.Description = description
		news.Image = images
		if err := news.Edit(connection); err != nil {
			log.Printf("User.Edit error: %s\n", err)
			return
		}
		http.Redirect(writer, request, "/admin", 302)
		return
	}

	tmpl := template.Must(template.ParseFiles("templates/admin/editNews.html"))
	if err = tmpl.Execute(writer, news); err != nil {
		log.Printf("Error in admin routes News Edit: %s\n", err)
		return
	}
}

func adminNewsDelete(writer http.ResponseWriter, request *http.Request) {
	session, err := sessionStore.Get(request, adminSession)
	if err != nil {
		log.Printf("Error in getting session: %s\n", err)
		return
	}
	if !isAuthAdmin(session) {
		http.Redirect(writer, request, "/admin/login", 302)
		return
	}

	var news model.News
	id, _ := strconv.Atoi(mux.Vars(request)["id"])
	if err := news.Get(connection, id); err != nil {
		log.Printf("Error in User.Get(id) method: %s\n", err)
		return
	}

	if err := news.Delete(connection); err != nil {
		log.Printf("Error in News.Delete() method: %s\n", err)
		return
	}

	DeleteImages(news.Image)
	http.Redirect(writer, request, "/admin", 302)
}
