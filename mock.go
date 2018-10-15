package main

import (
	"fmt"
	"net/http"
)

func rootHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "https://avatars3.githubusercontent.com/u/32389251?v=4")
}

func main() {
	http.HandleFunc("/", rootHandler)

	fmt.Println("Server is listening...")
	// http.ListenAndServeTLS(":8182", "server.crt", "server.key", nil)
	http.ListenAndServe(":8182", nil)
}
