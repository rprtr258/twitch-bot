package cmds

import (
	"fmt"
	"strconv"
	"strings"

	twitch "github.com/gempir/go-twitch-irc/v3"
	"github.com/karalef/balaboba"
)

// TODO: continue command
// TODO: read command

type BlabCmd struct{}

func (BlabCmd) Command() string {
	return "?blab"
}

func (BlabCmd) Description() string {
	return "Balaboba text generation neural network"
}

func (cmd BlabCmd) Run(s *Services, perms []string, message twitch.PrivateMessage) (string, error) {
	words := strings.Split(message.Message, " ")

	if len(words) == 1 {
		return fmt.Sprintf("Available styles: %d-standart, %d-user manual, %d-recipes, %d-short stories, %d-wikipedia simplified, %d-movie synopses, %d-folk wisdom",
			balaboba.Standart,
			balaboba.UserManual,
			balaboba.Recipes,
			balaboba.ShortStories,
			balaboba.WikipediaSipmlified,
			balaboba.MovieSynopses,
			balaboba.FolkWisdom,
		), nil
	}

	// remove command
	words = words[1:]

	styleIdx, err := strconv.Atoi(words[0])
	if err != nil {
		styleIdx = int(balaboba.Standart)
	} else {
		if styleIdx < 0 || styleIdx > int(balaboba.FolkWisdom) {
			return fmt.Sprintf("Invalid style, see %s for list of available styles", cmd.Command()), nil
		}
		// remove style
		words = words[1:]
	}

	text := strings.Join(words, " ")
	response, err := s.Balaboba.Generate(text, balaboba.Style(styleIdx))
	if err != nil {
		return "", err
	}

	id, err := s.Insert("blab", map[string]any{
		"text":          response.Text,
		"author_id":     message.User.ID,
		"continuations": 1,
	})
	if err != nil {
		return "", err
	}

	res := fmt.Sprintf("%s: %s", id, response.Text)
	// TODO: get back
	// if utf8.RuneCountInString(res) > _maxMessageLength {
	// 	return fmt.Sprintf("Читать продолжение в источнике: %s/blab/%s", s.Backend.Settings().Meta.AppUrl, id), nil
	// }

	return res, nil
}
