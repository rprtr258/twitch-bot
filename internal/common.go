package internal

import (
	"context"

	"github.com/rprtr258/twitch-bot/internal/message"
	"github.com/rprtr258/twitch-bot/internal/permissions"
	"github.com/rprtr258/twitch-bot/internal/services"
)

type Command interface {
	// RequiredPermissions - minimal permissions required to use the command
	RequiredPermissions() []string
	// Command - get command name
	Command() string
	// Description - get command description
	Description() string
	// Run - execute command
	Run(context.Context, *services.Services, permissions.PermissionsList, message.TwitchMessage) (string, error)
}
