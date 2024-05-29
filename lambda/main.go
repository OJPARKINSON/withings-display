package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	routes "github.com/ojparkinson/withings-display/routes"
)

var access_token string

func main() {
	godotenv.Load("../dev.env")
	r := mux.NewRouter()
	r.HandleFunc("/home", HomeHandler)
	r.HandleFunc("/login", routes.LoginHandler)
	r.HandleFunc("/welcome", routes.WelcomeHandler)
	r.HandleFunc("/callback", routes.CallbackHandler)
	http.Handle("/", r)

	fmt.Println("Server Started 🚀")
	log.Fatal(http.ListenAndServeTLS("localhost:8080", "../cert.pem", "../key.pem", r))
}

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Hello world"))
	fmt.Println("hello world")
}
