package fun

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	. "github.com/0supa/supa8/fun"
	"github.com/0supa/supa8/fun/api"
	. "github.com/0supa/supa8/fun/api/twitch"
	"github.com/0supa/supa8/fun/utils"
	"github.com/gempir/go-twitch-irc/v4"
)

type apiUser struct {
	ID    int    `json:"id"`
	Login string `json:"login"`
}

type apiStream struct {
	ID         int        `json:"id"`
	CreatedAt  time.Time  `json:"created_at"`
	EndedAt    *time.Time `json:"ended_at"`
	DurationMS *int       `json:"duration_ms"`
}

func init() {
	Fun.Register(&Cmd{
		Name: "vod",
		Handler: func(m twitch.PrivateMessage) (err error) {
			args := strings.Split(m.Message, " ")
			if len(args) < 1 || (args[0] != "`vod" && args[0] != "`vods") {
				return
			}

			c := m.Channel
			if len(args) >= 2 {
				c = strings.TrimPrefix(strings.ToLower(args[1]), "@")
			}

			res, err := api.Generic.Get("https://api-tv.supa.sh/user?login=" + url.QueryEscape(c))
			if err != nil || res.StatusCode != http.StatusOK {
				return
			}
			defer res.Body.Close()

			var user apiUser
			buf, _ := io.ReadAll(res.Body)
			err = json.Unmarshal(buf, &user)
			if err != nil {
				return
			}

			if args[0] == "`vods" {
				_, err = Say(m.RoomID, "https://tv.supa.sh/vods/"+user.Login, m.ID)
				return
			}

			res, err = api.Generic.Get("https://api-tv.supa.sh/streams?user_id=" + strconv.Itoa(user.ID))
			if err != nil || res.StatusCode != http.StatusOK {
				return
			}
			defer res.Body.Close()

			var streams []apiStream
			buf, _ = io.ReadAll(res.Body)
			err = json.Unmarshal(buf, &streams)
			if err != nil {
				return
			}

			if len(streams) == 0 {
				_, err = Say(m.RoomID, "no VODs available", m.ID)
				return
			}

			vod := streams[0]

			var msg string
			uptime := int(time.Since(vod.CreatedAt).Seconds())
			if vod.EndedAt == nil {
				msg += "(live) "
			} else {
				msg += fmt.Sprintf("%s ago, ", time.Since(*vod.EndedAt).Truncate(time.Second))
			}

			if vod.DurationMS == nil {
				msg += utils.FormatDuration(uptime)
			} else {
				msg += utils.FormatDuration(*vod.DurationMS / 1000)
			}

			msg += fmt.Sprintf(" https://tv.supa.sh/vods/%s/%v", user.Login, vod.ID)
			if uptime > 30 {
				msg += fmt.Sprintf("?t=%v", uptime-30)
			}

			_, err = Say(m.RoomID, msg, m.ID)
			return
		},
	})
}
