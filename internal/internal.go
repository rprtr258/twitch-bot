package internal

import (
	"fmt"
	"log"
	"strings"
	"time"

	twitch "github.com/gempir/go-twitch-irc/v3"
	"github.com/pocketbase/dbx"
)

func (s *Services) OnPrivateMessage(message twitch.PrivateMessage) {
	if message.User.Name == "rprtr258" && strings.HasPrefix(message.Message, intelCmdPrefix) {
		response, err := s.getIntelCmd(message)
		if err != nil {
			response = fmt.Sprintf("Internal error: %s", err.Error())
		}
		s.ChatClient.Reply(message.Channel, message.ID, response)
		db := s.Backend.DB()
		if _, err := db.Insert("chat_commands", dbx.Params{
			"command":  intelCmdPrefix,
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
