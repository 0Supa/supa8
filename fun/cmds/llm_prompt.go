package fun

import (
	"fmt"
	"strings"
	"time"

	. "github.com/0supa/supa8/fun"
	. "github.com/0supa/supa8/fun/api/cloudflare"
	. "github.com/0supa/supa8/fun/api/twitch"
	"github.com/gempir/go-twitch-irc/v4"
)

func systemPrompt(m twitch.PrivateMessage) string {
	return fmt.Sprintf(`The current time is: %s.
You are talking to "%s", in the channel "%s".
Try to be helpful and informative.
Match the user's tone, you can be informal.

You should:
- keep your response concise.
- not add opening or closing statements.
- not use markdown or special formatting.

Q:`,
		// prompt appended with query
		time.Now().Format(time.RFC1123),
		m.User.Name, m.Channel)
}

func init() {
	Fun.Register(&Cmd{
		Name: "llm",
		Handler: func(m twitch.PrivateMessage) (err error) {
			var model string
			var messages []QueryMessage

			args := strings.Split(m.Message, " ")
			if len(args) < 2 {
				return
			}

			if args[0] == "`deepseek" || args[0] == "`r1" {
				model = "@cf/deepseek-ai/deepseek-r1-distill-qwen-32b"
			} else if args[0] == "`ask" || args[0] == "`llm" {
				model = "@cf/meta/llama-3.3-70b-instruct-fp8-fast"
				messages = append(messages, QueryMessage{Role: "system", Content: systemPrompt(m)})
			} else {
				return
			}

			messages = append(messages, QueryMessage{Role: "user", Content: strings.Join(args[1:], " ")})

			query := TextQuery{
				Stream:   true,
				Messages: messages,
			}

			c := make(chan Result)
			go TextGeneration(c, query, model)

			var builder strings.Builder
			for data := range c {
				if err := data.Error; err != nil {
					return err
				}

				builder.WriteString(data.Response)
			}
			res := builder.String()

			_, err = Say(m.RoomID, res, m.ID)
			return
		},
	})
}
