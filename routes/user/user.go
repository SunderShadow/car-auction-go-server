package user

import (
	"car-auction/oauth/google"
	"encoding/json"
	"net/http"
)

type Env struct {
	GoogleOauth *google.Account
}

var env Env

func SetupEnv(_env Env) {
	env = _env
}

func HandleGetUserBids(w http.ResponseWriter, r *http.Request) {
	jsonResponse := json.NewEncoder(w)

	jsonResponse.Encode([]map[string]any{
		{
			"Id":  0,
			"Bid": "50000",
		},
		{
			"Id":  1,
			"Bid": "35600",
		},
	})
}

func HandleAuthByGoogleOauth(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Location", env.GoogleOauth.RedirectURL())

	w.WriteHeader(302)
}

func HandleFinishGoogleAuth(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	responseEncoder := json.NewEncoder(w)

	requestData := r.URL.Query()

	if requestData.Has("error") {
		responseEncoder.Encode(map[string]any{
			"message": "Unauthenticated",
		})

		w.WriteHeader(403)
		return
	}

	data, err := env.GoogleOauth.ExchangeToken(requestData.Get("code"))

	if err != nil {
		responseEncoder.Encode(map[string]any{
			"message": "Unauthenticated",
		})

		w.WriteHeader(403)
		return
	}

	userInfo, _ := env.GoogleOauth.UserInfo(data.AccessToken)

	responseEncoder.Encode(userInfo)
}
