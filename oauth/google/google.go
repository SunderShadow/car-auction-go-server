package google

import (
	"encoding/json"
	"net/http"
	"net/url"
	"os"
	"strings"
)

type ExchangeTokenResponse struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   int64  `json:"expires_in"`
	IdToken     string `json:"id_token"`
	Scope       string `json:"scope"`
	TokenType   string `json:"token_type"`
}

type UserInfo struct {
	Email         string `json:"email"`
	EmailVerified bool   `json:"email_verified"`
	FamilyName    string `json:"family_name"`
	GivenName     string `json:"given_name"`
	Name          string `json:"name"`
	Picture       string `json:"picture"`
	Sub           string `json:"sub"`
}

type Env struct {
	AuthURL          string
	ClientID         string
	RedirectURL      string
	TokenExchangeURL string
	UserInfoURL      string
}

type Account struct {
	env *Env
}

func NewAuthenticator(env *Env) *Account {
	return &Account{env}
}

func (auth *Account) RedirectURL() string {
	authURL := auth.env.AuthURL
	authURL += "?client_id=" + auth.env.ClientID
	authURL += "&redirect_uri=" + auth.env.RedirectURL
	authURL += "&response_type=code"
	authURL += "&scope=https://www.googleapis.com/auth/userinfo.email https://www.googleapis.com/auth/userinfo.profile"

	return authURL
}

func (auth *Account) ExchangeToken(code string) (*ExchangeTokenResponse, error) {
	requestData := "client_id=" + os.Getenv("GOOGLE_CLIENT_ID")
	requestData += "&client_secret=" + os.Getenv("GOOGLE_CLIENT_SECRET")
	requestData += "&code=" + code
	requestData += "&grant_type=authorization_code"
	requestData += "&redirect_uri=" + url.QueryEscape(os.Getenv("GOOGLE_OAUTH_REDIRECT_URL"))

	response, err := http.Post(auth.env.TokenExchangeURL, "application/x-www-form-urlencoded", strings.NewReader(requestData))

	if err != nil {
		return nil, err
	}

	data := new(ExchangeTokenResponse)
	err = json.NewDecoder(response.Body).Decode(data)

	if err != nil {
		return nil, err
	}

	return data, nil
}

func (auth *Account) UserInfo(accessToken string) (*UserInfo, error) {
	response, err := http.Get(auth.env.UserInfoURL + "?access_token=" + accessToken)

	if err != nil {
		return nil, err
	}

	data := new(UserInfo)

	err = json.NewDecoder(response.Body).Decode(data)

	if err != nil {
		return nil, err
	}

	return data, nil
}
