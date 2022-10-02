package services

import (
	"log"

	twitch "github.com/gempir/go-twitch-irc/v3"
	"github.com/karalef/balaboba"
	"github.com/nicklaw5/helix"
	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/forms"
	"github.com/pocketbase/pocketbase/models"

	"abobus/internal/permissions"
)

type Services struct {
	ChatClient      *twitch.Client
	TwitchApiClient *helix.Client
	Backend         *pocketbase.PocketBase
	Balaboba        *balaboba.Client
	Permissions     permissions.Permissions
}

func (s *Services) Insert(collectionName string, data map[string]any) (string, error) {
	collection, err := s.Backend.Dao().FindCollectionByNameOrId(collectionName)
	if err != nil {
		return "", err
	}

	record := models.NewRecord(collection)
	for k, v := range data {
		record.SetDataValue(k, v)
	}

	form := forms.NewRecordUpsert(s.Backend.App, record)

	if err := form.Validate(); err != nil {
		return "", err
	}

	if err := form.Submit(); err != nil {
		return "", err
	}

	return record.Id, nil
}

func (s *Services) LogMessage(message twitch.PrivateMessage) {
	_, err := s.Insert("messages", map[string]any{
		"user_id":           message.User.ID,
		"message":           message.Message,
		"at":                message.Time,
		"channel":           message.Channel,
		"user_name":         message.User.Name,
		"user_display_name": message.User.DisplayName,
	})
	// TODO: users table
	if err != nil {
		// TODO: save log
		log.Println(err.Error())
	}
}
