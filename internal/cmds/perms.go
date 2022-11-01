package cmds

import (
	"context"
	"fmt"
	"strings"
	"unicode/utf8"

	twitch "github.com/gempir/go-twitch-irc/v3"
	"github.com/samber/lo"

	"github.com/rprtr258/twitch-bot/internal/permissions"
	"github.com/rprtr258/twitch-bot/internal/services"
)

type PermsCmd struct{}

func (PermsCmd) Command() string {
	return "?perms"
}

func (PermsCmd) Description() string {
	return "Get permissions"
}

func (cmd PermsCmd) Run(ctx context.Context, s *services.Services, perms permissions.PermissionsList, message twitch.PrivateMessage) (string, error) {
	words := strings.Split(message.Message, " ")
	permissionsToShow := perms
	if len(words) > 1 {
		claims := permissions.Claims{}
		for _, word := range words[1:] {
			idx := strings.IndexRune(word, ':')
			if idx == -1 || idx == 0 || idx == utf8.RuneCountInString(word)-1 {
				return fmt.Sprintf(`Usage: "%[1]s" or "%[1]s key:value key:value"`, cmd.Command()), nil
			}
			key, value := word[:idx], word[idx+1:]
			claims[key] = value
		}
		permissionsToShow = s.Permissions.GetPermissions(claims)
	}
	return strings.Join(lo.Keys(permissionsToShow), ","), nil
}
