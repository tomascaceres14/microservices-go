package main

import (
	"log"
	"net/http"
	"ride-sharing/services/api-gateway/grpc_clients"
	"ride-sharing/shared/contracts"
	pb "ride-sharing/shared/proto/driver"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func handleRiderWebSocket(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("WebSocket upgrade failed: %s\n", err)
		return
	}

	defer conn.Close()

	userID := r.URL.Query().Get("userID")
	if userID == "" {
		log.Printf("No user ID provided\n")
		return
	}

	for {
		_, msg, err := conn.ReadMessage()
		if err != nil {
			log.Printf("Error reading WS message: %s", err)
			break
		}

		log.Printf("Message received: %s", msg)
	}
}

func handleDriverWebSocket(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("WebSocket upgrade failed: %s\n", err)
		return
	}

	defer conn.Close()

	driverID := r.URL.Query().Get("userID")
	if driverID == "" {
		log.Printf("No user ID provided\n")
		return
	}

	pkgSlug := r.URL.Query().Get("packageSlug")
	if pkgSlug == "" {
		log.Printf("No package slug provided\n")
		return
	}

	driverService, err := grpc_clients.NewDriverServiceClient()
	if err != nil {
		log.Printf("Error connecting to driver service: %s", err)
		conn.WriteMessage(websocket.TextMessage, []byte(err.Error()))
		return
	}

	// Unregister driver when connection gets closed
	defer func() {
		driverService.Client.UnregisterDriver(r.Context(), &pb.RegisterDriverRequest{
			DriverID:    driverID,
			PackageSlug: pkgSlug,
		})
		driverService.Close()
		log.Printf("Driver unregistered: %s", driverID)
	}()

	req := &pb.RegisterDriverRequest{
		DriverID:    driverID,
		PackageSlug: pkgSlug,
	}

	registerDriverResponse, err := driverService.Client.RegisterDriver(r.Context(), req)
	if err != nil {
		log.Printf("Error registering driver: %s", err)
		return
	}

	msg := contracts.WSMessage{
		Type: "driver.cmd.register",
		Data: registerDriverResponse.Driver,
	}

	if err := conn.WriteJSON(msg); err != nil {
		log.Printf("Error reading message: %s", err)
		return
	}

	// Listen incoming messages
	for {
		_, msg, err := conn.ReadMessage()
		if err != nil {
			log.Printf("Error reading WS message: %s", err)
			break
		}

		log.Printf("Message received: %s", msg)
	}
}
