package functions

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strings"
)

func Marshal(data interface{}) string {
	unFormated_JsonText, err := json.Marshal(data)
	if err != nil {
		fmt.Println("Unable to Marshal data")
	}
	marshaled_str := string(unFormated_JsonText)

	return marshaled_str
}

func Str_manip(marshaled_str interface{}) string {
	unFormated_text := Marshal(marshaled_str)

	Formated_text := strings.ReplaceAll(unFormated_text, `":`, "<br>")
	Formated_text = strings.ReplaceAll(Formated_text, `","`, "<br>")
	Formated_text = strings.ReplaceAll(Formated_text, `,"`, "<br><br>")
	Formated_text = strings.ReplaceAll(Formated_text, "_", " ")

	re := regexp.MustCompile(`[_*"\{\}\[\]]`)
	Formated_text = re.ReplaceAllString(Formated_text, "")

	Formated_text = strings.ToUpper(Formated_text)
	if Formated_text == "NULL" {
		Formated_text = ""
	}
	return Formated_text
}
