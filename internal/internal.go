package internal

import (
	"abobus/internal/cmds"
	"abobus/internal/permissions"
	"abobus/internal/services"
	"context"
	"fmt"
	"log"
	"strings"
	"time"
	"unicode/utf8"

	twitch "github.com/gempir/go-twitch-irc/v3"
	"github.com/samber/lo"
)

// TODO: list only available commands
type CommandsCmd struct{}

func (CommandsCmd) Command() string {
	return "?commands"
}

func (CommandsCmd) Description() string {
	return "List all commands"
}

func (CommandsCmd) Run(ctx context.Context, s *services.Services, perms permissions.PermissionsList, message twitch.PrivateMessage) (string, error) {
	parts := lo.Map(allCommands, func(cmd Command, _ int) string {
		return fmt.Sprintf("%s - %s", cmd.Cmd.Command(), cmd.Cmd.Description())
	})

	return strings.Join(parts, ", "), nil
}

type Command struct {
	PermissionsRequired []string
	Cmd                 cmds.Command
}

// TODO: permissions constants/store in db to dinamically update
// TODO: change to map
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
	Cmd:                 cmds.BlabGenCmd{},
}, {
	PermissionsRequired: []string{"execute_commands"},
	Cmd:                 cmds.BlabContinueCmd{},
}, {
	PermissionsRequired: []string{"execute_commands"},
	Cmd:                 cmds.BlabReadCmd{},
}, {
	PermissionsRequired: []string{"execute_commands"},
	Cmd:                 &cmds.PythCmd{},
}, {
	PermissionsRequired: []string{"execute_commands"},
	Cmd:                 CommandsCmd{},
}, {
	PermissionsRequired: []string{},
	Cmd:                 cmds.PermsCmd{},
}, {
	PermissionsRequired: []string{},
	Cmd:                 cmds.PastaSearchCmd{},
}}

// TODO: handle whispers
func OnPrivateMessage(s *services.Services) func(twitch.PrivateMessage) {
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
			if !userPermissions.Has(cmd.PermissionsRequired...) {
				break
			}

			whisper := !userPermissions.Has("say_response")

			ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
			defer cancel()

			response, err := cmd.Cmd.Run(ctx, s, userPermissions, message)
			if err != nil {
				response = fmt.Sprintf("Internal error: %s", err.Error())
			} else if response == "" {
				response = "Empty response"
			}
			response = strings.ReplaceAll(response, "\n", " ")
			// TODO: move responding to commands
			// TODO: rewrite
			// TODO: fix response: Must be less than 500 character(s). somewhere
			if utf8.RuneCountInString(response) > cmds.MaxMessageLength {
				runes := []rune(response)
				log.Println(len(runes), utf8.RuneCountInString(response))
				if message.User.Name == "rprtr258" { // TODO: change to permission check
					for i := 0; i < len(runes); i += cmds.MaxMessageLength {
						resp := string(lo.Subset(runes, i, cmds.MaxMessageLength))
						log.Println("trying to send msg of len(1)", utf8.RuneCountInString(resp))
						if whisper {
							s.ChatClient.Whisper(message.User.Name, resp)
						} else {
							s.ChatClient.Say(message.Channel, resp)
						}
						// TODO: check spam limit
						time.Sleep(time.Millisecond * 500)
					}
				} else {
					response = string(runes[:cmds.MaxMessageLength])

					// TODO: fix cutting
					log.Println("trying to send msg of len(2)", utf8.RuneCountInString(response))
					if whisper {
						s.ChatClient.Whisper(message.User.Name, response)
					} else {
						s.ChatClient.Reply(message.Channel, message.ID, response)
					}
				}
			} else {
				log.Println("trying to send msg of len(3)", utf8.RuneCountInString(response))
				if whisper {
					s.ChatClient.Whisper(message.User.Name, response)
				} else {
					s.ChatClient.Reply(message.Channel, message.ID, response)
				}
			}

			// TODO: fix not logging blab cmds
			if _, err := s.Insert("chat_commands", map[string]any{
				"command":  cmd.Cmd.Command(),
				"args":     message.Message,
				"at":       time.Now(),
				"response": response,
				"user":     userName,
				"channel":  message.Channel,
			}); err != nil {
				// TODO: save log/sentry
				// TODO: change logger to zap
				log.Println(err.Error())
			}

			break
		}
	}
}
