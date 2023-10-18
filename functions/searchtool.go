package functions

import (
	"net/http"
	"regexp"
	"strconv"
	"strings"
)

func Search(strc []ArtistInfo, r *http.Request) []ArtistInfo {
	input := r.URL.Query().Get("name")

	re := regexp.MustCompile(`\s\(.*`)
	search_input := re.ReplaceAllString(input, "")

	var newArtists []ArtistInfo
	membercount := 0

	for _, artist := range strc {
		if strings.Contains(strings.Title(artist.Name), search_input) || strings.Contains(strings.ToLower(artist.Name), search_input) {
			newArtists = append(newArtists, artist)
			continue
		}

		for _, member := range artist.Members {
			if strings.Contains(strings.Title(member), search_input) || strings.Contains(strings.ToLower(member), search_input) {
				membercount++
				newArtists = append(newArtists, artist)
				break
			}
		}

		for _, member := range artist.Locations_unpacked {
			if strings.Contains(strings.Title(member), search_input) || strings.Contains(strings.ToLower(member), search_input) {
				if len(search_input) < 2 && membercount > 0 {
					break
				}
				newArtists = append(newArtists, artist)
				break
			}
		}

		if strconv.Itoa(artist.CreationDate) == search_input {
			newArtists = append(newArtists, artist)
			continue
		}

		if artist.FirstAlbum == search_input {
			newArtists = append(newArtists, artist)
			continue
		}

	}
	return newArtists
}
