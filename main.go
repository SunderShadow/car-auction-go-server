package main

import (
	"car-auction/routes/auction"
	"car-auction/routes/user"
	"fmt"
	"github.com/joho/godotenv"
	"log"
	"net/http"
	"os"
)

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
	mux.HandleFunc("/user/bids", user.HandleGetUserBids)
}

func initEnvironment() {
	err := godotenv.Load()

	if err != nil {
		log.Fatal("Error loading .env file")
	}
}
