package functions

import (
	"fmt"
	"io/ioutil"
	"net/http"
)

func NotFoundHandler(w http.ResponseWriter, r *http.Request) {
	// Set the response status to 404 Not Found
	w.WriteHeader(http.StatusNotFound)

	// Read the contents of the 404.html file
	data, err := ioutil.ReadFile("template/errorpage.html")
	if err != nil {
		http.Error(w, "Error reading 404 page", http.StatusInternalServerError)
		return
	}

	// Write the contents of the 404.html file as the response body
	_, err = w.Write(data)
	if err != nil {
		fmt.Println("Error writing 404 page:", err)
	}
}
