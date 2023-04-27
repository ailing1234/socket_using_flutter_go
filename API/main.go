package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

// We'll need to define an Upgrader
// this will require a Read and Write buffer size
var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(r *http.Request) bool { return true },
}

// define a reader which will listen for
// new messages being sent to our WebSocket
// endpoint
func reader(conn *websocket.Conn) {
	for {
		// read in a message
		messageType, p, err := conn.ReadMessage()
		if err != nil {
			log.Println(err)
			return
		}
		// print out that message for clarity
		log.Println(string(p))

		// Echo the message back to the client
		if err := conn.WriteMessage(messageType, p); err != nil {
			log.Println(err)
			return
		}

		// Send a message to the client on channel1 after receiving a message
		message := map[string]string{"channel": "channel1", "message": "Hello client on channel1!"}
		jsonMessage, err := json.Marshal(message)
		if err != nil {
			log.Println(err)
			return
		}
		if err := conn.WriteMessage(messageType, jsonMessage); err != nil {
			log.Println(err)
			return
		}
	}
}

func homePage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Home Page")
}

func wsEndpoint(w http.ResponseWriter, r *http.Request) {
	// upgrade this connection to a WebSocket connection
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
	}

	log.Println("Client Connected")
	message := map[string]string{"channel": "notification", "message": "Hello client on notification!"}
	jsonMessage, _ := json.Marshal(message)
	err = ws.WriteMessage(1, jsonMessage)
	if err != nil {
		log.Println(err)
	}

	ticker := time.Tick(30 * time.Second)
	go func() {
		for range ticker {
			message := map[string]string{"channel": "update", "message": "Please update your information."}
			jsonMessage, err := json.Marshal(message)
			if err != nil {
				log.Println(err)
				continue
			}
			if err := ws.WriteMessage(websocket.TextMessage, jsonMessage); err != nil {
				log.Println(err)
				continue
			}
		}
	}()

	// listen indefinitely for new messages coming
	// through on our WebSocket connection
	reader(ws)
}

/*
func setupRoutes() {
	// Set up CORS headers
	cors := handlers.CORS(
		handlers.AllowedOrigins([]string{"http://localhost:57630"}),
		handlers.AllowCredentials(),
		handlers.AllowedHeaders([]string{"X-Requested-With"}),
		handlers.AllowedMethods([]string{"GET", "POST", "OPTIONS"}),
	)

	http.HandleFunc("/", cors.Wrap(http.HandlerFunc(homePage)))
	http.HandleFunc("/ws", cors.Wrap(http.HandlerFunc(wsEndpoint)))
}
*/

func main() {
	fmt.Println("Hello World")
	// setupRoutes()

	// Create a new router
	router := http.NewServeMux()

	// Add your routes here
	router.HandleFunc("/", homePage)
	router.HandleFunc("/ws", wsEndpoint)

	// Wrap the router with CORS middleware
	// corsHandler := handlers.CORS(
	// handlers.AllowedOrigins([]string{"http://localhost:52609/"}),
	// handlers.AllowCredentials(),
	// handlers.AllowedHeaders([]string{"X-Requested-With"}),
	// handlers.AllowedMethods([]string{"GET", "POST", "OPTIONS"}),
	// )(router)

	// Start the server
	http.ListenAndServe(":8080", router)
}
