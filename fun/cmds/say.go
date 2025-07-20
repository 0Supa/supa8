package fun

import (
	"strings"

	. "github.com/0supa/supa8/fun"
	. "github.com/0supa/supa8/fun/api/twitch"
	"github.com/0supa/supa8/fun/utils"
	"github.com/gempir/go-twitch-irc/v4"
)

func init() {
	Fun.Register(&Cmd{
		Name: "say",
		Handler: func(m twitch.PrivateMessage) (err error) {
			if !utils.IsPrivileged(m.User.ID) {
				return
			}

			args := strings.Split(m.Message, " ")
			if args[0] == "`echo" && len(args) >= 2 {
				parent := ""
				if m.Reply != nil {
					parent = m.Reply.ParentMsgID
				}
				_, err = Say(m.RoomID, strings.Join(args[1:], " "), parent)
				return
			}

			if args[0] != "`say" || len(args) < 3 {
				return
			}

			user, err := GetUser(args[1], "")
			if err != nil {
				_, err = Say(m.RoomID, "get user: "+err.Error(), m.ID)
				return
			}

			if user.ID == "" {
				_, err = Say(m.RoomID, "user not found", m.ID)
				return
			}

			res, err := Say(user.ID, strings.Join(args[2:], " "), "")
			if err != nil {
				_, err = Say(m.RoomID, err.Error(), m.ID)
				return
			}

			_, err = Say(m.RoomID, res.Data.Mutation.Message.ID, m.ID)
			return
		},
	})
}
