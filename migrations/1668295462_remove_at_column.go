package migrations

import (
	"github.com/pocketbase/dbx"
	m "github.com/pocketbase/pocketbase/migrations"
)

func init() {
	m.Register(func(db dbx.Builder) error {
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
	}, func(db dbx.Builder) error {
		return nil
	})
}
