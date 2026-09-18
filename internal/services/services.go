package services

import (
	"log"

	twitch "github.com/gempir/go-twitch-irc/v3"
	"github.com/nicklaw5/helix/v2"
	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/core"
	"github.com/rprtr258/balaboba"

	"github.com/rprtr258/twitch-bot/internal/message"
	"github.com/rprtr258/twitch-bot/internal/permissions"
)

type Services struct {
	ChatClient      *twitch.Client
	TwitchApiClient *helix.Client
	Backend         *pocketbase.PocketBase
	Balaboba        *balaboba.Client
	Permissions     permissions.Permissions
}

func (s *Services) Insert(collectionName string, data map[string]any) (string, error) {
	collection, err := s.Backend.FindCollectionByNameOrId(collectionName)
	if err != nil {
		return "", err
	}

	record := core.NewRecord(collection)
	for k, v := range data {
		record.Set(k, v)
	}

	// Save validates the record against the collection schema, like the
	// removed forms.RecordUpsert Validate+Submit pair did before v0.23.
	if err := s.Backend.Save(record); err != nil {
		return "", err
	}

	return record.Id, nil
}

func (s *Services) LogMessage(msg message.TwitchMessage) {
	_, err := s.Insert("messages", map[string]any{
		"user_id":           msg.User.ID,
		"message":           msg.Text,
		"channel":           msg.Channel,
		"user_name":         msg.User.Name,
		"user_display_name": msg.User.DisplayName,
	})
	// TODO: users table
	if err != nil {
		// TODO: save log
		log.Println(err.Error())
	}
}
