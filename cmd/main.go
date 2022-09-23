package main

import (
	"log"
	"net/http"
	"os"

	twitch "github.com/gempir/go-twitch-irc/v3"
	"github.com/nicklaw5/helix"

	abobus "abobus/internal"
)

func handleTwitchAuth(helixClient *helix.Client) {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		code, ok := r.URL.Query()["code"]
		if !ok {
			return
		}

		resp, err := helixClient.RequestUserAccessToken(code[0])
		if err != nil {
			log.Println(err.Error())
		}

		// Set the access token on the client
		// TODO: remember token, https://github.com/nicklaw5/helix/blob/main/docs/authentication_docs.md#refresh-user-access-token
		helixClient.SetUserAccessToken(resp.Data.AccessToken)
	})
	if err := http.ListenAndServe(":5000", nil); err != nil {
		log.Println(err.Error())
	}
}

func run() error {
	helixClient, err := helix.NewClient(&helix.Options{
		ClientID:     os.Getenv("TWITCH_CLIENT_ID"),
		ClientSecret: os.Getenv("TWITCH_CLIENT_SECRET"),
		RedirectURI:  os.Getenv("TWITCH_REDIRECT_URI"),
	})
	if err != nil {
		return err
	}

	go handleTwitchAuth(helixClient)

	url := helixClient.GetAuthorizationURL(&helix.AuthorizationURLParams{
		ResponseType: "code",
		Scopes:       []string{"user:read:email"},
		ForceVerify:  false,
	})
	log.Println(url)

	client := twitch.NewClient("trooba_bot", os.Getenv("TWITCH_OAUTH_TOKEN"))
	client.Join("vs_code")

	srvcs := abobus.Services{
		ChatClient:      client,
		TwitchApiClient: helixClient,
	}

	client.OnPrivateMessage(srvcs.OnPrivateMessage)

	if err := client.Connect(); err != nil {
		return err
	}

	return nil
}

func main() {
	if err := run(); err != nil {
		log.Fatal(err.Error())
	}
}
