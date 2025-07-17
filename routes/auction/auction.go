package auction

import (
	"encoding/json"
	"net/http"
)

func HandleGetActiveLots(w http.ResponseWriter, r *http.Request) {
	jsonReponse := json.NewEncoder(w)

	jsonReponse.Encode([]map[string]any{
		{
			"Id":          0,
			"Title":       "Audi 3Xs",
			"Bid":         "50000",
			"Description": "Some really important description",
		},
		{
			"Id":          1,
			"Title":       "Audi 999Xs",
			"Bid":         "80000",
			"Description": "Lorem ipsum dolor sit amet",
		},
		{
			"Id":          2,
			"Title":       "Audi 999Xs",
			"Bid":         "80000",
			"Description": "Lorem ipsum dolor sit amet",
		},
	})
}
