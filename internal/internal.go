package internal

import (
	"abobus/internal/cmds"
	"abobus/internal/permissions"
	"fmt"
	"log"
	"strings"
	"time"

	twitch "github.com/gempir/go-twitch-irc/v3"
	"github.com/samber/lo"
)

const _maxMessageLength = 500

// TODO: list only available commands
type CommandsCmd struct{}

func (CommandsCmd) Command() string {
	return "?commands"
}

func (CommandsCmd) Description() string {
	return "List all commands"
}

func (CommandsCmd) Run(s *cmds.Services, perms []string, message twitch.PrivateMessage) (string, error) {
	parts := make([]string, 0, len(allCommands))

	for _, cmd := range allCommands {
		parts = append(parts, fmt.Sprintf("%s - %s", cmd.Cmd.Command(), cmd.Cmd.Description()))
	}

	return strings.Join(parts, ", "), nil
}

type Command struct {
	PermissionsRequired []string
	Cmd                 cmds.Command
}

// TODO: permissions constants/store in db to dinamically update
var allCommands = []Command{{
	PermissionsRequired: []string{},
	Cmd:                 cmds.IntelCmd{},
}, {
	PermissionsRequired: []string{"global_admin"},
	Cmd:                 cmds.JoinCmd{},
}, {
	PermissionsRequired: []string{"global_admin"},
	Cmd:                 cmds.LeaveCmd{},
}, {
	PermissionsRequired: []string{"execute_commands"},
	Cmd:                 cmds.FedCmd{},
}, {
	PermissionsRequired: []string{"execute_commands"},
	Cmd:                 cmds.BlabCmd{},
}, {
	PermissionsRequired: []string{"execute_commands"},
	Cmd:                 &cmds.PythCmd{},
}, {
	PermissionsRequired: []string{"execute_commands"},
	Cmd:                 CommandsCmd{},
}, {
	PermissionsRequired: []string{},
	Cmd:                 cmds.PermsCmd{},
}}

// TODO: handle whispers
func OnPrivateMessage(s *cmds.Services) func(twitch.PrivateMessage) {
	return func(message twitch.PrivateMessage) {
		s.LogMessage(message)
		userPermissions := s.Permissions.GetPermissions(permissions.Claims{
			"username": message.User.Name, // TODO: replace with user id
			"channel":  message.Channel,
			// TODO: vips, moders
		})

		firstWord := strings.Split(message.Message, " ")[0]
		userName := message.User.Name

		for _, cmd := range allCommands {
			if cmd.Cmd.Command() != firstWord {
				continue
			}

			// TODO: add ban permissions
			if !lo.Every(userPermissions, cmd.PermissionsRequired) {
				break
			}

			whisper := !lo.Contains(userPermissions, "say_response")

			response, err := cmd.Cmd.Run(s, userPermissions, message)
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

			if _, err := s.Insert("chat_commands", map[string]any{
				"command":  cmd.Cmd.Command(),
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
}
