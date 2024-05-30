package routes

import (
	"fmt"
	"net/http"
	"os"
)

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "https://account.withings.com/oauth2_user/authorize2?response_type=code&client_id="+os.Getenv("CLIENTID")+"&scope=user.info,user.metrics,user.activity&state=1&redirect_uri=https://localhost:8080/callback", http.StatusMovedPermanently)
	fmt.Println("callback")
}
