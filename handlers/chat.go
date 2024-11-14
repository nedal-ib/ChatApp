package handlers

import (
	"chatApp/backend/encryption"
	"chatApp/backend/models"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

func ChatHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		http.Error(w, "Could not upgrade to WebSocket", http.StatusBadRequest)
		return
	}
	defer conn.Close()

	log.Println("New WebSocket connection established")

	for {
		var msg models.Message
		if err := conn.ReadJSON(&msg); err != nil {
			if websocket.IsCloseError(err, websocket.CloseNormalClosure, websocket.CloseGoingAway) {
				log.Println("WebSocket connection closed normally:", err)
				break
			}
			log.Println("Error reading JSON from WebSocket:", err)
			break
		}

		recipientPubKey, err := encryption.GetPublicKeyForUser(msg.RecipientID)
		if err != nil {
			log.Printf("Public key not found for recipient ID %s: %v", msg.RecipientID, err)
			continue
		}

		encryptedContent, err := encryption.EncryptMessage(recipientPubKey, []byte(msg.Content))
		if err != nil {
			log.Printf("Error encrypting message for recipient ID %s: %v", msg.RecipientID, err)
			continue
		}

		msg.Content = encryptedContent

		if err := conn.WriteJSON(msg); err != nil {
			log.Println("Error sending encrypted message over WebSocket:", err)
			break
		}
	}
	log.Println("WebSocket connection closed")
}
