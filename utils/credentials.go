package utils

import (
	"log"
	"os"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/calendar/v3"
	goauth "google.golang.org/api/oauth2/v2"
)

var GoogleOauthConfig *oauth2.Config

func SetGoogleAuthConfig() {
	b, err := os.ReadFile("credentials.json")
	if err != nil {
		log.Fatal(err)
	}
	GoogleOauthConfig, err = google.ConfigFromJSON(b, calendar.CalendarScope, goauth.UserinfoEmailScope)
	if err != nil {
		log.Fatal(err)
	}
}
