package main

import (
	"car-auction/models/lot"
	_ "modernc.org/sqlite"

	"github.com/joho/godotenv"
	"github.com/rs/cors"

	auctoinWebsocket "car-auction/websocket/auction"

	"car-auction/models/user"
	oauthGoogle "car-auction/models/user/oauth/google"

	googleOauthHelper "car-auction/oauth/google"
	auctionRoutes "car-auction/routes/auction"
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
		Google *googleOauthHelper.Account
	}
	Repositories struct {
		UserRepository        *user.Repository
		GoogleOauthRepository *oauthGoogle.Repository
		AuctionLotRepository  *lot.Repository
	}
	AuctionWebsocketServer *auctoinWebsocket.Server
}

var env Env

func main() {
	initEnvironment()

	httpMux := http.NewServeMux()

	c := cors.New(cors.Options{
		AllowedMethods:   []string{"GET", "POST"},
		AllowCredentials: true,
		Debug:            false,
		AllowedOrigins:   []string{"*"},
	})

	initHttpRoutes(httpMux)
	initWebsocketRoutes(httpMux)

	fmt.Println("Server listening on http://localhost:" + os.Getenv("HTTP_SERVER_PORT"))

	err := http.ListenAndServe(":"+os.Getenv("HTTP_SERVER_PORT"), c.Handler(httpMux))

	if err != nil {
		log.Fatal(err)
	}
}

func initWebsocketRoutes(httpMux *http.ServeMux) {
	httpMux.HandleFunc("/auction", env.AuctionWebsocketServer.ServeWebsocket)
}

func initHttpRoutes(mux *http.ServeMux) {
	userRoutes.SetupEnv(userRoutes.Env{
		GoogleOauth:           env.Oauth.Google,
		UserRepository:        env.Repositories.UserRepository,
		GoogleOauthRepository: env.Repositories.GoogleOauthRepository,
	})

	mux.HandleFunc("GET /user/oauth/google", userRoutes.HandleAuthByGoogleOauth)
	mux.HandleFunc("GET /oauth/google", userRoutes.HandleFinishGoogleAuth)
	mux.HandleFunc("GET /user/bids", userRoutes.HandleGetUserBids)

	auctionRoutes.SetupEnv(auctionRoutes.Env{
		AuctionLotRepository:   env.Repositories.AuctionLotRepository,
		AuctionWebsocketServer: env.AuctionWebsocketServer,
	})

	mux.HandleFunc("GET /auction/lot/all", auctionRoutes.HandleGetAllLots)
	mux.HandleFunc("PUT /auction/lot/add", auctionRoutes.HandleAddLot)
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
	env.Oauth.Google = googleOauthHelper.NewAuthenticator(&googleOauthHelper.Env{
		AuthURL:          "https://accounts.google.com/o/oauth2/auth",
		ClientID:         os.Getenv("GOOGLE_CLIENT_ID"),
		RedirectURL:      os.Getenv("GOOGLE_OAUTH_REDIRECT_URL"),
		TokenExchangeURL: "https://oauth2.googleapis.com/token",
		UserInfoURL:      "https://www.googleapis.com/oauth2/v3/userinfo",
	})

	env.AuctionWebsocketServer = auctoinWebsocket.NewServer()

	initRepositories()
}

func initRepositories() {
	env.Repositories.UserRepository = user.NewRepository(env.DB)
	if err := env.Repositories.UserRepository.CreateTable(); err != nil {
		log.Fatal(err)
	}

	env.Repositories.GoogleOauthRepository = oauthGoogle.NewRepository(env.DB)
	if err := env.Repositories.GoogleOauthRepository.CreateTable(); err != nil {
		log.Fatal(err)
	}

	env.Repositories.AuctionLotRepository = lot.NewRepository(env.DB)
	if err := env.Repositories.AuctionLotRepository.CreateTable(); err != nil {
		log.Fatal(err)
	}
}
