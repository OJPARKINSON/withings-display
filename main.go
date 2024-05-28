package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
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
	fmt.Println("Server Started 🚀")
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

	response, err := http.Get("https://wbsapi.withings.net/v2/oauth2?action=requesttoken&grant_type=authorization_code&client_id=" + os.Getenv("CLIENTID") + "&client_secret=" + os.Getenv("CLIENTSECRET") + "&code=" + params.Get("code") + "&redirect_uri=https://localhost:8080/callback")

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

	measurements := getMeasurements(access_token)
	fmt.Println(measurements)

	w.WriteHeader(http.StatusOK)
}

type responseMeasurements struct {
	Status int `json:"status"`
	Body   struct {
		Updatetime  int    `json:"updatetime"`
		Timezone    string `json:"timezone"`
		Measuregrps []struct {
			Grpid        int    `json:"grpid"`
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

func getMeasurements(accessToken string) []Measurements {
	req, err := http.NewRequest(http.MethodGet, "https://wbsapi.withings.net/measure?action=getmeas&meastypes=1,6,8,88&category=1&startdate=1709149524&enddate=1716921924", nil)
	if err != nil {
		log.Fatalln("Failed to create req")
	}

	req.Header.Set("authorization", "Bearer "+accessToken)

	c := &http.Client{}
	response, err := c.Do(req)
	if err != nil {
		log.Fatalln("Failed to fetch measurement")
	}

	responseData, err := io.ReadAll(response.Body)
	if err != nil {
		log.Fatal(err)
	}

	var responseMeasurements responseMeasurements

	err = json.Unmarshal(responseData, &responseMeasurements)
	if err != nil {
		fmt.Println("--------", err)
	}

	measurements := parseMeasurements(responseMeasurements)

	return measurements
}

type Measurements struct {
	ID       int `json:"id"`
	Date     int `json:"date"`
	Weight   int `json:"weight"`   // 1
	FatMass  int `json:"fatmass"`  // 6
	BodyFat  int `json:"bodyfat"`  // 8
	BoneMass int `json:"bonemass"` //88
}

func parseMeasurements(oldslice responseMeasurements) []Measurements {

	var measurements []Measurements

	for _, m := range oldslice.Body.Measuregrps {
		date := m.Date
		id := m.Grpid

		var weight int
		var fatmass int
		var BodyFat int
		var BoneMass int
		for _, p := range m.Measures {

			switch p.Type {
			case 1:
				weight = p.Value
			case 6:
				fatmass = p.Value
			case 8:
				BodyFat = p.Value
			case 88:
				BoneMass = p.Value
			}
		}

		measurements = append(measurements, Measurements{ID: id, Date: date, Weight: weight, FatMass: fatmass, BodyFat: BodyFat, BoneMass: BoneMass})
	}

	return measurements
}
