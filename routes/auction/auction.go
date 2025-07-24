package auction

import (
	"car-auction/models/lot"
	"car-auction/websocket/auction"
	"encoding/json"
	"net/http"
)

type Env struct {
	AuctionWebsocketServer *auction.Server
	AuctionLotRepository   *lot.Repository
}

var env Env

func SetupEnv(_env Env) {
	env = _env
}

func HandleGetAllLots(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	responseEncoder := json.NewEncoder(w)

	responseEncoder.Encode(env.AuctionLotRepository.FindAll())
}

func HandleAddLot(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var requestData lot.Model

	requestDecoder := json.NewDecoder(r.Body)
	requestDecoder.Decode(&requestData)

	responseEncoder := json.NewEncoder(w)

	requestData.Picture = "https://imageio.forbes.com/specials-images/imageserve/5d35eacaf1176b0008974b54/0x0.jpg?format=jpg&crop=4560,2565,x790,y784,safe&height=900&width=1600&fit=bounds"

	if len(requestData.Name) == 0 {
		w.WriteHeader(403)
		responseEncoder.Encode(map[string]any{
			"message": `field "name" must have at least 1 character`,
		})
		return
	}

	if len(requestData.Description) == 0 {
		w.WriteHeader(403)
		responseEncoder.Encode(map[string]any{
			"message": `field "description" must have at least 50 characters`,
		})
		return
	}

	if err := env.AuctionLotRepository.Create(&requestData); err != nil {
		w.WriteHeader(500)
		return
	}

	env.AuctionWebsocketServer.WriteJSON("lot.create", requestData)

	responseEncoder.Encode(requestData)
}
