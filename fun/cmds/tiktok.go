package fun

import (
	"fmt"
	"log"
	"net/url"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	. "github.com/0supa/supa8/fun"
	. "github.com/0supa/supa8/fun/api/twitch"
	"github.com/0supa/supa8/fun/utils"
	"github.com/gempir/go-twitch-irc/v4"
	regexp "github.com/wasilibs/go-re2"
)

var cooldown = map[string]struct{}{}

func init() {
	links := regexp.MustCompile(`(?i)\S*tiktok\.com\/\S+|\S*(instagram|facebook)\.com\/(reels?|p|share)\/\S+`)
	parentDir := "/var/www/fi.supa.sh/tiktok"

	urlFix := strings.NewReplacer("/reels/", "/reel/", "/share/reel/", "/share/")

	Fun.Register(&Cmd{
		Name: "tiktok",
		Handler: func(m twitch.PrivateMessage) (err error) {
			if utils.IsBot(m.User.ID) {
				return
			}

			link := urlFix.Replace(links.FindString(m.Message))
			if link == "" {
				return
			}

			if _, found := cooldown[m.User.ID]; found && m.Channel != "omuljake" {
				return
			}

			cooldown[m.User.ID] = struct{}{}

			defer func() {
				time.Sleep(10 * time.Second)
				delete(cooldown, m.User.ID)
			}()

			cmd := exec.Command("./bin/yt-dlp",
				"--force-ipv4",
				"-S", "vcodec:h264",
				"--max-filesize", "100M",
				"--match-filters", "!is_live & !was_live",
				"--write-info-json",
				"--embed-metadata",
				"-P", fmt.Sprintf("%s/%s", parentDir, m.User.Name),
				"-o", fmt.Sprintf("%v.%%(ext)s", time.Now().Unix()),
				"--restrict-filenames",
				"-q", "--exec", "echo {}",
				link,
			)
			out, err := cmd.Output()
			if err != nil {
				if exit, ok := err.(*exec.ExitError); ok && exit.ExitCode() == 1 {
					log.Printf("yt-dlp error: %s:\n%s\n", err.Error(), exit.Stderr)

					msg := string(exit.Stderr)
					if strings.Contains(msg, "Restricted Video") || strings.Contains(msg, "This post may not be comfortable") {
						Say(m.RoomID, "🔞 Restricted Video", m.ID)
					}

					return nil
				}
				return err
			}

			fileName := filepath.Base(strings.TrimSuffix(string(out), "\n"))

			if fileName == "." {
				return
			}

			_, err = Say(
				m.RoomID,
				fmt.Sprintf("mirror: https://fi.supa.sh/tiktok/%s/%s", m.User.Name, url.PathEscape(fileName)),
				m.ID,
			)

			return err
		},
	})
}
