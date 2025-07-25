package fun

import (
	"context"
	"fmt"
	"io"
	"strings"

	. "github.com/0supa/supa8/fun"
	logs_db "github.com/0supa/supa8/fun/api/clickhouse_db"
	api_kappa "github.com/0supa/supa8/fun/api/kappa"
	. "github.com/0supa/supa8/fun/api/twitch"
	"github.com/0supa/supa8/fun/utils"
	"github.com/gempir/go-twitch-irc/v4"
	"github.com/olekukonko/tablewriter"
)

func init() {
	Fun.Register(&Cmd{
		Name: "toplines",
		Handler: func(m twitch.PrivateMessage) (err error) {
			args := strings.Split(m.Message, " ")
			if !utils.IsPrivileged(m.User.ID) || len(args) < 2 || args[0] != "`toplines" {
				return
			}

			user, err := GetUser(args[1], "")
			if err != nil {
				Say(m.RoomID, "failed getting user: "+err.Error(), m.ID)
				return
			}

			if user.ID == "" {
				Say(m.RoomID, "user not found", m.ID)
				return
			}

			rows, err := logs_db.Clickhouse.Query(context.Background(),
				`
SELECT
    channel_id,
    anyLast(channel_login) AS channel_login,
    count() AS lines
FROM rustlog.message_structured
WHERE user_id = ? AND message_type = 1
GROUP BY channel_id
`, user.ID, user.ID)
			if err != nil {
				return
			}
			defer rows.Close()

			tableString := &strings.Builder{}
			table := tablewriter.NewWriter(tableString)
			table.SetHeader([]string{"cID", "Channel", "Lines"})

			var cID string
			var cLogin string
			var lineCount uint64
			var totalLines uint64
			for rows.Next() {
				if err := rows.Scan(&cID, &cLogin, &lineCount); err != nil {
					return err
				}
				totalLines += lineCount
				table.Append([]string{cID, cLogin, fmt.Sprintf("%v", lineCount)})
			}

			tableString.WriteString(fmt.Sprintf("Top lines by @%s, total: %v\n\n", user.Login, totalLines))
			table.Render()

			rc := io.NopCloser(strings.NewReader(tableString.String()))
			defer rc.Close()

			upload, err := api_kappa.UploadFile(rc, "dat.txt", "text/plain")

			_, err = Say(m.RoomID, upload.Link, m.ID)

			return err
		},
	})
}
