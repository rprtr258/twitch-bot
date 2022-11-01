package main

import (
	"html/template"
	"log"
	"math/rand"
	"time"

	"github.com/nicklaw5/helix/v2"
	"github.com/pocketbase/pocketbase"

	// TODO: rename packages

	_ "github.com/rprtr258/twitch-bot/migrations"
)

// TODO: show partially content
var idsPage = template.Must(template.New("blab-list").Parse(`
{{range .}}
	<a style="padding: 10% 15%; font-size: 1.8em;" href={{.URL}}>{{.ID}}</a>
{{end}}
`))

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

	return app.Start()
}

func main() {
	if err := run(); err != nil {
		log.Fatal(err.Error())
	}
}
