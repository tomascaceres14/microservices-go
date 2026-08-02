package main

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"ride-sharing/services/api-gateway/grpc_clients"
	"ride-sharing/shared/contracts"
	json_utils "ride-sharing/shared/json"
)

func handleTripReview(w http.ResponseWriter, r *http.Request) {

	var reqBody previewTripReq
	w.Header().Set("Content-Type", "application/json")

	if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	defer r.Body.Close()

	if reqBody.UserID == "" {
		http.Error(w, "userID is required.", http.StatusBadRequest)
		return
	}

	log.Println("Body accepted:", reqBody)

	// TODO: Call trip service.
	jsonBody, _ := json.Marshal(reqBody)
	reader := bytes.NewReader(jsonBody)

	tripService, err := grpc_clients.NewTripServiceClient()
	if err != nil {
		log.Fatal(err)
	}

	defer tripService.Close()

	resp, err := http.Post("http://trip-service:8083/preview", "application/json", reader)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	var trip any
	if err := json.NewDecoder(resp.Body).Decode(&trip); err != nil {
		http.Error(w, "failed to parse JSON data from trip service", http.StatusInternalServerError)
		return
	}

	res := contracts.APIResponse{Data: trip}
	json_utils.WriteJSON(w, http.StatusCreated, res)
}
