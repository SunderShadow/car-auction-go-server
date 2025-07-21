package main

import (
	"car-auction/models/user"
	"github.com/joho/godotenv"

	_ "modernc.org/sqlite"

	"car-auction/oauth/google"
	"car-auction/routes/auction"
	userRoutes "car-auction/routes/user"

	"database/sql"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"path"
)

type Env struct {
	DB    *sql.DB
	Oauth struct {
		Google *google.Account
	}
	Repositories struct {
		UserRepository *user.Repository
	}
}

var env Env

func main() {
	initEnvironment()

	httpMux := http.NewServeMux()

	initHttpRoutes(httpMux)

	fmt.Print("Server listening on http://localhost:" + os.Getenv("HTTP_SERVER_PORT"))

	err := http.ListenAndServe(":"+os.Getenv("HTTP_SERVER_PORT"), httpMux)

	if err != nil {
		log.Fatal(err)
	}
}

func initHttpRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/auction/lots", auction.HandleGetActiveLots)

	userRoutes.SetupEnv(userRoutes.Env{
		GoogleOauth:    env.Oauth.Google,
		UserRepository: env.Repositories.UserRepository,
	})

	mux.HandleFunc("/user/oauth/google", userRoutes.HandleAuthByGoogleOauth)
	mux.HandleFunc("/oauth/google", userRoutes.HandleFinishGoogleAuth)

	mux.HandleFunc("/user/bids", userRoutes.HandleGetUserBids)
}

func initDatabase() *sql.DB {
	filepath := os.Getenv("SQLITE_DATABASE_FILE")

	if filepath[0] != '/' {
		filepath = "/" + filepath
	}

	filepath = "." + filepath

	if _, err := os.Stat(filepath); errors.Is(err, os.ErrNotExist) {
		dir := path.Dir(filepath)

		if dir != "." {
			err := os.MkdirAll(dir, 0750)

			if err != nil {
				log.Fatal("Unable to create database file " + filepath)
			}
		}

		file, dbCreateErr := os.Create(filepath)

		file.Chmod(0775)
		if dbCreateErr != nil {
			log.Fatal("Unable to create database file " + filepath)
		}
	}

	db, err := sql.Open("sqlite", filepath)

	if err != nil {
		log.Println("Unable to connect to " + filepath + " file")
		log.Fatal(err)
	}

	if err := db.Ping(); err != nil {
		log.Println("Unable to connect to " + filepath + " file")
		log.Fatal(err)
	}

	return db
}

func initEnvironment() {
	err := godotenv.Load()

	if err != nil {
		log.Fatal("Error loading .env file")
	}

	env.DB = initDatabase()
	env.Oauth.Google = google.NewAuthenticator(&google.Env{
		AuthURL:          "https://accounts.google.com/o/oauth2/auth",
		ClientID:         os.Getenv("GOOGLE_CLIENT_ID"),
		RedirectURL:      os.Getenv("GOOGLE_OAUTH_REDIRECT_URL"),
		TokenExchangeURL: "https://oauth2.googleapis.com/token",
		UserInfoURL:      "https://www.googleapis.com/oauth2/v3/userinfo",
	})

	initRepositories()
}

func initRepositories() {
	env.Repositories.UserRepository = user.NewRepository(env.DB)
	if err := env.Repositories.UserRepository.CreateTable(); err != nil {
		log.Fatal(err)
	}
}
