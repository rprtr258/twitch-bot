package internal

import (
	"fmt"
	"log"
	"strings"
	"time"

	twitch "github.com/gempir/go-twitch-irc/v3"
	"github.com/karalef/balaboba"
	"github.com/nicklaw5/helix"
	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/forms"
	"github.com/pocketbase/pocketbase/models"
	"github.com/samber/lo"
)

const (
	commandsCmd = "?commands"
)

// TODO: move to lib
type Claims = map[string]string
type Permissions map[string][]Claims

func (perms Permissions) getPermissions(providedClaims Claims) []string {
	res := []string{}
	for permission, requiredClaimsSlice := range perms {
		if lo.SomeBy(requiredClaimsSlice, func(requiredClaims Claims) bool {
			for k, v := range requiredClaims {
				if providedClaims[k] != v {
					return false
				}
			}
			return true
		}) {
			res = append(res, permission)
		}
	}
	return res
}

type Services struct {
	ChatClient      *twitch.Client
	TwitchApiClient *helix.Client
	Backend         *pocketbase.PocketBase
	Balaboba        *balaboba.Client
	Permissions     Permissions
}

type command struct {
	PermissionsRequired []string
	Command             string
	Run                 func(*Services, []string, twitch.PrivateMessage) (string, error)
	Description         string
}

var cmds []command

func init() {
	// TODO: permissions constants
	cmds = []command{{
		PermissionsRequired: []string{},
		Command:             intelCmd,
		Run:                 (*Services).getIntelCmd,
		Description:         "Gather intel on user",
	}, {
		PermissionsRequired: []string{"global_admin"},
		Command:             joinCmd,
		Run:                 (*Services).join,
		Description:         "Join channel",
	}, {
		PermissionsRequired: []string{"global_admin"},
		Command:             leaveCmd,
		Run:                 (*Services).leave,
		Description:         "Leave channel",
	}, {
		PermissionsRequired: []string{"execute_commands"},
		Command:             fedCmd,
		Run:                 (*Services).fed,
		Description:         "Show how many times word has been used",
	}, {
		PermissionsRequired: []string{"execute_commands"},
		Command:             blabCmd,
		Run:                 (*Services).blab,
		Description:         "Balaboba text generation neural network",
	}, {
		PermissionsRequired: []string{"execute_commands"},
		Command:             pythCmd,
		Run:                 (*Services).pyth,
		Description:         "Eval in pyth",
	}, {
		PermissionsRequired: []string{"execute_commands"},
		Command:             commandsCmd,
		Run:                 (*Services).commands,
		// TODO: list only available commands
		Description: "List all commands",
	}, {
		PermissionsRequired: []string{},
		Command:             permsCmd,
		Run:                 (*Services).getPermissions,
		Description:         "Get permissions",
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

// TODO: handle whispers
func (s *Services) OnPrivateMessage(message twitch.PrivateMessage) {
	s.logMessage(message)
	userPermissions := s.Permissions.getPermissions(Claims{
		"username": message.User.Name, // TODO: replace with user id
		"channel":  message.Channel,
		// TODO: vips, moders
	})

	firstWord := strings.Split(message.Message, " ")[0]
	userName := message.User.Name

	for _, cmd := range cmds {
		if cmd.Command != firstWord {
			continue
		}

		// TODO: add ban permissions
		if !lo.Every(userPermissions, cmd.PermissionsRequired) {
			break
		}

		whisper := !lo.Contains(userPermissions, "say_response")

		response, err := cmd.Run(s, userPermissions, message)
		if err != nil {
			response = fmt.Sprintf("Internal error: %s", err.Error())
		} else if response == "" {
			response = "Empty response"
		}
		response = strings.ReplaceAll(response, "\n", " ")
		// TODO: move responding to commands
		// TODO: rewrite
		if len(response) > _maxMessageLength {
			if message.User.Name == "rprtr258" {
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
		} else {
			if whisper {
				s.ChatClient.Whisper(message.User.Name, response)
			} else {
				s.ChatClient.Reply(message.Channel, message.ID, response)
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

func (s *Services) commands(perms []string, message twitch.PrivateMessage) (string, error) {
	parts := make([]string, 0, len(cmds))

	for _, cmd := range cmds {
		parts = append(parts, fmt.Sprintf("%s - %s", cmd.Command, cmd.Description))
	}

	return strings.Join(parts, ", "), nil
}
