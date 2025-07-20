package fun

import (
	logger "log"
	"time"

	api_twitch "github.com/0supa/supa8/fun/api/twitch"
	"github.com/gempir/go-twitch-irc/v4"
)

type Cmd struct {
	Name    string
	Handler func(m twitch.PrivateMessage) error
}

type FunData struct {
	Cmds           []Cmd
	BlockedUserIDs map[string]struct{}
}

func (f *FunData) Register(c *Cmd) {
	f.Cmds = append(f.Cmds, *c)
}

var Fun = FunData{
	Cmds:           []Cmd{},
	BlockedUserIDs: map[string]struct{}{},
}

var InitTime = time.Now()

func LoadBlocklist() {
	user, err := api_twitch.GetOwner()
	if err != nil {
		logger.Panicln("failed getting owner account", err)
		return
	}

	if user.BlockedUsers != nil {
		clear(Fun.BlockedUserIDs)
		for _, u := range *user.BlockedUsers {
			Fun.BlockedUserIDs[u.ID] = struct{}{}
		}
	}
}

func IsBlocked(userID string) bool {
	_, exists := Fun.BlockedUserIDs[userID]
	return exists
}
