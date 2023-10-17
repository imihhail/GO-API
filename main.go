package main

import (
	"fmt"
	"groupie/functions"
	"net/http"
)

func main() {
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	http.HandleFunc("/home", functions.Homepage)
	http.HandleFunc("/locations", Single_Page)
	http.HandleFunc("/dates", Single_Page)
	http.HandleFunc("/artists", Single_Page)
	http.HandleFunc("/", functions.NotFoundHandler)

	fmt.Println("Server is running at: http://localhost:8080/home")
	http.ListenAndServe("0.0.0.0:8080", nil)
}
