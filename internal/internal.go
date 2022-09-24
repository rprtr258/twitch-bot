package internal

import (
	"database/sql"
	"fmt"
	"log"
	"strings"
	"time"

	twitch "github.com/gempir/go-twitch-irc/v3"
	"github.com/karalef/balaboba"
	"github.com/nicklaw5/helix"
	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/forms"
	"github.com/pocketbase/pocketbase/models"
)

const (
	everyoneRelation  = ""
	executeRelation   = "execute"
	forbiddenRelation = "forbidden"
)

type Services struct {
	ChatClient      *twitch.Client
	TwitchApiClient *helix.Client
	Backend         *pocketbase.PocketBase
	Balaboba        *balaboba.Client
}

type command struct {
	Relation string
	Command  string
	Run      func(*Services, twitch.PrivateMessage) (string, error)
	Whisper  bool
}

var cmds []command = []command{{
	Relation: executeRelation,
	Command:  intelCmd,
	Run:      (*Services).getIntelCmd,
	Whisper:  true,
}, {
	Relation: executeRelation,
	Command:  joinCmd,
	Run:      (*Services).join,
}, {
	Relation: executeRelation,
	Command:  leaveCmd,
	Run:      (*Services).leave,
}, {
	Relation: everyoneRelation,
	Command:  fedCmd,
	Run:      (*Services).fed,
}, {
	Relation: everyoneRelation,
	Command:  blabCmd,
	Run:      (*Services).blab,
}, {
	Relation: everyoneRelation,
	Command:  pythCmd,
	Run:      (*Services).pyth,
}}

func (s *Services) logMessage(message twitch.PrivateMessage) {
	_, err := s._insert("messages", map[string]any{
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

func (s *Services) OnPrivateMessage(message twitch.PrivateMessage) {
	db := s.Backend.DB()

	s.logMessage(message)

	firstWord := strings.Split(message.Message, " ")[0]
	userName := message.User.Name

	var relation string
	if err := db.Select("relation").From("permissions").Where(dbx.And(
		dbx.NewExp("user_id={:user_id}", dbx.Params{"user_id": message.User.ID}),
		dbx.NewExp("object={:cmd}", dbx.Params{"cmd": fmt.Sprintf("cmd:%s", firstWord)}),
	)).Build().Row(&relation); err != nil {
		if err != sql.ErrNoRows {
			log.Println(err.Error())
			return
		}
	}

	for _, cmd := range cmds {
		if cmd.Command != firstWord {
			continue
		}

		if relation == forbiddenRelation {
			break
		}

		if cmd.Relation != everyoneRelation && cmd.Relation != relation {
			break
		}

		// TODO: peredelat'
		if (cmd.Command == fedCmd || cmd.Command == blabCmd) && message.Channel != "rprtr258" {
			break
		}

		// TODO: peredelat'
		whisper := cmd.Whisper
		if cmd.Command == intelCmd && message.Channel == "rprtr258" {
			whisper = false
		}

		response, err := cmd.Run(s, message)
		if err != nil {
			response = fmt.Sprintf("Internal error: %s", err.Error())
		} else if response == "" {
			response = "Empty response"
		}
		response = strings.ReplaceAll(response, "\n", " ")
		if len(response) > _maxMessageLength {
			runes := []rune(response)
			response = string(runes[:_maxMessageLength])
		}

		if whisper {
			s.ChatClient.Whisper(message.User.Name, response)
		} else {
			s.ChatClient.Reply(message.Channel, message.ID, response)
		}

		if _, err := s._insert("chat_commands", map[string]any{
			"command":  cmd.Command,
			"args":     message.Message,
			"at":       time.Now(),
			"response": response,
			"user":     userName,
			"channel":  message.Channel,
		}); err != nil {
			// TODO: save log
			log.Println(err.Error())
		}

		break
	}
}

func (s *Services) _insert(collectionName string, data map[string]any) (string, error) {
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
