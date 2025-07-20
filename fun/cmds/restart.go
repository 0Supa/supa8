package fun

import (
	"os"

	. "github.com/0supa/supa8/fun"
	. "github.com/0supa/supa8/fun/api/twitch"
	"github.com/0supa/supa8/fun/utils"
	"github.com/gempir/go-twitch-irc/v4"
)

func init() {
	Fun.Register(&Cmd{
		Name: "restart",
		Handler: func(m twitch.PrivateMessage) (err error) {
			if !utils.IsPrivileged(m.User.ID) || m.Message != "`restart" {
				return
			}

			_, err = Say(m.RoomID, "exiting...", m.ID)
			os.Exit(0)
			return
		},
	})
}
