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
	everyoneRelation  = ""
	executeRelation   = "execute"
	forbiddenRelation = "forbidden"
)

type command struct {
	Relation string
	Command  string
	Run      func(*Services, twitch.PrivateMessage) (string, error)
	Whisper  bool
}

var cmds []command = []command{{
	Relation: executeRelation,
	Command:  intelCmd,
	Run:      (*Services).getIntelCmd,
	Whisper:  true,
}, {
	Relation: executeRelation,
	Command:  joinCmd,
	Run:      (*Services).join,
}, {
	Relation: executeRelation,
	Command:  leaveCmd,
	Run:      (*Services).leave,
}, {
	Relation: everyoneRelation,
	Command:  fedCmd,
	Run:      (*Services).fed,
}}

func (s *Services) logMessage(message twitch.PrivateMessage) {
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
}

func (s *Services) OnPrivateMessage(message twitch.PrivateMessage) {
	db := s.Backend.DB()

	s.logMessage(message)

	firstWord := strings.Split(message.Message, " ")[0]
	userName := message.User.Name

	var relation string
	if err := db.Select("relation").From("permissions").Where(dbx.And(
		dbx.NewExp("user_id={:user_id}", dbx.Params{"user_id": message.User.ID}),
		dbx.NewExp("object={:cmd}", dbx.Params{"cmd": fmt.Sprintf("cmd:%s", firstWord)}),
	)).Build().Row(&relation); err != nil {
		if err != sql.ErrNoRows {
			log.Println(err.Error())
			return
		}
	}

	for _, cmd := range cmds {
		if cmd.Command != firstWord {
			continue
		}

		if relation == forbiddenRelation {
			break
		}

		if cmd.Relation != everyoneRelation && cmd.Relation != relation {
			break
		}

		// TODO: peredelat'
		if cmd.Command == fedCmd && message.Channel != "rprtr258" {
			break
		}

		// TODO: peredelat'
		whisper := cmd.Whisper
		if cmd.Command == intelCmd && message.Channel == "rprtr258" {
			whisper = false
		}

		response, err := cmd.Run(s, message)
		if err != nil {
			response = fmt.Sprintf("Internal error: %s", err.Error())
		}

		if whisper {
			s.ChatClient.Whisper(message.User.Name, response)
		} else {
			s.ChatClient.Reply(message.Channel, message.ID, response)
		}

		if _, err := db.Insert("chat_commands", dbx.Params{
			"command":  cmd.Command,
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
