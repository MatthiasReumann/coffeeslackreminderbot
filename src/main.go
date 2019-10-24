package main

import "os"

func main() {
	var token= os.Getenv("SLACK_TOKEN")
	bot := NewBot(token)
	bot.StartListening()
}