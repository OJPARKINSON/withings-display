package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var access_token string

func main() {
	godotenv.Load("dev.env")

	r := mux.NewRouter()
	r.HandleFunc("/home", HomeHandler)
	r.HandleFunc("/login", LoginHandler)
	r.HandleFunc("/welcome", WelcomeHandler)
	r.HandleFunc("/callback", CallbackHandler)
	http.Handle("/", r)

	log.Fatal(http.ListenAndServeTLS("localhost:8080", "cert.pem", "key.pem", r))
}

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Hello world"))
	fmt.Println("hello world")
}

func CallbackHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println(r.URL)

	params, err := url.ParseQuery(r.URL.RawQuery)

	if err != nil {
		log.Fatalln("couldn't parse url")
	}

	responseCode := params.Get("code")

	fmt.Println(responseCode)

	response, err := http.Get("https://wbsapi.withings.net/v2/oauth2?action=requesttoken&grant_type=authorization_code&client_id=" + os.Getenv("CLIENTID") + "&client_secret=" + os.Getenv("CLIENTSECRET") + "&code=" + responseCode + "&redirect_uri=https://localhost:8080/callback")

	if err != nil {
		log.Fatalln("Auth fetch fail")
	}

	responseData, err := io.ReadAll(response.Body)
	if err != nil {
		log.Fatal(err)
	}

	type Body struct {
		Userid        string `json:"userid"`
		Access_token  string `json:"access_token"`
		Refresh_token string `json:"refresh_token"`
		Scope         string `json:"scope"`
		Expires_in    int    `json:"expires_in"`
		Token_type    string `json:"token_type"`
	}

	type Response struct {
		Status int  `json:"status"`
		Body   Body `json:"body"`
	}

	var responseBody Response

	err = json.Unmarshal(responseData, &responseBody)
	if err != nil {
		fmt.Println("--------", err)
	}

	fmt.Println("resp access: " + responseBody.Body.Access_token)
	fmt.Println("resp refresh: " + responseBody.Body.Refresh_token)

	access_token = responseBody.Body.Access_token

	w.WriteHeader(http.StatusOK)

	http.Redirect(w, r, "https://localhost:8080/welcome", http.StatusMovedPermanently)
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "https://account.withings.com/oauth2_user/authorize2?response_type=code&client_id="+os.Getenv("CLIENTID")+"&scope=user.info,user.metrics,user.activity&state=1&redirect_uri=https://localhost:8080/callback", http.StatusMovedPermanently)
	fmt.Println("callback")
}

func WelcomeHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte(access_token))

	serverAPI := options.ServerAPI(options.ServerAPIVersion1)
	opts := options.Client().ApplyURI(os.Getenv("CONNECTIONSTRING")).SetServerAPIOptions(serverAPI)

	client, err := mongo.Connect(context.TODO(), opts)
	if err != nil {
		panic(err)
	}
	defer func() {
		if err = client.Disconnect(context.TODO()); err != nil {
			panic(err)
		}
	}()

	if err := client.Database("admin").RunCommand(context.TODO(), bson.D{{"ping", 1}}).Err(); err != nil {
		panic(err)
	}
	fmt.Println("Pinged your deployment. You successfully connected to MongoDB!")

	w.WriteHeader(http.StatusOK)
}

type measureStrut struct {
	Status int `json:"status"`
	Body   struct {
		Updatetime  int    `json:"updatetime"`
		Timezone    string `json:"timezone"`
		Measuregrps []struct {
			Grpid        int64  `json:"grpid"`
			Attrib       int    `json:"attrib"`
			Date         int    `json:"date"`
			Created      int    `json:"created"`
			Modified     int    `json:"modified"`
			Category     int    `json:"category"`
			Deviceid     string `json:"deviceid"`
			HashDeviceid string `json:"hash_deviceid"`
			Measures     []struct {
				Value int `json:"value"`
				Type  int `json:"type"`
				Unit  int `json:"unit"`
				Algo  int `json:"algo,omitempty"`
				Fm    int `json:"fm,omitempty"`
			} `json:"measures"`
			Modelid int    `json:"modelid"`
			Model   string `json:"model"`
			Comment any    `json:"comment"`
		} `json:"measuregrps"`
	} `json:"body"`
}
