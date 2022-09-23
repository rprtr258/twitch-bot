package internal

import (
	"fmt"
	"strings"

	twitch "github.com/gempir/go-twitch-irc/v3"
)

func (s *Services) OnPrivateMessage(message twitch.PrivateMessage) {
	if message.User.Name == "rprtr258" && strings.HasPrefix(message.Message, intelCmdPrefix) {
		response, err := s.getIntelCmd(message)
		if err != nil {
			response = fmt.Sprintf("Internal error: %s", err.Error())
		}
		s.ChatClient.Reply(message.Channel, message.ID, response)
	}
}
