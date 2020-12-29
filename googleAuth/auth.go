package googleAuth

import (
	"context"
	"log"
	"net/http"

	"golang.org/x/oauth2"
)

type GoogleAuth struct {
	Config *oauth2.Config
}

func (auth GoogleAuth) GetAuthCodeURL() string {
	authURL := auth.Config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
	return authURL
}

func (auth GoogleAuth) GetTokenFormCode(code string) *oauth2.Token {
	tok, err := auth.Config.Exchange(context.TODO(), code)
	if err != nil {
		log.Fatalf("Unable to retrieve token from web: %v", err)
	}
	return tok
}

func (auth GoogleAuth) GetHtpClient(token *oauth2.Token) *http.Client {
	return auth.Config.Client(context.Background(), token)
}
