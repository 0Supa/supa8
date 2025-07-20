package fun

import (
	"fmt"
	"strings"

	. "github.com/0supa/supa8/fun"
	. "github.com/0supa/supa8/fun/api/twitch"
	"github.com/gempir/go-twitch-irc/v4"
)

func init() {
	Fun.Register(&Cmd{
		Name: "logs",
		Handler: func(m twitch.PrivateMessage) (err error) {
			args := strings.Split(m.Message, " ")
			if len(args) < 1 || args[0] != "`logs" {
				return
			}

			c := m.Channel
			u := m.User.Name
			if len(args) >= 2 {
				u = args[1]
			}
			if len(args) >= 3 {
				c = args[2]
			}

			_, err = Say(m.RoomID, fmt.Sprintf("https://tv.supa.sh/logs?c=%s&u=%s", c, u), m.ID)
			return
		},
	})
}
