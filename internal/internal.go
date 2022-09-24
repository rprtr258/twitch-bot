package internal

import (
	"fmt"
	"log"
	"strings"
	"time"

	twitch "github.com/gempir/go-twitch-irc/v3"
	"github.com/pocketbase/dbx"
)

const (
	joinCmd  = "?join"
	leaveCmd = "?leave"
)

func (s *Services) join(message twitch.PrivateMessage) (string, error) {
	words := strings.Split(message.Message, " ")
	if len(words) != 2 {
		return fmt.Sprintf("Usage: %s <channel>", joinCmd), nil
	}

	channel := words[1]

	s.ChatClient.Join(channel)

	db := s.Backend.DB()
	if _, err := db.Insert("joined_channels", dbx.Params{
		"channel": channel,
	}).Execute(); err != nil {
		// TODO: save log
		log.Println(err.Error())
	}

	return fmt.Sprintf("joined %s successfully", channel), nil
}

func (s *Services) leave(message twitch.PrivateMessage) (string, error) {
	words := strings.Split(message.Message, " ")
	if len(words) != 2 {
		return fmt.Sprintf("Usage: %s <channel>", joinCmd), nil
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

func (s *Services) OnPrivateMessage(message twitch.PrivateMessage) {
	db := s.Backend.DB()
	if _, err := db.Insert("messages", dbx.Params{
		"user_id":           message.User.ID,
		"message":           message.Message,
		"at":                message.Time,
		"channel":           message.Channel,
		"user_name":         message.User.Name,
		"user_display_name": message.User.DisplayName,
	}).Execute(); err != nil {
		// TODO: save log
		log.Println(err.Error())
	}

	cmd := strings.Split(message.Message, " ")[0]
	if message.User.Name == "rprtr258" && cmd == intelCmd {
		response, err := s.getIntelCmd(message)
		if err != nil {
			response = fmt.Sprintf("Internal error: %s", err.Error())
		}
		s.ChatClient.Whisper(message.User.Name, response)
		db := s.Backend.DB()
		if _, err := db.Insert("chat_commands", dbx.Params{
			"command":  intelCmd,
			"args":     message.Message,
			"at":       time.Now(),
			"response": response,
			"user":     message.User.Name,
			"channel":  message.Channel,
		}).Execute(); err != nil {
			// TODO: save log
			log.Println(err.Error())
		}
	}
	if message.User.Name == "rprtr258" && cmd == joinCmd {
		response, err := s.join(message)
		if err != nil {
			response = fmt.Sprintf("Internal error: %s", err.Error())
		}
		s.ChatClient.Reply(message.Channel, message.ID, response)
		db := s.Backend.DB()
		if _, err := db.Insert("chat_commands", dbx.Params{
			"command":  joinCmd,
			"args":     message.Message,
			"at":       time.Now(),
			"response": response,
			"user":     message.User.Name,
			"channel":  message.Channel,
		}).Execute(); err != nil {
			// TODO: save log
			log.Println(err.Error())
		}
	}
	if message.User.Name == "rprtr258" && cmd == leaveCmd {
		response, err := s.leave(message)
		if err != nil {
			response = fmt.Sprintf("Internal error: %s", err.Error())
		}
		s.ChatClient.Reply(message.Channel, message.ID, response)
		db := s.Backend.DB()
		if _, err := db.Insert("chat_commands", dbx.Params{
			"command":  leaveCmd,
			"args":     message.Message,
			"at":       time.Now(),
			"response": response,
			"user":     message.User.Name,
			"channel":  message.Channel,
		}).Execute(); err != nil {
			// TODO: save log
			log.Println(err.Error())
		}
	}
}
