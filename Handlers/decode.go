package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func DecodeData(data any, w http.ResponseWriter, r *http.Request) any {
	decoder := json.NewDecoder(r.Body)
	decoder.Decode(data)
	err := Validation(data)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error %s", err), http.StatusBadRequest)
	}
	return data
}
