package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/gorilla/websocket"
)

const convertAPIURL = "http://142.93.174.191:3000/api/convert"

func rootHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.Error(w, "Wrong URL", http.StatusNotFound)
		return
	}
	http.ServeFile(w, r, "index.html")
}

func wsHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := websocket.Upgrade(w, r, w.Header(), 1024, 1024)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	go sockConnectionHandler(conn)
}

func sockConnectionHandler(conn *websocket.Conn) {
	for {
		msgType, data, err := conn.ReadMessage()
		if err != nil {
			fmt.Println(err)
			return
		}
		if msgType != websocket.BinaryMessage {
			fmt.Println("Invalid message type.")
			continue
		}

		convertedImageURL, err := getConvertedImageURL(data)
		if err != nil {
			fmt.Println(err)
			return
		}
		conn.WriteMessage(websocket.TextMessage, convertedImageURL)
	}
}

func getConvertedImageURL(data []byte) ([]byte, error) {
	r := bytes.NewReader(data)
	resp, err := http.Post(convertAPIURL, "image/png", r)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return body, nil
}

func main() {
	http.HandleFunc("/ws", wsHandler)
	http.HandleFunc("/", rootHandler)

	fmt.Println("Server is listening...")
	http.ListenAndServe(":3000", nil)
}
