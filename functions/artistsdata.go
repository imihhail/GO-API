package functions

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"net/http"
	"strings"
)

type LocationData struct {
	Index []struct {
		ID        int      `json:"id"`
		Locations []string `json:"locations"`
	} `json:"index"`
}

type ArtistInfo struct {
	Id                 int      `json:"id"`
	Image              string   `json:"image"`
	Name               string   `json:"name"`
	Members            []string `json:"members"`
	CreationDate       int      `json:"creationDate"`
	FirstAlbum         string   `json:"firstAlbum"`
	Locations_unpacked []string
}

func Homepage(w http.ResponseWriter, r *http.Request) {
	response, err := http.Get("https://groupietrackers.herokuapp.com/api/artists")
	if err != nil {
		http.Error(w, "Error fetching artist data", http.StatusInternalServerError)
		fmt.Printf("Error fetching artist data: %v\n", err)
		return
	}
	defer response.Body.Close()

	data, err := io.ReadAll(response.Body)
	if err != nil {
		http.Error(w, "Error reading artist data", http.StatusInternalServerError)
		fmt.Printf("Error reading artist data: %v\n", err)
		return
	}

	locations_response, err := http.Get("https://groupietrackers.herokuapp.com/api/locations")
	if err != nil {
		fmt.Printf("Error fetching location data: %v\n", err)
	}
	defer locations_response.Body.Close()

	locations_data, err := io.ReadAll(locations_response.Body)
	if err != nil {
		fmt.Printf("Error reading location data: %v\n", err)
	}

	var Location_format LocationData
	if err := json.NewDecoder(strings.NewReader(string(locations_data))).Decode(&Location_format); err != nil {
		fmt.Printf("Error decoding JSON: %v\n", err)

	}

	var Artists []ArtistInfo
	err = json.Unmarshal(data, &Artists)
	if err != nil {
		http.Error(w, "Error unmarshaling artist data", http.StatusInternalServerError)
		fmt.Printf("Error unmarshaling artist data: %v\n", err)
		return
	}

	//Connect locations with artists
	for i := range Artists {
		for _, item := range Location_format.Index {
			if item.ID == Artists[i].Id {
				Artists[i].Locations_unpacked = item.Locations
				break
			}
		}
	}

	type TemplateData struct {
		Search_Suggestion []ArtistInfo
		ExportData        []ArtistInfo
	}

	templateData := TemplateData{
		Search_Suggestion: Artists,
		ExportData:        Search(Artists, r),
	}

	custom_tmpl, err := template.ParseFiles("template/index.html")
	if err != nil {
		http.Error(w, "Error parsing template", http.StatusInternalServerError)
		fmt.Printf("Error parsing template: %v\n", err)
		return
	}

	err = custom_tmpl.Execute(w, templateData)
	if err != nil {
		// Check if the error is due to the client disconnecting
		if !strings.Contains(err.Error(), "broken pipe") {
			// If it's not a broken pipe error, log the error
			fmt.Println("Error executing template:", err)
			http.Error(w, "Error executing template", http.StatusInternalServerError)
		}
		return
	}
}
