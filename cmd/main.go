package main

import (
	"database/sql"
	"fmt"
	"html/template"
	"log"
	"math/rand"
	"net"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/gempir/go-twitch-irc/v3"
	"github.com/labstack/echo/v5"
	"github.com/nicklaw5/helix/v2"
	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/core"
	"github.com/rprtr258/balaboba"
	"github.com/samber/lo"

	"github.com/rprtr258/twitch-bot/internal"
	"github.com/rprtr258/twitch-bot/internal/permissions"
	"github.com/rprtr258/twitch-bot/internal/services"
	_ "github.com/rprtr258/twitch-bot/migrations"
)

var style = `body {
  background-color: #09090b;
  color: #fff;
}
a {
  color: #ececec;
  text-decoration: underline;
  font-size: 1.8em;
}
a:visited {
  color: #333;
}
table {
  margin: 0 auto;
  text-align: center;
}`

// TODO: show partially content
var (
	idsPage = template.Must(template.New("blab-list").Parse(`
<head>
  <style>` + style + // TODO: external css file
		`</style>
</head>
<body>
  <table>
  {{range .}}
    <tr>
	  {{range .}}
        <td><a href={{.URL}}>{{.ID}}</a></td>
	  {{end}}
    </tr>
  {{end}}
  </table>
</body>
`))
	blabPage = template.Must(template.New("blab-list").Parse(`
<head>
  <style>` + style + // TODO: external css file
		`</style>
</head>
<body>
  <p style="padding: 10% 15%; font-size: 1.8em;">{{.}}</p>
</body>
`))
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
	rand.Seed(time.Now().Unix())

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

		res, err := app.
			DB().
			Select("channel").
			From("joined_channels").
			Build().
			Rows()
		if err != nil {
			return err
		}
		var channel string
		for res.Next() {
			res.Scan(&channel)
			client.Join(channel)
		}

		// TODO: provide proxy in env
		proxyURL, err := url.Parse(os.Getenv("BALABOBA_PROXY"))
		if err != nil {
			panic(err)
		}

		d := net.Dialer{
			Timeout: time.Minute,
		}
		balabobaClient := balaboba.New(balaboba.ClientConfig{
			Lang: balaboba.Rus,
			HTTP: &http.Client{
				Timeout: d.Timeout,
				Transport: &http.Transport{
					DialTLSContext:      d.DialContext,
					TLSHandshakeTimeout: d.Timeout,
					Proxy:               http.ProxyURL(proxyURL),
				},
			},
		})

		permissions, err := permissions.LoadFromJSONFile("permissions.json")
		if err != nil {
			return err
		}

		services := services.Services{
			ChatClient:      client,
			TwitchApiClient: helixClient,
			Backend:         app,
			Balaboba:        balabobaClient,
			Permissions:     permissions,
		}

		client.OnPrivateMessage(internal.OnPrivateMessage(&services))
		client.OnWhisperMessage(internal.OnWhisperMessage(&services))

		go func() {
			if err := client.Connect(); err != nil {
				log.Fatal(err.Error())
			}
		}()

		return nil
	})

	app.OnBeforeServe().Add(func(e *core.ServeEvent) error {
		if _, err := e.Router.AddRoute(echo.Route{
			Method: http.MethodGet,
			Path:   "/blab",
			Handler: func(c echo.Context) error {
				db := app.DB()

				type row struct {
					ID string
				}
				var ids []row
				err := db.
					Select("id").
					From("blab").
					All(&ids)
				if err != nil && err != sql.ErrNoRows {
					// TODO: save log
					log.Println("failed extracting ids from db", err.Error())
					return echo.ErrInternalServerError
				}

				type pageEntry struct {
					ID  string
					URL string
				}
				entries := lo.Map(ids, func(row row, _ int) pageEntry {
					return pageEntry{
						ID:  row.ID,
						URL: fmt.Sprintf("/%s/%s", "blab", row.ID),
					}
				})
				entriesChunks := lo.Chunk(entries, 4)

				var page strings.Builder
				idsPage.Execute(&page, entriesChunks)
				return c.HTML(http.StatusOK, page.String())
			},
		}); err != nil {
			return err
		}

		if _, err := e.Router.AddRoute(echo.Route{
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
				var page strings.Builder
				blabPage.Execute(&page, text)
				return c.HTML(http.StatusOK, page.String())
			},
		}); err != nil {
			return err
		}

		return nil
	})

	return app.Start()
}

func main() {
	if err := run(); err != nil {
		log.Fatal(err.Error())
	}
}
