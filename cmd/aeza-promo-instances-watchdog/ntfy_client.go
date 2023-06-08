package main

import (
	"net/http"
	"net/url"
	"strings"
)

const (
	minPriority = "1"
	maxPriority = "5"
)

// https://docs.ntfy.sh/emojis/
var (
	mapPriorityTag = map[string]string{
		minPriority: "white_check_mark,rocket",
		maxPriority: "skull,warning",
	}
)

func sendNtfyIfEnabled(exitCode int, msg string) {
	if ntfyChannel == "" {
		return
	}
	if exitCode == passedCode && msg != "" {
		sendNtfyMessage(msg, minPriority)
	}
	if exitCode == errorCode {
		sendNtfyMessage(msg, maxPriority)
	}
}

func sendNtfyMessage(message string, priority string) {
	u, err := url.Parse(ntfyUrl)
	if err != nil {
		exit(1, "ERROR parse url ntfy: %s, for send message: %s", err.Error(), message)
	}
	u.Path = ntfyChannel

	req, err := http.NewRequest(http.MethodPost, u.String(), strings.NewReader(message))
	if err != nil {
		exit(1, "ERROR generate request: %s, for send message: %s", err.Error(), message)
	}
	req.Header.Set("Priority", priority)
	req.Header.Set("Tags", mapPriorityTag[priority])

	if _, err = http.DefaultClient.Do(req); err != nil {
		exit(1, "ERROR send request to ntfy: %s, for send message: %s", err.Error(), message)
	}
}
