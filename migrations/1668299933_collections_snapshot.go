package migrations

import (
	"github.com/pocketbase/pocketbase/core"
	m "github.com/pocketbase/pocketbase/migrations"
)

// Auto generated migration with the most recent collections configuration.
func init() {
	m.Register(func(app core.App) error {
		jsonData := `
		[
			{
				"id": "_pb_users_auth_",
				"name": "users",
				"type": "auth",
				"fields": [
					{
						"id": "pbfielduser",
						"name": "userId",
						"required": true,
						"type": "relation",
						"collectionId": "_pb_users_auth_",
						"cascadeDelete": true,
						"minSelect": null,
						"maxSelect": 1,
						"display": "id"
					},
					{
						"id": "pbfieldname",
						"name": "name",
						"required": false,
						"type": "text",
						"min": null,
						"max": null,
						"pattern": ""
					},
					{
						"id": "pbfieldavatar",
						"name": "avatar",
						"required": false,
						"type": "file",
						"maxSelect": 1,
						"maxSize": 5242880,
						"mimeTypes": [
							"image/jpg",
							"image/jpeg",
							"image/png",
							"image/svg+xml",
							"image/gif"
						],
						"thumbs": null,
						"protected": false
					},
					{
						"id": "created000001",
						"name": "created",
						"type": "autodate",
						"onCreate": true,
						"onUpdate": false
					},
					{
						"id": "updated000001",
						"name": "updated",
						"type": "autodate",
						"onCreate": true,
						"onUpdate": true
					}
				],
				"indexes": [
					"CREATE UNIQUE INDEX idx_users_userId ON users (userId)"
				],
				"listRule": "userId = @request.auth.id",
				"viewRule": "userId = @request.auth.id",
				"createRule": "userId = @request.auth.id",
				"updateRule": "userId = @request.auth.id",
				"deleteRule": null
			},
			{
				"id": "xyqfe2absuqms10",
				"name": "chat_commands",
				"type": "base",
				"fields": [
					{
						"id": "rz0gibsj",
						"name": "command",
						"required": true,
						"type": "text",
						"min": null,
						"max": null,
						"pattern": ""
					},
					{
						"id": "jmsyjohs",
						"name": "response",
						"required": true,
						"type": "text",
						"min": null,
						"max": null,
						"pattern": ""
					},
					{
						"id": "mm36zyq2",
						"name": "user",
						"required": true,
						"type": "text",
						"min": null,
						"max": null,
						"pattern": ""
					},
					{
						"id": "jre7r0ow",
						"name": "args",
						"required": true,
						"type": "text",
						"min": null,
						"max": null,
						"pattern": ""
					},
					{
						"id": "7uogu1bv",
						"name": "channel",
						"required": true,
						"type": "text",
						"min": null,
						"max": null,
						"pattern": ""
					},
					{
						"id": "created000001",
						"name": "created",
						"type": "autodate",
						"onCreate": true,
						"onUpdate": false
					},
					{
						"id": "updated000001",
						"name": "updated",
						"type": "autodate",
						"onCreate": true,
						"onUpdate": true
					}
				],
				"indexes": null,
				"listRule": null,
				"viewRule": null,
				"createRule": null,
				"updateRule": null,
				"deleteRule": null
			},
			{
				"id": "ji0v4q3n0le3s8i",
				"name": "joined_channels",
				"type": "base",
				"fields": [
					{
						"id": "8a85du96",
						"name": "channel",
						"required": true,
						"type": "text",
						"min": null,
						"max": null,
						"pattern": "^[a-zA-Z0-9_]+$"
					},
					{
						"id": "created000001",
						"name": "created",
						"type": "autodate",
						"onCreate": true,
						"onUpdate": false
					},
					{
						"id": "updated000001",
						"name": "updated",
						"type": "autodate",
						"onCreate": true,
						"onUpdate": true
					}
				],
				"indexes": [
					"CREATE UNIQUE INDEX idx_joined_channels_channel ON joined_channels (channel)"
				],
				"listRule": null,
				"viewRule": null,
				"createRule": null,
				"updateRule": null,
				"deleteRule": null
			},
			{
				"id": "vvifj1r17756fvq",
				"name": "messages",
				"type": "base",
				"fields": [
					{
						"id": "u7wth8tl",
						"name": "message",
						"required": true,
						"type": "text",
						"min": null,
						"max": null,
						"pattern": ""
					},
					{
						"id": "fm0quhrb",
						"name": "channel",
						"required": true,
						"type": "text",
						"min": null,
						"max": null,
						"pattern": ""
					},
					{
						"id": "ae9chw4d",
						"name": "user_id",
						"required": true,
						"type": "number",
						"min": null,
						"max": null,
						"onlyInt": false
					},
					{
						"id": "jizud6gh",
						"name": "user_name",
						"required": true,
						"type": "text",
						"min": null,
						"max": null,
						"pattern": ""
					},
					{
						"id": "ul51abnp",
						"name": "user_display_name",
						"required": true,
						"type": "text",
						"min": null,
						"max": null,
						"pattern": ""
					},
					{
						"id": "created000001",
						"name": "created",
						"type": "autodate",
						"onCreate": true,
						"onUpdate": false
					},
					{
						"id": "updated000001",
						"name": "updated",
						"type": "autodate",
						"onCreate": true,
						"onUpdate": true
					}
				],
				"indexes": null,
				"listRule": null,
				"viewRule": null,
				"createRule": null,
				"updateRule": null,
				"deleteRule": null
			},
			{
				"id": "rgkya276bix2c7v",
				"name": "blab",
				"type": "base",
				"fields": [
					{
						"id": "npnh0uak",
						"name": "text",
						"required": true,
						"type": "text",
						"min": null,
						"max": null,
						"pattern": ""
					},
					{
						"id": "w9ro7ake",
						"name": "author_id",
						"required": true,
						"type": "number",
						"min": null,
						"max": null,
						"onlyInt": false
					},
					{
						"id": "mimpxbvg",
						"name": "continuations",
						"required": true,
						"type": "number",
						"min": null,
						"max": null,
						"onlyInt": false
					},
					{
						"id": "5kywk8ep",
						"name": "style",
						"required": false,
						"type": "number",
						"min": 0,
						"max": null,
						"onlyInt": false
					},
					{
						"id": "xuetf5tz",
						"name": "channel",
						"required": true,
						"type": "text",
						"min": null,
						"max": null,
						"pattern": ""
					},
					{
						"id": "created000001",
						"name": "created",
						"type": "autodate",
						"onCreate": true,
						"onUpdate": false
					},
					{
						"id": "updated000001",
						"name": "updated",
						"type": "autodate",
						"onCreate": true,
						"onUpdate": true
					}
				],
				"indexes": null,
				"listRule": null,
				"viewRule": null,
				"createRule": null,
				"updateRule": null,
				"deleteRule": null
			}
		]
		`

		return app.ImportCollectionsByMarshaledJSON([]byte(jsonData), true)
	}, func(app core.App) error {
		// no revert since the configuration on the environment, on which
		// the migration was executed, could have changed via the UI/API
		return nil
	})
}
