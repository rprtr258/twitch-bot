package cmds

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/rprtr258/twitch-bot/internal/message"
	"github.com/rprtr258/twitch-bot/internal/permissions"
	"github.com/rprtr258/twitch-bot/internal/services"

	"github.com/pocketbase/dbx"
)

type JoinCmd struct{}

func (JoinCmd) RequiredPermissions() []string {
	return []string{"global_admin"}
}

func (JoinCmd) Command() string {
	return "?join"
}

func (JoinCmd) Description() string {
	return "Join channel"
}

func (cmd JoinCmd) Run(ctx context.Context, s *services.Services, perms permissions.PermissionsList, msg message.TwitchMessage) (string, error) {
	words := strings.Split(msg.Text, " ")
	if len(words) != 2 {
		return fmt.Sprintf("Usage: %s <channel>", cmd.Command()), nil
	}

	channel := words[1]

	s.ChatClient.Join(channel)

	if _, err := s.Insert("joined_channels", map[string]any{
		"channel": channel,
	}); err != nil {
		// TODO: save log
		log.Println(err.Error())
	}

	return fmt.Sprintf("joined %s successfully", channel), nil
}

type LeaveCmd struct{}

func (LeaveCmd) RequiredPermissions() []string {
	return []string{"global_admin"}
}

func (LeaveCmd) Command() string {
	return "?leave"
}

func (LeaveCmd) Description() string {
	return "Leave channel"
}

func (cmd LeaveCmd) Run(ctx context.Context, s *services.Services, perms permissions.PermissionsList, msg message.TwitchMessage) (string, error) {
	words := strings.Split(msg.Text, " ")
	if len(words) != 2 {
		return fmt.Sprintf("Usage: %s <channel>", cmd.Command()), nil
	}

	channel := words[1]

	s.ChatClient.Depart(channel)

	db := s.Backend.DB()
	if _, err := db.Delete("joined_channels", dbx.NewExp("channel={:channel}", dbx.Params{"channel": channel})).Execute(); err != nil {
		// TODO: save log
		log.Println(err.Error())
	}

	return fmt.Sprintf("left %s successfully", channel), nil
}
