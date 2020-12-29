package googleAuth

import (
	"io/ioutil"
	"log"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/calendar/v3"
)

func readFileCredentials() ([]byte, error) {
	b, err := ioutil.ReadFile("credentials.json")
	return b, err
}

func configGoogle() (*oauth2.Config, error) {
	b, err := readFileCredentials()
	if err != nil {
		log.Fatalf("Unable to read client secret file: %v", err)
	}

	config, err := google.ConfigFromJSON(b, calendar.CalendarScope)
	return config, err
}
