package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gempir/go-twitch-irc/v3"
	"github.com/karalef/balaboba"
	"github.com/labstack/echo/v5"
	"github.com/nicklaw5/helix"
	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/core"

	// TODO: rename module, packages
	"abobus/internal"
	"abobus/internal/permissions"
	"abobus/internal/services"
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

		services := services.Services{
			ChatClient:      client,
			TwitchApiClient: helixClient,
			Backend:         app,
			Balaboba:        balabobaClient,
			// TODO: move out to file
			Permissions: permissions.Permissions{
				"global_admin": []permissions.Claims{{
					"username": "rprtr258",
				}},
				"admin": []permissions.Claims{{
					"username": "rprtr258",
					"channel":  "rprtr258",
				}},
				"say_response": []permissions.Claims{{
					// "username": "rprtr258",
					"channel": "rprtr258",
				}},
				"execute_commands": []permissions.Claims{{
					"channel": "rprtr258",
				}},
			},
		}

		client.OnPrivateMessage(internal.OnPrivateMessage(&services))

		go func() {
			if err := client.Connect(); err != nil {
				log.Fatal(err.Error())
			}
		}()

		return nil
	})

	app.OnBeforeServe().Add(func(e *core.ServeEvent) error {
		e.Router.AddRoute(echo.Route{
			Method: http.MethodGet,
			Path:   "/blab/:id",
			Handler: func(c echo.Context) error {
				blabID := c.PathParam("id")
				db := app.DB()

				var text string
				err := db.Select("text").From("blab").Where(dbx.NewExp(
					"id={:id}", dbx.Params{"id": blabID},
				)).Row(&text)
				if err != nil && err != sql.ErrNoRows {
					// TODO: save log
					log.Println(err.Error())
					return echo.ErrInternalServerError
				} else if err == sql.ErrNoRows {
					return echo.ErrNotFound
				}

				// TODO: use html template with text safeguarding
				return c.HTML(http.StatusOK, fmt.Sprintf(
					`<p style="padding: 10%% 15%%; font-size: 1.8em;">%s</p>`,
					text,
				))
			},
		})

		return nil
	})

	return app.Start()
}

func main() {
	if err := run(); err != nil {
		log.Fatal(err.Error())
	}
}
