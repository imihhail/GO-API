package main

import (
	"encoding/json"
	"fmt"
	"groupie/functions"
	"io"
	"net/http"
	"text/template"
)

type ArtistInfo struct {
	Image          string              `json:"image"`
	Name           string              `json:"name"`
	Members        []string            `json:"members"`
	CreationDate   int                 `json:"creationDate"`
	FirstAlbum     string              `json:"firstAlbum"`
	Dates          []string            `json:"dates"`
	Locations      []string            `json:"locations"`
	Relationlink   string              `json:"relations"`
	DatesLocations map[string][]string `json:"datesLocations"`
	//Formated data
	Str_Relations string
	Str_Locations string
	Str_Dates     string
}

func Single_Page(w http.ResponseWriter, r *http.Request) {
	artist_ID := r.URL.Query().Get("name")
	url_path := r.URL.Path

	var Export_data ArtistInfo

	// If user goes on locations or dates page, then start extracting data.
	datesLocations_URL, err := http.Get("https://groupietrackers.herokuapp.com/api" + url_path + "/" + artist_ID)
	if err != nil {
		panic(err)
	}
	defer datesLocations_URL.Body.Close()

	datesLocation_reader, err := io.ReadAll(datesLocations_URL.Body)
	if err != nil {
		panic(err)
	}

	var dates_locations ArtistInfo

	err = json.Unmarshal(datesLocation_reader, &dates_locations)
	if err != nil && url_path != "/locations" && url_path != "/artists" {
		panic(err)
	}

	// If relationlink exists in a page, then execute then start extracting single artistpage data and data inside relationlink.
	if dates_locations.Relationlink != "" {
		relations_URL, _ := http.Get(dates_locations.Relationlink)

		relations_reader, _ := io.ReadAll(relations_URL.Body)

		defer relations_URL.Body.Close()

		var full_data ArtistInfo
		err = json.Unmarshal(relations_reader, &full_data)
		if err != nil {
			fmt.Println("Error unmarshaling JSON:", err)
		}
		
		Export_data = ArtistInfo{
			Image:         dates_locations.Image,
			Name:          dates_locations.Name,
			Members:       dates_locations.Members,
			CreationDate:  dates_locations.CreationDate,
			FirstAlbum:    dates_locations.FirstAlbum,
			Str_Relations: functions.Str_manip(full_data.DatesLocations), // Str_manip makes JSON data into a string and makes data more user-friendly
		}
		

	} else {
		// If there's no relationlink on a page, then you're in a location or dates page and start extracting data.
		Export_data = ArtistInfo{

			Str_Locations: functions.Str_manip(dates_locations.Locations),
			Str_Dates:     functions.Str_manip(dates_locations.Dates),
		}
	}

	custom_tmpl, err := template.ParseFiles("template/infopage.html")
	if err != nil {
		http.Error(w, "500, Error parsing template.", http.StatusInternalServerError)
		return
	}

	err = custom_tmpl.Execute(w, Export_data)
	if err != nil {
		http.Error(w, "Error executing template", http.StatusInternalServerError)
		return
	}

}
