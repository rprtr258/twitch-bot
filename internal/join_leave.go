package internal

import (
	"fmt"
	"log"
	"strings"

	twitch "github.com/gempir/go-twitch-irc/v3"
	"github.com/pocketbase/dbx"
)

const (
	joinCmd  = "?join"
	leaveCmd = "?leave"
)

func (s *Services) join(perms []string, message twitch.PrivateMessage) (string, error) {
	words := strings.Split(message.Message, " ")
	if len(words) != 2 {
		return fmt.Sprintf("Usage: %s <channel>", joinCmd), nil
	}

	channel := words[1]

	s.ChatClient.Join(channel)

	if _, err := s._insert("joined_channels", map[string]any{
		"channel": channel,
	}); err != nil {
		// TODO: save log
		log.Println(err.Error())
	}

	return fmt.Sprintf("joined %s successfully", channel), nil
}

func (s *Services) leave(perms []string, message twitch.PrivateMessage) (string, error) {
	words := strings.Split(message.Message, " ")
	if len(words) != 2 {
		return fmt.Sprintf("Usage: %s <channel>", leaveCmd), nil
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
