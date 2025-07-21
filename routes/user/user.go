package user

import (
	"car-auction/models/user"
	googleOauth "car-auction/models/user/oauth/google"
	"car-auction/oauth/google"

	"encoding/json"
	"net/http"
)

type Env struct {
	GoogleOauth *google.Account

	UserRepository        *user.Repository
	GoogleOauthRepository *googleOauth.Repository
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

	exchangeToken, err := env.GoogleOauth.ExchangeToken(requestData.Get("code"))

	if err != nil {
		responseEncoder.Encode(map[string]any{
			"message": "Unauthenticated",
		})

		w.WriteHeader(403)
		return
	}

	userInfo, _ := env.GoogleOauth.UserInfo(exchangeToken.AccessToken)

	responseEncoder.Encode(userInfo)

	userModel := new(user.Model)

	userModel.Name = userInfo.Name
	userModel.Picture = userInfo.Picture

	env.UserRepository.Register(userModel)

	googleOauthModel := new(googleOauth.Model)
	googleOauthModel.AccessToken = exchangeToken.AccessToken
	googleOauthModel.AccessTokenExpiresIn = exchangeToken.ExpiresIn

	env.GoogleOauthRepository.Register(userModel, googleOauthModel)
}
