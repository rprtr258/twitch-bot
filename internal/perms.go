package internal

import (
	"fmt"
	"strings"
	"unicode/utf8"

	twitch "github.com/gempir/go-twitch-irc/v3"
)

const (
	permsCmd = "?perms"
)

func (s *Services) getPermissions(perms []string, message twitch.PrivateMessage) (string, error) {
	words := strings.Split(message.Message, " ")
	permissionsToShow := perms
	if len(words) > 1 {
		claims := Claims{}
		for _, word := range words[1:] {
			idx := strings.IndexRune(word, ':')
			if idx == -1 || idx == 0 || idx == utf8.RuneCountInString(word)-1 {
				return fmt.Sprintf(`Usage: "%[1]s" or "%[1]s key:value key:value"`, permsCmd), nil
			}
			key, value := word[:idx], word[idx+1:]
			claims[key] = value
		}
		permissionsToShow = s.Permissions.getPermissions(claims)
	}
	return strings.Join(permissionsToShow, ","), nil
}
