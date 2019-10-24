package main

import (
	"fmt"
	"github.com/nlopes/slack"
	"github.com/yanzay/log"
	"strings"
	"time"
)

type Bot struct {
	api *slack.Client
	rtm *slack.RTM
}

func NewBot(token string) Bot {
	var api = slack.New(token)
	return Bot{
		api,
		api.NewRTM()}
}

func (b *Bot) StartListening() {
	log.Info("Start listening for messages")

	go b.rtm.ManageConnection()

	for {
		select {
		case msg := <-b.rtm.IncomingEvents:
			switch ev := msg.Data.(type) {
			case *slack.MessageEvent:
				log.Info("Message received")

				info := b.rtm.GetInfo()
				prefix := fmt.Sprintf("<@%s> ", info.User.ID)

				if ev.User != info.User.ID && strings.HasPrefix(ev.Text, prefix) {
					b.respond(ev)
				}

			case *slack.RTMError:
				fmt.Printf("Error: %s\n", ev.Error())

			case *slack.InvalidAuthEvent:
				fmt.Printf("Invalid credentials")
			}
		}
	}
}

func (b *Bot) sendCoffeeMessage(channel string, user string) {
	time.AfterFunc(BREAKDOWN_TIME, func() {
		msg := fmt.Sprintf("<@%s> Time to get coffee!", user)
		_, _, err := b.api.PostMessage(channel, slack.MsgOptionText(msg, false))
		if err != nil {
			log.Error(err)
		}
	});
}

func (b *Bot) respond(msg *slack.MessageEvent) {
	go b.sendCoffeeMessage(msg.Channel, msg.User)
	b.rtm.SendMessage(b.rtm.NewOutgoingMessage(MSG_ADDED_COFFEE, msg.Channel))
}
