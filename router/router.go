package router

import (
	"avtoSell/model"
	"database/sql"
	"encoding/gob"
	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"net/http"
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

	// Api functions
	router.HandleFunc("/api/checkLogin", checkLogin).Methods(http.MethodPost)
	router.HandleFunc("/api/checkEmail", checkEmail).Methods(http.MethodPost)
	router.HandleFunc("/api/checkPhone", checkPhone).Methods(http.MethodPost)
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
