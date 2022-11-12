package internal

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/rprtr258/twitch-bot/internal/cmds"
	"github.com/rprtr258/twitch-bot/internal/message"
	"github.com/rprtr258/twitch-bot/internal/permissions"
	"github.com/rprtr258/twitch-bot/internal/services"

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

func (cmd CommandsCmd) Run(context.Context, *services.Services, permissions.PermissionsList, message.TwitchMessage) (string, error) {
	parts := lo.Map(allCommands, func(cmd Command2, _ int) string {
		return fmt.Sprintf("%s - %s", cmd.Cmd.Command(), cmd.Cmd.Description())
	})

	return strings.Join(parts, ", "), nil
}

type Command2 struct {
	PermissionsRequired []string
	Cmd                 Command
}

// TODO: permissions constants/store in db to dinamically update
// TODO: change to map
var allCommands = []Command2{{
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
	PermissionsRequired: []string{},
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

func OnPrivateMessage(s *services.Services) func(twitch.PrivateMessage) {
	return func(msg twitch.PrivateMessage) {
		msgMy := message.TwitchMessage{
			Text:    msg.Message,
			User:    msg.User,
			At:      msg.Time,
			Channel: msg.Channel,
		}
		s.LogMessage(msgMy)
		userPermissions := s.Permissions.GetPermissions(permissions.Claims{
			"username": msg.User.Name, // TODO: replace with user id
			"channel":  msg.Channel,
			// TODO: vips, moders
		})

		firstWord := strings.Split(msg.Message, " ")[0]
		userName := msg.User.Name

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

			response, err := cmd.Cmd.Run(ctx, s, userPermissions, msgMy)
			if err != nil {
				response = fmt.Sprintf("Internal error: %s", err.Error())
			} else if response == "" {
				response = "Empty response"
			}
			response = strings.ReplaceAll(response, "\n", " ")
			// TODO: move responding to commands
			// TODO: rewrite
			// TODO: fix response: Must be less than 500 character(s). somewhere
			if utf8.RuneCountInString(response) > message.MaxMessageLength {
				runes := []rune(response)
				log.Println(len(runes), utf8.RuneCountInString(response))
				if msg.User.Name == "rprtr258" { // TODO: change to permission check
					for i := 0; i < len(runes); i += message.MaxMessageLength {
						resp := string(lo.Subset(runes, i, message.MaxMessageLength))
						log.Println("trying to send msg of len(1)", utf8.RuneCountInString(resp))
						if whisper {
							s.ChatClient.Whisper(msg.User.Name, resp)
						} else {
							s.ChatClient.Say(msg.Channel, resp)
						}
						// TODO: check spam limit
						time.Sleep(time.Millisecond * 500)
					}
				} else {
					response = string(runes[:message.MaxMessageLength])

					// TODO: fix cutting
					log.Println("trying to send msg of len(2)", utf8.RuneCountInString(response))
					if whisper {
						s.ChatClient.Whisper(msg.User.Name, response)
					} else {
						s.ChatClient.Reply(msg.Channel, msg.ID, response)
					}
				}
			} else {
				log.Println("trying to send msg of len(3)", utf8.RuneCountInString(response))
				if whisper {
					s.ChatClient.Whisper(msg.User.Name, response)
				} else {
					s.ChatClient.Reply(msg.Channel, msg.ID, response)
				}
			}

			// TODO: fix not logging blab cmds
			if _, err := s.Insert("chat_commands", map[string]any{
				"command":  cmd.Cmd.Command(),
				"args":     msg.Message,
				"at":       msg.Time,
				"response": response,
				"user":     userName,
				"channel":  msg.Channel,
			}); err != nil {
				// TODO: save log/sentry
				// TODO: change logger to zap
				log.Println(err.Error())
			}

			break
		}
	}
}

func OnWhisperMessage(s *services.Services) func(twitch.WhisperMessage) {
	return func(msg twitch.WhisperMessage) {
		msgMy := message.TwitchMessage{
			Text:    msg.Message,
			User:    msg.User,
			At:      time.Now(),
			Channel: msg.User.Name,
		}
		s.LogMessage(msgMy)
		userPermissions := s.Permissions.GetPermissions(permissions.Claims{
			"username": msg.User.Name, // TODO: replace with user id
			"channel":  msgMy.Channel,
		})

		firstWord := strings.Split(msg.Message, " ")[0]
		userName := msg.User.Name

		for _, cmd := range allCommands {
			if cmd.Cmd.Command() != firstWord {
				continue
			}

			// TODO: add ban permissions
			if !userPermissions.Has(cmd.PermissionsRequired...) {
				break
			}

			ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
			defer cancel()

			response, err := cmd.Cmd.Run(ctx, s, userPermissions, msgMy)
			if err != nil {
				response = fmt.Sprintf("Internal error: %s", err.Error())
			} else if response == "" {
				response = "Empty response"
			}
			response = strings.ReplaceAll(response, "\n", " ")
			// TODO: move responding to commands
			// TODO: rewrite
			// TODO: fix response: Must be less than 500 character(s). somewhere
			if utf8.RuneCountInString(response) > message.MaxMessageLength {
				runes := []rune(response)
				log.Println(len(runes), utf8.RuneCountInString(response))
				if msg.User.Name == "rprtr258" { // TODO: change to permission check
					for i := 0; i < len(runes); i += message.MaxMessageLength {
						resp := string(lo.Subset(runes, i, message.MaxMessageLength))
						log.Println("trying to send msg of len(1)", utf8.RuneCountInString(resp))
						s.ChatClient.Whisper(msg.User.Name, resp)
						// TODO: check spam limit
						time.Sleep(time.Millisecond * 500)
					}
				} else {
					response = string(runes[:message.MaxMessageLength])

					// TODO: fix cutting
					log.Println("trying to send msg of len(2)", utf8.RuneCountInString(response))
					s.ChatClient.Whisper(msg.User.Name, response)
				}
			} else {
				log.Println("trying to send msg of len(3)", utf8.RuneCountInString(response))
				s.ChatClient.Whisper(msg.User.Name, response)
			}

			// TODO: fix not logging blab cmds
			if _, err := s.Insert("chat_commands", map[string]any{
				"command":  cmd.Cmd.Command(),
				"args":     msg.Message,
				"at":       time.Now(),
				"response": response,
				"user":     userName,
				"channel":  msgMy.Channel,
			}); err != nil {
				// TODO: save log/sentry
				// TODO: change logger to zap
				log.Println(err.Error())
			}

			break
		}
	}
}
