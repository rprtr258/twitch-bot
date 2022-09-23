package main

import (
	"log"
	"os"

	"github.com/gempir/go-twitch-irc/v3"
	"github.com/nicklaw5/helix"
	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/core"

	abobus "abobus/internal"
)

func handleTwitchAuth(helixClient *helix.Client) {
	resp, err := helixClient.RequestAppAccessToken([]string{"user:read:email"})
	if err != nil {
		log.Println(err.Error())
		return
	}

	// TODO: somehow refresh after ExpiresIn seconds
	helixClient.SetAppAccessToken(resp.Data.AccessToken)
}

func run() error {
	app := pocketbase.New()

	app.OnBeforeServe().Add(func(data *core.ServeEvent) error {
		helixClient, err := helix.NewClient(&helix.Options{
			ClientID:     os.Getenv("TWITCH_CLIENT_ID"),
			ClientSecret: os.Getenv("TWITCH_CLIENT_SECRET"),
			RedirectURI:  os.Getenv("TWITCH_REDIRECT_URI"),
		})
		if err != nil {
			return err
		}

		handleTwitchAuth(helixClient)

		client := twitch.NewClient("trooba_bot", os.Getenv("TWITCH_OAUTH_TOKEN"))
		client.Join("vs_code")

		srvcs := abobus.Services{
			ChatClient:      client,
			TwitchApiClient: helixClient,
			Backend:         app,
		}

		client.OnPrivateMessage(srvcs.OnPrivateMessage)

		go func() {
			if err := client.Connect(); err != nil {
				log.Fatal(err.Error())
			}
		}()

		return nil
	})

	return app.Start()
}

func main() {
	if err := run(); err != nil {
		log.Fatal(err.Error())
	}
}
