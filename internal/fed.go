package internal

import (
	"database/sql"
	"fmt"
	"log"
	"strconv"
	"strings"

	twitch "github.com/gempir/go-twitch-irc/v3"
	"github.com/pocketbase/dbx"
)

const (
	fedCmd = "?fed"
)

func (s *Services) fed(message twitch.PrivateMessage) (string, error) {
	userID, err := strconv.Atoi(message.User.ID)
	if err != nil {
		return "", err
	}

	words := strings.Split(message.Message, " ")
	if len(words) != 2 {
		return fmt.Sprintf(`Usage: "%s <emote>"`, fedCmd), nil
	}

	db := s.Backend.DB()

	theWord := words[1]

	var count int
	err = db.Select("COUNT(*) AS count").From("messages").Where(dbx.And(
		dbx.NewExp("channel={:channel}", dbx.Params{"channel": message.Channel}),
		// TODO: fix theWord = '%%' escaping
		dbx.Like("message", theWord),
		dbx.NotLike("message", fedCmd).Match(false, true),
		dbx.NewExp("user_id={:user_id}", dbx.Params{"user_id": userID}),
	)).Row(&count)
	if err != nil && err != sql.ErrNoRows {
		// TODO: save log
		log.Println(err.Error())
		return "", err
	}

	// TODO: add feed after duration
	return fmt.Sprintf("Ты использовал слово %s %d раз.", theWord, count), nil
}
