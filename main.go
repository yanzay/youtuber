package main

import (
	"os"
	"os/exec"
	"strings"

	"github.com/yanzay/log"
	"github.com/yanzay/tbot"
)

func main() {
	token := os.Getenv("YOUTUBER_TOKEN")
	bot, err := tbot.NewServer(token)
	if err != nil {
		log.Fatal(err)
	}

	bot.HandleFunc("/youtube {link}", youtubeHandler, "Download video from youtube")
	bot.HandleFunc("/magnet {link}", magnetHandler, "Add magnet link to transmission")

	bot.AddMiddleware(tbot.NewAuth([]string{"yanzay", "katyabedryk"}))

	log.Fatal(bot.ListenAndServe())
}

func youtubeHandler(message tbot.Message) {
	link := message.Vars["link"]
	if !strings.Contains(link, "youtube.com") && !strings.Contains(link, "youtu.be") {
		message.Replyf("Error: '%s' is not a valid youtube link", link)
		return
	}
	cmd := exec.Command("/usr/local/bin/youtube-dl", "--no-mtime", "--restrict-filenames", "-o", "/home/osmc/downloads/%(title)s-%(id)s.%(ext)s", link)
	message.Replyf("Downloading %s", link)
	err := cmd.Run()
	if err != nil {
		message.Replyf("Error downloading %s: %s", link, err)
		return
	}
	message.Replyf("Download complete: %s", link)
}

// transmission-remote --add "magnet:?xt=urn:btih:13a27bdb9848a2e90514d5e6cc4e70f57e0ad11f&dn=Doctor+Strange+2016+HD-TS+x264+AC3-CPG&tr=udp%3A%2F%2Ftracker.leechers-paradise.org%3A6969&tr=udp%3A%2F%2Fzer0day.ch%3A1337&tr=udp%3A%2F%2Fopen.demonii.com%3A1337&tr=udp%3A%2F%2Ftracker.coppersurfer.tk%3A6969&tr=udp%3A%2F%2Fexodus.desync.com%3A6969" --auth "transmission:transmission"

func magnetHandler(message tbot.Message) {
	link := message.Vars["link"]
	if !strings.HasPrefix(link, "magnet:") {
		message.Replyf("Error: '%s' is not a valid magnet link", link)
		return
	}
	cmd := exec.Command("/usr/bin/transmission-remote", "--add", link, "--auth", "transmission:transmission")
	out, err := cmd.CombinedOutput()
	if err != nil {
		message.Replyf("Error adding magnet link: %s", err)
		message.Replyf("Details:\n%s", string(out))
		return
	}
	message.Reply("Magnet link added")
}
