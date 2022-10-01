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

	commandsCmd = "?commands"
)

type Services struct {
	ChatClient      *twitch.Client
	TwitchApiClient *helix.Client
	Backend         *pocketbase.PocketBase
	Balaboba        *balaboba.Client
}

type command struct {
	Relation    string
	Command     string
	Run         func(*Services, twitch.PrivateMessage) (string, error)
	Whisper     bool
	Description string
}

var cmds []command

func init() {
	cmds = []command{{
		Relation:    executeRelation,
		Command:     intelCmd,
		Run:         (*Services).getIntelCmd,
		Whisper:     true,
		Description: "Gather intel on user",
	}, {
		Relation:    executeRelation,
		Command:     joinCmd,
		Run:         (*Services).join,
		Description: "Join channel",
	}, {
		Relation:    executeRelation,
		Command:     leaveCmd,
		Run:         (*Services).leave,
		Description: "Leave channel",
	}, {
		Relation:    everyoneRelation,
		Command:     fedCmd,
		Run:         (*Services).fed,
		Description: "Show how many times word has been used",
	}, {
		Relation:    everyoneRelation,
		Command:     blabCmd,
		Run:         (*Services).blab,
		Description: "Balaboba text generation neural network",
	}, {
		Relation:    everyoneRelation,
		Command:     pythCmd,
		Run:         (*Services).pyth,
		Description: "Eval in pyth",
	}, {
		Relation:    everyoneRelation,
		Command:     commandsCmd,
		Run:         (*Services).commands,
		Description: "List all commands",
	}}
}

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
		if (cmd.Command == fedCmd || cmd.Command == blabCmd || cmd.Command == pythCmd || cmd.Command == commandsCmd) && message.Channel != "rprtr258" {
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
			if message.User.Name == "rprtr258" {
				// TODO: fix: not showing short responses to me
				runes := []rune(response)
				for i := 0; i < len(runes); i += _maxMessageLength {
					resp := string(runes[i : i+_maxMessageLength])
					if whisper {
						s.ChatClient.Whisper(message.User.Name, resp)
					} else {
						s.ChatClient.Say(message.Channel, resp)
					}
				}
			} else {
				runes := []rune(response)
				response = string(runes[:_maxMessageLength])

				if whisper {
					s.ChatClient.Whisper(message.User.Name, response)
				} else {
					s.ChatClient.Reply(message.Channel, message.ID, response)
				}
			}
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

func (s *Services) commands(message twitch.PrivateMessage) (string, error) {
	parts := make([]string, 0, len(cmds))

	for _, cmd := range cmds {
		parts = append(parts, fmt.Sprintf("%s - %s", cmd.Command, cmd.Description))
	}

	return strings.Join(parts, ", "), nil
}
