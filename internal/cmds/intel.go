package cmds

import (
	"context"
	"errors"
	"fmt"
	"math"
	"strings"
	"time"

	"github.com/rprtr258/twitch-bot/internal/message"
	"github.com/rprtr258/twitch-bot/internal/permissions"
	"github.com/rprtr258/twitch-bot/internal/services"

	"github.com/nicklaw5/helix/v2"
)

func formatDuration(d time.Time) string {
	var parts []string
	now := time.Now()

	{
		var years int
		for d.AddDate(1, 0, 0).Before(now) {
			years++
			d = d.AddDate(1, 0, 0)
		}
		if years > 0 {
			parts = append(parts, fmt.Sprintf("%dy", years))
		}
	}

	{
		var months int
		for d.AddDate(0, 1, 0).Before(now) {
			months++
			d = d.AddDate(0, 1, 0)
		}
		if months > 0 {
			parts = append(parts, fmt.Sprintf("%dmo", months))
		}
	}

	{
		var days int
		for d.AddDate(0, 0, 1).Before(now) {
			days++
			d = d.AddDate(0, 0, 1)
		}
		if days > 0 {
			parts = append(parts, fmt.Sprintf("%dd", days))
		}
	}

	minutes := (int64)(math.Floor(now.Sub(d).Minutes()))
	if minutes/60 > 0 {
		parts = append(parts, fmt.Sprintf("%dh", minutes/60))
	}

	parts = append(parts, fmt.Sprintf("%dm", minutes%60))

	return strings.Join(parts, "")
}

type IntelCmd struct{}

func (IntelCmd) RequiredPermissions() []string {
	return []string{}
}

func (IntelCmd) Command() string {
	return "?intel"
}

func (IntelCmd) Description() string {
	return "Gather intel on user"
}

func (IntelCmd) Run(ctx context.Context, s *services.Services, perms permissions.PermissionsList, msg message.TwitchMessage) (string, error) {
	words := strings.Split(msg.Text, " ")
	if len(words) < 2 {
		return "No username provided", nil
	}

	login := words[1]
	resp, err := s.TwitchApiClient.GetUsers(&helix.UsersParams{
		Logins: []string{login},
	})
	if err != nil {
		return "", err
	}
	if resp.ErrorMessage != "" {
		return "", errors.New(resp.ErrorMessage)
	}

	users := resp.Data.Users
	if len(users) == 0 {
		return fmt.Sprintf("%s not found", login), nil
	}

	user := users[0]
	createdAt := user.CreatedAt.Time
	delta := formatDuration(createdAt)
	return fmt.Sprintf(
		"created=%s ago=%s broadcaster_type=%q id=%s login=%q display_name=%q view_count=%d",
		createdAt.Format("15:04.2.1.2006"),
		delta,
		user.BroadcasterType,
		user.ID,
		user.Login,
		user.DisplayName,
		user.ViewCount,
	), nil
}
