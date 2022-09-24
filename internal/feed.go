package internal

import (
	"database/sql"
	"fmt"
	"log"
	"strings"

	twitch "github.com/gempir/go-twitch-irc/v3"
	"github.com/pocketbase/dbx"
)

const (
	fedCmd = "?fed"
)

func weight(count int) int {
	return 100 + 13*count
}

func (s *Services) fed(message twitch.PrivateMessage) (string, error) {
	words := strings.Split(message.Message, " ")
	if len(words) != 1 && len(words) != 2 {
		return fmt.Sprintf(`Usage: "%[1]s" or "%[1]s <emote>"`, joinCmd), nil
	}

	db := s.Backend.DB()

	if len(words) == 2 {
		emote := words[1]

		var count int
		err := db.Select("count").From("feed_counts AS fc").InnerJoin("feed_emotes AS fe", dbx.NewExp("fc.emote_id = fe.id")).Where(dbx.And(
			dbx.NewExp("fe.emote={:emote}", dbx.Params{"emote": emote}),
			dbx.NewExp("fc.channel={:channel}", dbx.Params{"channel": message.Channel}),
			dbx.NewExp("fc.user_id={:user_id}", dbx.Params{"user_id": message.User.ID}),
		)).Row(&count)
		if err != nil && err != sql.ErrNoRows {
			// TODO: save log
			log.Println(err.Error())
			return "", err
		}

		// TODO: add feed after duration
		return fmt.Sprintf("Ты покормил %s %d раз. Теперь он весит %d грамм.", emote, count, weight(count)), nil
	}

	//	current_minute = datetime.now().minute
	//	emote = request.args.get('e')
	//	if emote == '0':
	//	    return "Ты можешь покормить: " + " ".join(EMOTES.keys())
	//	if emote not in EMOTES:
	//	    return f"{emote} никого не кормит FeelsWeakMan"
	//	encoded = request.args.get('c')
	//	user = request.args.get('u')
	//	if user == "Gekmi":
	//	    return "Гекми наказан за читерство с талончиками PauseFish"
	//	if len(encoded) > 1:
	//	    try:
	//	        username, the_emote, count, minute = unpack(encoded)
	//	    except Exception as e:
	//	        return "Некорректный талончик"
	//	    if the_emote[0] != emote:
	//	        return f"Ты используешь талончик {the_emote[0]} для {emote} . Так нельзя. Получить талончик: !feed {emote}"
	//	    if username != user:
	//	        return f"Ты используешь талончик {username}, но тебя зовут {user}. Это называется воровство. Ты приговариваешься к расстрелу AAUGH"
	//	    if minute <= current_minute < minute + 5:
	//	        return f"Ты еще не можешь покормить {EMOTES[emote][1]}"
	//	    count += 1
	//	    w = weight(count)
	//	    response = pack(username, the_emote[1][0], count, minute)
	//	    return f"Ты покормил {the_emote[1][1]} {count} раз. Теперь он весит {w} грамм. Талончик на следующую кормежку: {response}"
	//	response = pack(user, EMOTES[emote][0], 1, current_minute)
	//	print(unpack(response))
	//	w = weight(1)
	return fmt.Sprintf("Ты покормил %s в первый раз. Теперь он весит %d грамм.", "smth", weight(1)), nil
}
