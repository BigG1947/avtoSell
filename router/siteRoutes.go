package router

import (
	"html/template"
	"log"
	"net/http"
)

func index(writer http.ResponseWriter, request *http.Request) {
	tmpl := template.Must(template.ParseFiles("templates/index.html"))
	err := tmpl.Execute(writer, nil)
	if err != nil{
		log.Printf("Error in site routes 'index': %s\n", err)
	}
	return
}


func news(writer http.ResponseWriter, request *http.Request) {
	tmpl := template.Must(template.ParseFiles("templates/news.html"))
	err := tmpl.Execute(writer, map[string]interface{}{
		"repeat": []int{1,2,3,4,5},
	})
	if err != nil{
		log.Printf("Error in site routes 'news': %s\n", err)
	}
	return
}

func post(writer http.ResponseWriter, request *http.Request) {
	tmpl := template.Must(template.ParseFiles("templates/post.html"))
	err := tmpl.Execute(writer, nil)
	if err != nil{
		log.Printf("Error in site routes 'post': %s\n", err)
	}
	return
}

func product(writer http.ResponseWriter, request *http.Request) {
	tmpl := template.Must(template.ParseFiles("templates/product.html"))
	err := tmpl.Execute(writer, nil)
	if err != nil{
		log.Printf("Error in site routes 'product': %s\n", err)
	}
	return
}

func catalog(writer http.ResponseWriter, request *http.Request) {
	tmpl := template.Must(template.ParseFiles("templates/catalog.html"))
	err := tmpl.Execute(writer, map[string]interface{}{
		"repeat": []int{1,2,3,4,5,6,7,8,9,10,11,12,13,14,15,16},
	})
	if err != nil{
		log.Printf("Error in site routes 'catalog': %s\n", err)
	}
	return
}