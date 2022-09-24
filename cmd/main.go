package main

import (
	"log"
	"os"

	"github.com/gempir/go-twitch-irc/v3"
	"github.com/karalef/balaboba"
	"github.com/nicklaw5/helix"
	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/core"

	abobus "abobus/internal"
)

func twitchAuth(helixClient *helix.Client) {
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

		twitchAuth(helixClient)

		client := twitch.NewClient(os.Getenv("TWITCH_USERNAME"), os.Getenv("TWITCH_OAUTH_TOKEN"))

		res, err := app.DB().Select("channel").From("joined_channels").Build().Rows()
		if err != nil {
			return err
		}
		var channel string
		for res.Next() {
			res.Scan(&channel)
			client.Join(channel)
		}

		balabobaClient := balaboba.New(balaboba.Rus)

		srvcs := abobus.Services{
			ChatClient:      client,
			TwitchApiClient: helixClient,
			Backend:         app,
			Balaboba:        balabobaClient,
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
