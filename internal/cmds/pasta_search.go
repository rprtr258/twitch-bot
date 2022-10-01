package cmds

import (
	"abobus/internal/services"
	"bufio"
	"errors"
	"fmt"
	"math/rand"
	"os"
	"strings"
	"time"

	twitch "github.com/gempir/go-twitch-irc/v3"
	"github.com/lithammer/fuzzysearch/fuzzy"
	"github.com/samber/lo"
)

type PastaSearchCmd struct{}

func (PastaSearchCmd) Command() string {
	return "?pasta"
}

func (PastaSearchCmd) Description() string {
	return "Search for copypasta"
}

// TODO: cmd to add pastas
func (cmd PastaSearchCmd) Run(s *services.Services, perms []string, message twitch.PrivateMessage) (string, error) {
	// TODO: check empty query
	query := strings.TrimPrefix(message.Message, cmd.Command()+" ")

	// TODO: move pastes to database
	file, err := os.Open("pastes.txt")
	if err != nil {
		return "", err
	}
	defer file.Close()

	pastes := []string{}
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		pastes = append(pastes, scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		return "", err
	}

	matches := lo.Filter(fuzzy.RankFindNormalized(query, pastes), func(match fuzzy.Rank, _ int) bool {
		return match.Distance < 500
	})
	totalDistance := lo.SumBy(matches, func(match fuzzy.Rank) int {
		return match.Distance
	})
	if totalDistance <= 0 {
		return "No pastes found", nil
	}

	rand.Seed(time.Now().Unix())
	rng := rand.Intn(totalDistance)
	for _, match := range matches {
		rng -= match.Distance
		if rng <= 0 {
			return fmt.Sprintf("%d: %s", match.Distance, match.Target), nil
		}
	}
	return "", errors.New("SOMETHING VERY UNLIKELY HAPPENED PLEASE CONTACT @rprtr258 TO FIX THIS")
}
