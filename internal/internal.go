package internal

import (
	"database/sql"
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

type command struct {
	Relation string
	Command  string
	Run      func(*Services, twitch.PrivateMessage) (string, error)
}

func (s *Services) OnPrivateMessage(message twitch.PrivateMessage) {
	cmds := []command{{
		Relation: "execute",
		Command:  intelCmd,
		Run:      (*Services).getIntelCmd,
	}, {
		Relation: "execute",
		Command:  joinCmd,
		Run:      (*Services).join,
	}, {
		Relation: "execute",
		Command:  leaveCmd,
		Run:      (*Services).leave,
	}}

	db := s.Backend.DB()

	// TODO: users table
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
	userName := message.User.Name

	var relation string
	if err := db.Select("relation").From("permissions").Where(dbx.And(
		dbx.NewExp("user_id={:user_id}", dbx.Params{"user_id": message.User.ID}),
		dbx.NewExp("object={:cmd}", dbx.Params{"cmd": fmt.Sprintf("cmd:%s", cmd)}),
	)).Build().Row(&relation); err != nil {
		if err == sql.ErrNoRows {
			return
		}
		log.Println(err.Error())
	}

	for _, command := range cmds {
		if command.Command != cmd || command.Relation != relation {
			continue
		}

		response, err := command.Run(s, message)
		if err != nil {
			response = fmt.Sprintf("Internal error: %s", err.Error())
		}

		s.ChatClient.Whisper(message.User.Name, response)
		if _, err := db.Insert("chat_commands", dbx.Params{
			"command":  command.Command,
			"args":     message.Message,
			"at":       time.Now(),
			"response": response,
			"user":     userName,
			"channel":  message.Channel,
		}).Execute(); err != nil {
			// TODO: save log
			log.Println(err.Error())
		}

		break
	}
}
