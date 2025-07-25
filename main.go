package main

import (
	"bytes"
	"encoding/json"
	"io"
	logger "log"
	"net/http"
	"os"
	"slices"
	"strings"
	"sync"
	"time"

	"github.com/0supa/supa8/config"
	_ "github.com/0supa/supa8/fun"
	_ "github.com/0supa/supa8/fun/cmds"
)

type LogChannels struct {
	Channels []struct {
		Name   string `json:"name"`
		UserID string `json:"userID"`
	} `json:"channels"`
}

type LiveChannels []struct {
	Name    string `json:"login"`
	UserID  string `json:"uid"`
	Viewers int    `json:"viewers"`
	Type    string `json:"type"`
}

type JoinPayload struct {
	Channels []string `json:"channels"`
}

var wg = &sync.WaitGroup{}

var httpClient = http.Client{Timeout: time.Minute}

var log = logger.New(os.Stdout, "System ", logger.LstdFlags)

func init() {
	func() {
		res, err := httpClient.Get("https://github.com/yt-dlp/yt-dlp/releases/latest/download/yt-dlp_linux")
		if err != nil {
			log.Println("failed downloading yt-dlp bin", err)
			return
		}
		defer res.Body.Close()

		out, err := os.Create("./bin/yt-dlp")
		if err != nil {
			log.Println("failed creating yt-dlp bin", err)
			return
		}
		defer out.Close()

		_, err = io.Copy(out, res.Body)
		if err != nil {
			log.Println("failed writing yt-dlp bin", err)
			return
		}

		err = os.Chmod("./bin/yt-dlp", 0755)
		if err != nil {
			log.Println("failed setting +x on yt-dlp bin", err)
			return
		}
	}()
}

func main() {
	wg.Add(1)

	go func() {
		for range time.Tick(time.Minute * 10) {
			res, err := httpClient.Get("https://logs.supa.codes/channels")
			if err != nil {
				log.Println(err)
				continue
			}

			rustlog := LogChannels{}
			err = json.NewDecoder(res.Body).Decode(&rustlog)
			if err != nil {
				log.Println(err)
				continue
			}
			res.Body.Close()

			ignored := []string{}
			for _, ch := range rustlog.Channels {
				ignored = append(ignored, ch.UserID)
			}

			res, err = httpClient.Get("https://api-tv.supa.sh/tags/ro")
			if err != nil {
				log.Println(err)
				continue
			}

			liveChannels := LiveChannels{}
			err = json.NewDecoder(res.Body).Decode(&liveChannels)
			if err != nil {
				log.Println(err)
				continue
			}
			res.Body.Close()

			var resMsg strings.Builder
			resMsg.WriteString("logs.supa.codes -> joining new channels:")

			joinPayload := JoinPayload{}
			for _, ch := range liveChannels {
				if ch.Viewers < 10 || slices.Contains(ignored, ch.UserID) {
					continue
				}
				resMsg.WriteString(" @" + ch.Name)
				joinPayload.Channels = append(joinPayload.Channels, ch.UserID)
			}
			if len(joinPayload.Channels) == 0 {
				continue
			}

			body, err := json.Marshal(joinPayload)
			if err != nil {
				log.Println(err)
				continue
			}

			req, _ := http.NewRequest(
				"POST", "https://logs.supa.codes/admin/channels",
				bytes.NewBuffer(body),
			)
			req.Header.Set("Content-Type", "application/json")
			req.Header.Set("X-Api-Key", config.Auth.Rustlog.Key)

			res, err = httpClient.Do(req)
			if err != nil {
				log.Println(err)
				continue
			}

			if res.StatusCode != http.StatusOK {
				log.Println("failed joining new rustlog channels")

				b, err := io.ReadAll(res.Body)
				log.Println(string(b))
				if err != nil {
					log.Println(err)
				}
				continue
			}
			res.Body.Close()

			log.Println(resMsg.String())
			// fun.Say("675052240", resMsg.String(), "")
		}
	}()

	wg.Wait()
}
