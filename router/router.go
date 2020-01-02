package router

import (
	"avtoSell/model"
	"database/sql"
	"encoding/gob"
	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"os"
)

var connection *sql.DB

var sessionStore = sessions.NewCookieStore([]byte("SECRET-KEY"))

const (
	userSession  string = "user-session"
	adminSession string = "admin-session"
)

func Init(db *sql.DB) *mux.Router {
	connection = db
	sessionStore.MaxAge(0)

	gob.Register(&model.User{})

	router := mux.NewRouter()

	// Serve static files
	router.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	router.PathPrefix("/upload-images/").Handler(http.StripPrefix("/upload-images/", http.FileServer(http.Dir("upload-images"))))

	// Site routes
	router.HandleFunc("/", index)

	router.HandleFunc("/news", news)
	router.HandleFunc("/news/{id:[0-9]+}", post)

	router.HandleFunc("/catalog", catalog)
	router.HandleFunc("/product/{id:[0-9]+}", product)

	// User routes
	router.HandleFunc("/login", signIn).Methods(http.MethodGet, http.MethodPost)
	router.HandleFunc("/registration", signUp).Methods(http.MethodGet, http.MethodPost)
	router.HandleFunc("/cabinet", cabinet)
	router.HandleFunc("/cabinet/exit", exitUser)

	// Admin routes
	router.HandleFunc("/admin/login", adminLogin).Methods(http.MethodGet, http.MethodPost)
	router.HandleFunc("/admin", admin)
	router.HandleFunc("/admin/exit", adminExit)

	router.HandleFunc("/admin/news/add", adminNewsAdd).Methods(http.MethodGet, http.MethodPost)
	router.HandleFunc("/admin/news/{id:[0-9]+}/edit", adminNewsEdit).Methods(http.MethodGet, http.MethodPost)
	router.HandleFunc("/admin/news/{id:[0-9]+}/delete", adminNewsDelete).Methods(http.MethodGet)

	// Api functions
	router.HandleFunc("/api/checkLogin", checkLogin).Methods(http.MethodPost)
	router.HandleFunc("/api/checkEmail", checkEmail).Methods(http.MethodPost)
	router.HandleFunc("/api/checkPhone", checkPhone).Methods(http.MethodPost)
	router.HandleFunc("/api/news/all", ApiNewsAll)
	return router
}

func isAuthUser(session *sessions.Session) bool {
	if _, ok := session.Values["user"]; ok {
		return true
	}
	return false
}

func isAuthAdmin(session *sessions.Session) bool {
	if _, ok := session.Values["admin"]; ok {
		return true
	}
	return false
}

func UploadImages(file multipart.File) (string, error) {
	tempFile, err := ioutil.TempFile("upload-images", "upload-*.png")
	if err != nil {
		return "", err
	}
	defer tempFile.Close()

	fileBytes, err := ioutil.ReadAll(file)
	if err != nil {
		return "", err
	}

	_, err = tempFile.Write(fileBytes)
	if err != nil {
		return "", err
	}
	return tempFile.Name(), nil
}

func DeleteImages(src string) error {
	err := os.Remove(src)
	if err != nil && !os.IsNotExist(err) {
		return err
	}
	return nil
}
