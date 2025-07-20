package fun

import (
	"strings"

	. "github.com/0supa/supa8/fun"
	. "github.com/0supa/supa8/fun/api/cloudflare"
	api_kappa "github.com/0supa/supa8/fun/api/kappa"
	. "github.com/0supa/supa8/fun/api/twitch"
	"github.com/gempir/go-twitch-irc/v4"
)

func init() {
	Fun.Register(&Cmd{
		Name: "stable_diffusion",
		Handler: func(m twitch.PrivateMessage) (err error) {
			args := strings.Split(m.Message, " ")
			if (args[0] != "`sd") || len(args) < 2 {
				return
			}

			Say(m.RoomID, "ppHop", m.ID)

			body, err := StableDiffusionImage(strings.Join(args[1:], " "))
			if err != nil {
				return
			}

			upload, err := api_kappa.UploadFile(body, "image.png", "image/png")
			if err != nil {
				return
			}

			if upload.Checksum == "9b00921685afb8cc77218cdf39ce78c2" { // blank image
				_, err = Say(m.RoomID, "prompt rejected", m.ID)
				return
			}

			_, err = Say(m.RoomID, upload.Link, m.ID)
			return
		},
	})
}
