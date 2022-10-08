package cmds

import (
	"abobus/internal/services"

	twitch "github.com/gempir/go-twitch-irc/v3"
)

const MaxMessageLength = 500

type Command interface {
	// Command - get command name
	Command() string
	// Description - get command description
	Description() string
	// Run - execute command
	Run(*services.Services, []string, twitch.PrivateMessage) (string, error)
}
