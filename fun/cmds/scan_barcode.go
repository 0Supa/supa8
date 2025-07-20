package fun

import (
	"log"
	"net/http"
	"net/url"
	"os/exec"
	"strings"

	. "github.com/0supa/supa8/fun"
	"github.com/0supa/supa8/fun/api"
	. "github.com/0supa/supa8/fun/api/twitch"
	"github.com/0supa/supa8/fun/utils"
	"github.com/gempir/go-twitch-irc/v4"
	regexp "github.com/wasilibs/go-re2"
)

func init() {
	links := regexp.MustCompile(`(?i)\S*(i\.supa\.codes|kappa\.lol|gachi\.gay|femboy\.beauty)\/(\S+)`)

	Fun.Register(&Cmd{
		Name: "scan_barcode",
		Handler: func(m twitch.PrivateMessage) (err error) {
			if utils.IsBot(m.User.ID) {
				return
			}

			match := links.FindStringSubmatch(m.Message)
			if len(match) < 3 {
				return
			}

			id := match[2]
			fileURL := "https://kappa.lol/" + url.PathEscape(id)
			res, err := api.Generic.Head(fileURL)
			res.Body.Close()
			if err != nil || res.StatusCode != http.StatusOK || !strings.HasPrefix(res.Header.Get("Content-Type"), "image/") {
				log.Println(err)
				return nil
			}

			res, err = api.Generic.Get(fileURL)
			if err != nil {
				log.Println(err)
				return nil
			}

			cmd := exec.Command("zbarimg", "-q", "-")
			cmd.Stdin = res.Body
			out, err := cmd.Output()
			res.Body.Close()
			if err != nil {
				if _, ok := err.(*exec.ExitError); ok {
					return nil
				}
				return
			}

			for _, barcode := range strings.Split(string(out), "\n") {
				dat := strings.SplitN(barcode, ":", 2)
				if len(dat) < 2 || (!strings.HasPrefix(dat[0], "EAN") && !strings.HasPrefix(dat[0], "UPC")) {
					continue
				}

				foodfactsLink := "https://world.openfoodfacts.org/product/" + url.PathEscape(dat[1])
				res, _err := api.Generic.Head(foodfactsLink)
				if _err == nil && res.StatusCode == http.StatusOK {
					_, err = Say(m.RoomID, foodfactsLink, m.ID)
				}
			}

			return
		},
	})
}
