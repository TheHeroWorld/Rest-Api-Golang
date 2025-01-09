package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/sirupsen/logrus"
)

func DecodeData(data any, w http.ResponseWriter, r *http.Request) any {
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(data)
	if err != nil {
		log.WithFields(logrus.Fields{
			"error": err.Error(),
		}).Error("Error decoding request body")
		http.Error(w, fmt.Sprintf("Error %s", err), http.StatusBadRequest)
		return nil
	}

	err = Validation(data)
	if err != nil {
		log.WithFields(logrus.Fields{
			"error": err.Error(),
		}).Warn("Data validation failed")
		http.Error(w, fmt.Sprintf("Error %s", err), http.StatusBadRequest)
		return nil
	}
	return data
}
