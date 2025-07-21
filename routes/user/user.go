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
		w.WriteHeader(403)

		responseEncoder.Encode(map[string]any{
			"message": "Unauthenticated",
		})

		return
	}

	exchangeToken, err := env.GoogleOauth.ExchangeToken(requestData.Get("code"))

	if err != nil {
		w.WriteHeader(403)

		responseEncoder.Encode(map[string]any{
			"message": "Unauthenticated",
		})

		return
	}

	userInfo, userInfoErr := env.GoogleOauth.UserInfo(exchangeToken.AccessToken)

	if userInfoErr != nil {
		w.WriteHeader(403)

		responseEncoder.Encode(map[string]any{
			"message": "Unauthenticated",
		})

		return
	}

	googleOauthModel := env.GoogleOauthRepository.FindByGoogleUserId(userInfo.Sub)

	var userModel *user.Model

	if googleOauthModel != nil {
		userModel = env.UserRepository.FindById(googleOauthModel.UserId)

		if userModel != nil {
			responseEncoder.Encode(userInfo)
			return
		}
	}

	userModel = new(user.Model)

	userModel.Name = userInfo.Name
	userModel.Picture = userInfo.Picture

	err = env.UserRepository.Register(userModel)

	if err != nil {
		w.WriteHeader(500)

		responseEncoder.Encode(map[string]any{
			"message": "Server error",
		})

		return
	}

	if googleOauthModel == nil {
		googleOauthModel = new(googleOauth.Model)
		googleOauthModel.AccessToken = exchangeToken.AccessToken
		googleOauthModel.AccessTokenExpiresIn = exchangeToken.ExpiresIn
		googleOauthModel.GoogleUserId = userInfo.Sub

		err := env.GoogleOauthRepository.Register(userModel, googleOauthModel)

		if err != nil {
			w.WriteHeader(500)

			responseEncoder.Encode(map[string]any{
				"message": "Server error",
			})

			return
		}
	}

	responseEncoder.Encode(userInfo)
}
