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

	var colors model.ColorList
	if err := colors.GetAll(connection); err != nil {
		log.Printf("Error in ColorList.GetAll() method: %s\n", err)
		return
	}

	var manufacturer model.ManufacturerList
	if err := manufacturer.GetAll(connection); err != nil {
		log.Printf("Error in ManufacturerList.GetAll() method: %s\n", err)
		return
	}

	var category model.CategoryList
	if err := category.GetAll(connection); err != nil {
		log.Printf("Error in CategoryList.GetAll() method: %s\n", err)
		return
	}

	var cancelOrders model.OrderList
	cancelOrders.GetUserCancelOrders(connection, 0)

	var checkOrders model.OrderList
	checkOrders.GetUserCheckOrders(connection, 0)

	var newOrders model.OrderList
	newOrders.GetUserNewOrders(connection, 0)

	var cars model.CarList
	if err := cars.GetAll(connection); err != nil {
		log.Printf("Error in CarList.GetAll() method: %s\n", err)
		return
	}

	tmpl := template.Must(template.ParseFiles("templates/admin/index.html"))
	err = tmpl.Execute(writer, map[string]interface{}{
		"news":         news,
		"colors":       colors,
		"manufacturer": manufacturer,
		"category":     category,
		"cars":         cars,
		"newOrders":    newOrders,
		"checkOrders":  checkOrders,
		"cancelOrders": cancelOrders,
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

func adminCarAdd(writer http.ResponseWriter, request *http.Request) {
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
		var car model.Car

		car.Model = request.FormValue("model")
		car.MiniDesc = request.FormValue("mini-desc")
		car.Description = request.FormValue("description")
		if car.Year, err = strconv.Atoi(request.FormValue("year")); err != nil {
			log.Printf("%s\n", err)
			return
		}
		if car.Price, err = strconv.Atoi(request.FormValue("price")); err != nil {
			log.Printf("%s\n", err)
			return
		}
		if car.Manufacturer.Id, err = strconv.Atoi(request.FormValue("manufacturer")); err != nil {
			log.Printf("%s\n", err)
			return
		}
		if car.Color.Id, err = strconv.Atoi(request.FormValue("color")); err != nil {
			log.Printf("%s\n", err)
			return
		}
		if car.Category.Id, err = strconv.Atoi(request.FormValue("category")); err != nil {
			log.Printf("%s\n", err)
			return
		}
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

		car.Images, err = UploadImages(file)
		if err != nil {
			http.Redirect(writer, request, "/500", 302)
			log.Printf("%s\n", err)
			return
		}

		photos := request.MultipartForm.File["photos"]

		for i := range photos {
			if (photos[i].Header.Get("Content-Type") != "image/jpeg" && photos[i].Header.Get("Content-Type") != "image/png") || photos[i].Size > 2<<20 {
				continue
			}

			file, err := photos[i].Open()
			if err != nil {
				http.Redirect(writer, request, "/500", 302)
				log.Printf("%s\n", err)
				return
			}

			src, err := UploadImages(file)
			if err != nil {
				http.Redirect(writer, request, "/500", 302)
				log.Printf("%s\n", err)
				return
			}

			car.SecondImages = append(car.SecondImages, src)
		}

		if err := car.Add(connection); err != nil {
			log.Printf("%s\n", err)
			return
		}

		http.Redirect(writer, request, "/admin", 302)
		return
	}

	var colors model.ColorList
	if err := colors.GetAll(connection); err != nil {
		log.Printf("%s\n", err)
		return
	}

	var manufacturer model.ManufacturerList
	if err := manufacturer.GetAll(connection); err != nil {
		log.Printf("%s\n", err)
		return
	}

	var category model.CategoryList
	if err := category.GetAll(connection); err != nil {
		log.Printf("%s\n", err)
		return
	}

	tmpl := template.Must(template.ParseFiles("templates/admin/addCars.html"))
	if err := tmpl.Execute(writer, map[string]interface{}{
		"colors":       colors,
		"manufacturer": manufacturer,
		"category":     category,
	}); err != nil {
		log.Printf("%s\n", err)
	}
}

func adminCarEdit(writer http.ResponseWriter, request *http.Request) {
	session, err := sessionStore.Get(request, adminSession)
	if err != nil {
		log.Printf("Error in getting session: %s\n", err)
		return
	}
	if !isAuthAdmin(session) {
		http.Redirect(writer, request, "/admin/login", 302)
		return
	}

	var car model.Car
	if car.Id, err = strconv.Atoi(mux.Vars(request)["id"]); err != nil {
		log.Printf("%s\n", err)
		return
	}
	if err := car.Get(connection, car.Id); err != nil {
		log.Printf("%s\n", err)
		return
	}

	if request.Method == http.MethodPost {
		var images string
		car.Model = request.FormValue("model")
		car.MiniDesc = request.FormValue("mini-desc")
		car.Description = request.FormValue("description")
		if car.Year, err = strconv.Atoi(request.FormValue("year")); err != nil {
			log.Printf("%s\n", err)
			return
		}
		if car.Price, err = strconv.Atoi(request.FormValue("price")); err != nil {
			log.Printf("%s\n", err)
			return
		}
		if car.Manufacturer.Id, err = strconv.Atoi(request.FormValue("manufacturer")); err != nil {
			log.Printf("%s\n", err)
			return
		}
		if car.Color.Id, err = strconv.Atoi(request.FormValue("color")); err != nil {
			log.Printf("%s\n", err)
			return
		}
		if car.Category.Id, err = strconv.Atoi(request.FormValue("category")); err != nil {
			log.Printf("%s\n", err)
			return
		}

		err = request.ParseMultipartForm(22 << 20)
		if err != nil {
			http.Redirect(writer, request, "/500", 302)
			log.Printf("%s\n", err)
			return
		}

		file, header, err := request.FormFile("images")
		if err == nil {
			log.Printf("images: %s\n", header)
			defer file.Close()
			if (header.Header.Get("Content-Type") != "image/jpeg" && header.Header.Get("Content-Type") != "image/png") || (header.Size > 2<<20) {
				log.Printf("AddProject: Неверный тип формата файла: %s\n", header.Header.Get("Content-Type"))
				http.Redirect(writer, request, "/500", 302)
				return
			}
			images, err = UploadImages(file)
			if err != nil {
				http.Redirect(writer, request, "/500", 302)
				log.Printf("%s\n", err)
				return
			}

			err = DeleteImages(car.Images)
			if err != nil {
				http.Redirect(writer, request, "/500", 302)
				log.Printf("%s\n", err)
				return
			}
			car.Images = images
		} else if err != http.ErrMissingFile {
			http.Redirect(writer, request, "/500", 302)
			log.Printf("%s\n", err)
			return
		}

		photos := request.MultipartForm.File["photos"]
		if len(photos) > 0 {
			for i := range car.SecondImages {
				if err := DeleteImages(car.SecondImages[i]); err != nil {
					log.Printf("%s\n", err)
					return
				}
			}
			for i := range photos {
				if (photos[i].Header.Get("Content-Type") != "image/jpeg" && photos[i].Header.Get("Content-Type") != "image/png") || photos[i].Size > 2<<20 {
					continue
				}

				file, err := photos[i].Open()
				if err != nil {
					http.Redirect(writer, request, "/500", 302)
					log.Printf("%s\n", err)
					return
				}

				src, err := UploadImages(file)
				if err != nil {
					http.Redirect(writer, request, "/500", 302)
					log.Printf("%s\n", err)
					return
				}

				car.SecondImages = append(car.SecondImages, src)
			}
		}

		err = car.Edit(connection)
		if err != nil {
			http.Redirect(writer, request, "/500", 302)
			log.Printf("%s\n", err)
			return
		}

		http.Redirect(writer, request, "/admin", 302)
		return
	}
	var colors model.ColorList
	if err := colors.GetAll(connection); err != nil {
		log.Printf("%s\n", err)
		return
	}

	var manufacturer model.ManufacturerList
	if err := manufacturer.GetAll(connection); err != nil {
		log.Printf("%s\n", err)
		return
	}

	var category model.CategoryList
	if err := category.GetAll(connection); err != nil {
		log.Printf("%s\n", err)
		return
	}

	tmpl := template.Must(template.ParseFiles("templates/admin/editCars.html"))
	if err := tmpl.Execute(writer, map[string]interface{}{
		"colors":       colors,
		"manufacturer": manufacturer,
		"category":     category,
		"car":          car,
	}); err != nil {
		log.Printf("%s\n", err)
	}
}

func adminCarDelete(writer http.ResponseWriter, request *http.Request) {
	session, err := sessionStore.Get(request, adminSession)
	if err != nil {
		log.Printf("Error in getting session: %s\n", err)
		return
	}
	if !isAuthAdmin(session) {
		http.Redirect(writer, request, "/admin/login", 302)
		return
	}

	var car model.Car
	if car.Id, err = strconv.Atoi(mux.Vars(request)["id"]); err != nil {
		log.Printf("%s\n", err)
		return
	}
	if err := car.Get(connection, car.Id); err != nil {
		log.Printf("%s\n", err)
		return
	}

	if err := car.Delete(connection); err != nil {
		log.Printf("%s\n", err)
		return
	}
	DeleteImages(car.Images)
	for i := range car.SecondImages {
		DeleteImages(car.SecondImages[i])
	}

	http.Redirect(writer, request, "/admin", 302)
}

// Foreign function
func adminColorsAdd(writer http.ResponseWriter, request *http.Request) {
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
		var c model.Color
		c.Name = request.FormValue("name")
		if err := c.Add(connection); err != nil {
			log.Printf("%s\n", err)
			return
		}
		http.Redirect(writer, request, "/admin", 302)
		return
	}

	tmpl := template.Must(template.ParseFiles("templates/admin/addForeign.html"))
	if err := tmpl.Execute(writer, nil); err != nil {
		log.Printf("%s\n", err)
	}
	return
}

func adminColorsEdit(writer http.ResponseWriter, request *http.Request) {
	session, err := sessionStore.Get(request, adminSession)
	if err != nil {
		log.Printf("Error in getting session: %s\n", err)
		return
	}
	if !isAuthAdmin(session) {
		http.Redirect(writer, request, "/admin/login", 302)
		return
	}
	var c model.Color
	if c.Id, err = strconv.Atoi(mux.Vars(request)["id"]); err != nil {
		log.Printf("%s\n", err)
		return
	}
	if err := c.Get(connection, c.Id); err != nil {
		log.Printf("%s\n", err)
		return
	}

	if request.Method == http.MethodPost {
		c.Name = request.FormValue("name")
		if err := c.Edit(connection); err != nil {
			log.Printf("%s\n", err)
			return
		}

		http.Redirect(writer, request, "/admin", 302)
		return
	}

	tmpl := template.Must(template.ParseFiles("templates/admin/editForeign.html"))
	if err := tmpl.Execute(writer, c); err != nil {
		log.Printf("%s\n", err)
	}
	return
}

func adminColorsDelete(writer http.ResponseWriter, request *http.Request) {
	session, err := sessionStore.Get(request, adminSession)
	if err != nil {
		log.Printf("Error in getting session: %s\n", err)
		return
	}
	if !isAuthAdmin(session) {
		http.Redirect(writer, request, "/admin/login", 302)
		return
	}
	var c model.Color
	if c.Id, err = strconv.Atoi(mux.Vars(request)["id"]); err != nil {
		log.Printf("%s\n", err)
		return
	}
	if err := c.Get(connection, c.Id); err != nil {
		log.Printf("%s\n", err)
		return
	}

	if err := c.Delete(connection); err != nil {
		log.Printf("%s\n", err)
		return
	}

	http.Redirect(writer, request, "/admin", 302)
	return
}

func adminCategoryAdd(writer http.ResponseWriter, request *http.Request) {
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
		var c model.Category
		c.Name = request.FormValue("name")
		if err := c.Add(connection); err != nil {
			log.Printf("%s\n", err)
			return
		}
		http.Redirect(writer, request, "/admin", 302)
		return
	}

	tmpl := template.Must(template.ParseFiles("templates/admin/addForeign.html"))
	if err := tmpl.Execute(writer, nil); err != nil {
		log.Printf("%s\n", err)
	}
	return
}

func adminCategoryEdit(writer http.ResponseWriter, request *http.Request) {
	session, err := sessionStore.Get(request, adminSession)
	if err != nil {
		log.Printf("Error in getting session: %s\n", err)
		return
	}
	if !isAuthAdmin(session) {
		http.Redirect(writer, request, "/admin/login", 302)
		return
	}
	var c model.Category
	if c.Id, err = strconv.Atoi(mux.Vars(request)["id"]); err != nil {
		log.Printf("%s\n", err)
		return
	}
	if err := c.Get(connection, c.Id); err != nil {
		log.Printf("%s\n", err)
		return
	}

	if request.Method == http.MethodPost {
		c.Name = request.FormValue("name")
		if err := c.Edit(connection); err != nil {
			log.Printf("%s\n", err)
			return
		}

		http.Redirect(writer, request, "/admin", 302)
		return
	}

	tmpl := template.Must(template.ParseFiles("templates/admin/editForeign.html"))
	if err := tmpl.Execute(writer, c); err != nil {
		log.Printf("%s\n", err)
	}
	return
}

func adminCategoryDelete(writer http.ResponseWriter, request *http.Request) {
	session, err := sessionStore.Get(request, adminSession)
	if err != nil {
		log.Printf("Error in getting session: %s\n", err)
		return
	}
	if !isAuthAdmin(session) {
		http.Redirect(writer, request, "/admin/login", 302)
		return
	}
	var c model.Category
	if c.Id, err = strconv.Atoi(mux.Vars(request)["id"]); err != nil {
		log.Printf("%s\n", err)
		return
	}
	if err := c.Get(connection, c.Id); err != nil {
		log.Printf("%s\n", err)
		return
	}

	if err := c.Delete(connection); err != nil {
		log.Printf("%s\n", err)
		return
	}

	http.Redirect(writer, request, "/admin", 302)
	return
}

func adminManufacturerAdd(writer http.ResponseWriter, request *http.Request) {
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
		var m model.Manufacturer
		m.Name = request.FormValue("name")
		if err := m.Add(connection); err != nil {
			log.Printf("%s\n", err)
			return
		}
		http.Redirect(writer, request, "/admin", 302)
		return
	}

	tmpl := template.Must(template.ParseFiles("templates/admin/addForeign.html"))
	if err := tmpl.Execute(writer, nil); err != nil {
		log.Printf("%s\n", err)
	}
	return
}

func adminManufacturerEdit(writer http.ResponseWriter, request *http.Request) {
	session, err := sessionStore.Get(request, adminSession)
	if err != nil {
		log.Printf("Error in getting session: %s\n", err)
		return
	}
	if !isAuthAdmin(session) {
		http.Redirect(writer, request, "/admin/login", 302)
		return
	}
	var m model.Manufacturer
	if m.Id, err = strconv.Atoi(mux.Vars(request)["id"]); err != nil {
		log.Printf("%s\n", err)
		return
	}
	if err := m.Get(connection, m.Id); err != nil {
		log.Printf("%s\n", err)
		return
	}

	if request.Method == http.MethodPost {
		m.Name = request.FormValue("name")
		if err := m.Edit(connection); err != nil {
			log.Printf("%s\n", err)
			return
		}

		http.Redirect(writer, request, "/admin", 302)
		return
	}

	tmpl := template.Must(template.ParseFiles("templates/admin/editForeign.html"))
	if err := tmpl.Execute(writer, m); err != nil {
		log.Printf("%s\n", err)
	}
	return
}

func adminManufacturerDelete(writer http.ResponseWriter, request *http.Request) {
	session, err := sessionStore.Get(request, adminSession)
	if err != nil {
		log.Printf("Error in getting session: %s\n", err)
		return
	}
	if !isAuthAdmin(session) {
		http.Redirect(writer, request, "/admin/login", 302)
		return
	}
	var m model.Manufacturer
	if m.Id, err = strconv.Atoi(mux.Vars(request)["id"]); err != nil {
		log.Printf("%s\n", err)
		return
	}
	if err := m.Get(connection, m.Id); err != nil {
		log.Printf("%s\n", err)
		return
	}

	if err := m.Delete(connection); err != nil {
		log.Printf("%s\n", err)
		return
	}

	http.Redirect(writer, request, "/admin", 302)
	return
}

func cancelOrderAdmin(writer http.ResponseWriter, request *http.Request) {
	session, err := sessionStore.Get(request, adminSession)
	if err != nil {
		log.Printf("Error with session in admin routes cancelOrderAdmin: %s\n", err)
		return
	}
	if !isAuthAdmin(session) {
		http.Redirect(writer, request, "/admin/login", 302)
		return
	}

	var order model.Order

	order.Id, _ = strconv.Atoi(mux.Vars(request)["id"])
	order.Cancel(connection, 0)
	http.Redirect(writer, request, "/admin", 302)
	return
}

func checkOrderAdmin(writer http.ResponseWriter, request *http.Request) {
	session, err := sessionStore.Get(request, adminSession)
	if err != nil {
		log.Printf("Error with session in admin routes cancelOrderAdmin: %s\n", err)
		return
	}
	if !isAuthAdmin(session) {
		http.Redirect(writer, request, "/admin/login", 302)
		return
	}

	var order model.Order

	order.Id, _ = strconv.Atoi(mux.Vars(request)["id"])
	order.Check(connection, 0)
	http.Redirect(writer, request, "/admin", 302)
	return
}
