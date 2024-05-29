package routes

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
)

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

	w.WriteHeader(http.StatusOK)
	http.Redirect(w, r, "https://localhost:8080/welcome?access_token="+responseBody.Body.Access_token, http.StatusPermanentRedirect)
}
