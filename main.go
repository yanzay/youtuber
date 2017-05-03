package main

import (
	"fmt"
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
		message.Reply(fmt.Sprintf("Error: '%s' is not a valid youtube link", link), tbot.DisablePreview)
		return
	}
	cmd := exec.Command("/usr/local/bin/youtube-dl", "--no-mtime", "--restrict-filenames", "-o", "/home/osmc/downloads/%(id)s-%(title)s.%(ext)s", link)
	message.Reply(fmt.Sprintf("Downloading %s", link), tbot.DisablePreview)
	err := cmd.Run()
	if err != nil {
		message.Reply(fmt.Sprintf("Error downloading %s: %s", link, err), tbot.DisablePreview)
		return
	}
	message.Reply(fmt.Sprintf("Download complete: %s", link), tbot.DisablePreview)
}

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
