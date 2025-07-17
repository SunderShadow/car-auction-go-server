package user

import (
	"encoding/json"
	"net/http"
)

func HandleGetUserBids(w http.ResponseWriter, r *http.Request) {
	jsonReponse := json.NewEncoder(w)

	jsonReponse.Encode([]map[string]any{
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
