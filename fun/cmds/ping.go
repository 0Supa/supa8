package fun

import (
	"fmt"
	"math/rand"
	"runtime"
	"time"

	"github.com/0supa/supa8/config"
	. "github.com/0supa/supa8/fun"
	. "github.com/0supa/supa8/fun/api/twitch"
	"github.com/gempir/go-twitch-irc/v4"
)

func init() {
	Fun.Register(&Cmd{
		Name: "ping",
		Handler: func(m twitch.PrivateMessage) (err error) {
			if m.Message != "`ping" {
				return
			}

			var mem runtime.MemStats
			runtime.ReadMemStats(&mem)

			e := []string{":3", ":x", "c:", ":$", ">.<", "^.^", "^o^"}
			_, err = Say(m.RoomID, fmt.Sprintf("%s %vms, %s, %vMiB, up:%s, channels:%v, blocked:%v",
				e[rand.Intn(len(e))],
				time.Since(m.Time).Milliseconds(),
				runtime.Version(),
				mem.Alloc/1024/1024,
				time.Since(InitTime).Truncate(time.Second),
				len(config.Meta.Channels),
				len(Fun.BlockedUserIDs)), m.ID)
			return
		},
	})
}
