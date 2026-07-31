package http_h

import (
	"encoding/json"
	"log"
	"net/http"
	"ride-sharing/services/trip-service/internal/domain"
	json_utils "ride-sharing/shared/json"
	"ride-sharing/shared/types"
)

type TripHandler struct {
	svc domain.TripService
}

func NewTripHandler(s domain.TripService) *TripHandler {
	return &TripHandler{
		svc: s,
	}
}

type previewTripRequest struct {
	UserID      string           `json:"userID"`
	Pickup      types.Coordinate `json:"pickup"`
	Destination types.Coordinate `json:"destination"`
}

func (h *TripHandler) HandleCreateTrip(w http.ResponseWriter, r *http.Request) {
	var preview previewTripRequest
	if err := json.NewDecoder(r.Body).Decode(&preview); err != nil {
		http.Error(w, "TRIP-SERVICE: Failed to parse body from request", http.StatusBadRequest)
		return
	}

	routes, err := h.svc.GetRoute(r.Context(), &preview.Pickup, &preview.Destination)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := json_utils.WriteJSON(w, http.StatusAccepted, routes); err != nil {
		log.Print(err)
	}
}
