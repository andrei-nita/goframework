package routes

import (
	fk "github.com/andrei-nita/goframework/framework"
	"log"
	"net/http"
	"strings"
)

func livereload(w http.ResponseWriter, r *http.Request) {
	//upgrade get request to websocket protocol
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	defer ws.Close()
	for {
		//Read Message from a client
		mt, message, err := ws.ReadMessage()
		if err != nil {
			if !strings.Contains(err.Error(), "close 1005") {
				log.Println(err)
			}
			break
		}

		// Ciao is the message from the browser sent when the connection is started
		if string(message) == "Ciao" {
			//
		}

		if <-fk.Reload {
			message = []byte("reload")
			log.Println("🌐 Websocket reloaded")
		}

		//Response message to a client
		err = ws.WriteMessage(mt, message)
		if err != nil {
			log.Println(err)
			break
		}
	}
}
