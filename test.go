package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/gorilla/websocket"
)

const transformAPIURL = "http://localhost:8182/"

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

	go connectionHandler(conn)
}

func connectionHandler(conn *websocket.Conn) {
	for {
		msgType, data, err := conn.ReadMessage()
		if err != nil {
			fmt.Println(err)
			conn.WriteMessage(websocket.TextMessage, []byte(err.Error()))
			return
		}
		conn.WriteMessage(websocket.TextMessage, []byte("Server received data."))
		if msgType != websocket.BinaryMessage {
			conn.WriteMessage(websocket.TextMessage, []byte("Invalid data format."))
			continue
		}

		transformedImageURL, err := getTransformedImageURL(data)
		if err != nil {
			fmt.Println(err)
			conn.WriteMessage(websocket.TextMessage, []byte(err.Error()))
			return
		}
		conn.WriteMessage(websocket.TextMessage, transformedImageURL)
	}
}

func getTransformedImageURL(data []byte) ([]byte, error) {
	r := bytes.NewReader(data)
	resp, err := http.Post(transformAPIURL, "image/png", r)
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
	http.ListenAndServe("localhost:8181", nil)
}
