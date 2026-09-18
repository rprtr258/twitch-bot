package migrations

import (
	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase/core"
	m "github.com/pocketbase/pocketbase/migrations"
)

func init() {
	m.Register(func(app core.App) error {
		db := app.DB()
		if _, err := db.
			Update(
				"chat_commands",
				dbx.Params{
					"created": dbx.NewExp("substring(at, 0, 24)"),
				},
				dbx.NewExp("created=''"),
			).
			Execute(); err != nil {
			return err
		}
		if _, err := db.
			Update(
				"messages",
				dbx.Params{
					"created": dbx.NewExp("substring(at, 0, 24)"),
				},
				dbx.NewExp("created=''"),
			).
			Execute(); err != nil {
			return err
		}
		return nil
	}, func(app core.App) error {
		return nil
	})
}
