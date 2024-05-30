package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/awslabs/aws-lambda-go-api-proxy/core"
	"github.com/awslabs/aws-lambda-go-api-proxy/gorillamux"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	routes "github.com/ojparkinson/withings-display/routes"
)

var gorillaLambda *gorillamux.GorillaMuxAdapter
var r *mux.Router

func init() {
	godotenv.Load("../dev.env")
	r = mux.NewRouter()

	r.HandleFunc("/login", routes.LoginHandler)
	r.HandleFunc("/welcome", routes.WelcomeHandler)
	r.HandleFunc("/callback", routes.CallbackHandler)
	r.HandleFunc("/", HomeHandler)
	http.Handle("/", r)

	if os.Getenv("ENV") != "development" {
		gorillaLambda = gorillamux.New(r)
	}
}

var access_token string

func main() {

	if os.Getenv("ENV") == "development" {
		fmt.Println("Server Started 🚀")
		log.Fatal(http.ListenAndServeTLS("localhost:8080", "../cert.pem", "../key.pem", r))
	} else {
		lambda.Start(Handler)
	}

}

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Hello world"))
	fmt.Println("hello world")
}

func Handler(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	r, err := gorillaLambda.ProxyWithContext(ctx, *core.NewSwitchableAPIGatewayRequestV1(&req))
	return *r.Version1(), err
}
