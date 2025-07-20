package fun

import (
	"fmt"
	"os"

	. "github.com/0supa/supa8/fun"
	. "github.com/0supa/supa8/fun/api/twitch"
	"github.com/0supa/supa8/fun/utils"
	"github.com/gempir/go-twitch-irc/v4"
)

func init() {
	Fun.Register(&Cmd{
		Name: "reload_blocklist",
		Handler: func(m twitch.PrivateMessage) (err error) {
			if !utils.IsPrivileged(m.User.ID) || m.Message != "`rb" {
				return
			}

			blockedOld := len(Fun.BlockedUserIDs)
			LoadBlocklist()
			blockedNew := len(Fun.BlockedUserIDs)

			diff := blockedNew - blockedOld
			var prefix string
			if diff >= 0 {
				prefix = "+"
			}

			_, err = Say(m.RoomID, fmt.Sprintf("%v blocked users (%s%v)", blockedNew, prefix, diff), m.ID)
			os.Exit(0)
			return
		},
	})
}
