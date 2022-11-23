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

func (CommandsCmd) RequiredPermissions() []string {
	return []string{"execute_commands"}
}

func (CommandsCmd) Command() string {
	return "?commands"
}

func (CommandsCmd) Description() string {
	return "List all commands"
}

func (cmd CommandsCmd) Run(context.Context, *services.Services, permissions.PermissionsList, message.TwitchMessage) (string, error) {
	parts := lo.Map(allCommands, func(cmd Command, _ int) string {
		return fmt.Sprintf("%s - %s", cmd.Command(), cmd.Description())
	})

	return strings.Join(parts, ", "), nil
}

// TODO: permissions constants/store in db to dinamically update
// TODO: change to map
var allCommands = []Command{
	cmds.IntelCmd{},
	cmds.JoinCmd{},
	cmds.LeaveCmd{},
	cmds.FedCmd{},
	cmds.BlabGenCmd{},
	cmds.BlabContinueCmd{},
	cmds.BlabReadCmd{},
	&cmds.PythCmd{},
	CommandsCmd{},
	cmds.PermsCmd{},
	cmds.PastaSearchCmd{},
}

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
			if cmd.Command() != firstWord {
				continue
			}

			// TODO: add ban permissions
			if !userPermissions.Has(cmd.RequiredPermissions()...) {
				break
			}

			whisper := !userPermissions.Has("say_response")

			ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
			defer cancel()

			response, err := cmd.Run(ctx, s, userPermissions, msgMy)
			if err != nil {
				response = fmt.Sprintf("Internal error: %s", err.Error())
			} else if response == "" {
				response = "Empty response"
			}
			response = strings.ReplaceAll(response, "\n", " ")

			// TODO: move responding to commands
			// TODO: rewrite
			// TODO: fix response: Must be less than 500 character(s). somewhere
			// TODO: change to permission check, to send message by parts
			croppedResponse := cropMessage(response)
			log.Println("trying to send msg of len(2)", utf8.RuneCountInString(croppedResponse))
			if whisper {
				s.ChatClient.Whisper(msg.User.Name, croppedResponse)
			} else {
				s.ChatClient.Reply(msg.Channel, msg.ID, croppedResponse)
			}

			// TODO: fix not logging blab cmds
			if _, err := s.Insert("chat_commands", map[string]any{
				"command":  cmd.Command(),
				"args":     msg.Message,
				"response": croppedResponse,
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
			if cmd.Command() != firstWord {
				continue
			}

			// TODO: add ban permissions
			if !userPermissions.Has(cmd.RequiredPermissions()...) {
				break
			}

			ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
			defer cancel()

			response, err := cmd.Run(ctx, s, userPermissions, msgMy)
			if err != nil {
				response = fmt.Sprintf("Internal error: %s", err.Error())
			} else if response == "" {
				response = "Empty response"
			}
			response = strings.ReplaceAll(response, "\n", " ")

			// TODO: move responding to commands
			// TODO: rewrite
			// TODO: fix response: Must be less than 500 character(s). somewhere
			// TODO: change to permission check, to send message by parts
			croppedResponse := cropMessage(response)
			log.Println("trying to whisper msg of len(3)", utf8.RuneCountInString(croppedResponse))
			s.ChatClient.Whisper(msg.User.Name, croppedResponse)

			// TODO: fix not logging blab cmds
			if _, err := s.Insert("chat_commands", map[string]any{
				"command":  cmd.Command(),
				"args":     msg.Message,
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

// TODO: fix cutting
func cropMessage(response string) string {
	if utf8.RuneCountInString(response) <= message.MaxMessageLength {
		return response
	}

	runes := []rune(response)
	return string(runes[:message.MaxMessageLength])
}
