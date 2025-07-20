package fun

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	. "github.com/0supa/supa8/fun"
	logs_db "github.com/0supa/supa8/fun/api/clickhouse_db"
	. "github.com/0supa/supa8/fun/api/twitch"
	"github.com/gempir/go-twitch-irc/v4"
)

func init() {
	Fun.Register(&Cmd{
		Name: "find",
		Handler: func(m twitch.PrivateMessage) (err error) {
			args := strings.Split(m.Message, " ")
			if len(args) < 1 || args[0] != "`find" {
				return
			}

			if len(args) < 2 {
				_, err := Say(m.RoomID, "must specify a match pattern", m.ID)
				return err
			}

			pattern := strings.Join(args[1:], " ")
			rows, err := logs_db.Clickhouse.Query(context.Background(), "SELECT timestamp, user_login, text FROM message_structured WHERE channel_id=? AND match(text, ?) ORDER BY timestamp DESC LIMIT 1000", m.RoomID, pattern)
			if err != nil {
				return
			}
			defer rows.Close()

			tt := time.Now().UnixMilli()
			filePath := fmt.Sprintf("/var/www/fi.supa.sh/dump/%v.txt", tt)
			f, err := os.Create(filePath)
			if err != nil {
				return
			}
			defer f.Close()

			fmt.Fprintf(f, "#%s logs.supa.codes dump - query: `%s`\n\n", m.Channel, pattern)

			var messageCount int
			var timestamp time.Time
			var user, text string

			for rows.Next() {
				if err := rows.Scan(&timestamp, &user, &text); err != nil {
					return err
				}
				fmt.Fprintf(f, "[%s] %s: %s\n", timestamp.Format("2006-01-02 15:04:05"), user, text)
				messageCount++
			}

			if err := rows.Err(); err != nil {
				return err
			}

			var _messageCount string
			if messageCount >= 1000 {
				_messageCount = ">999"
			} else {
				_messageCount = strconv.Itoa(messageCount)
			}

			_, err = Say(m.RoomID, fmt.Sprintf("found %s messages: %s", _messageCount, fmt.Sprintf("https://fi.supa.sh/dump/%v.txt", tt)), m.ID)
			return err
		},
	})
}
