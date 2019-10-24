package main

import (
	"fmt"
	"github.com/nlopes/slack"
	"github.com/yanzay/log"
	"os"
	"strings"
	"time"
)

const MSG_ADDED_COFFEE = "Acknowledged! I will notify you when it's time to get some coffee again!"
const BREAKDOWN_TIME = time.Minute * 210

var token = os.Getenv("SLACK_TOKEN")
var api = slack.New(token)

func main() {
	rtm := api.NewRTM()
	go rtm.ManageConnection()

	for {
		select {
		case msg := <-rtm.IncomingEvents:
			switch ev := msg.Data.(type) {
			case *slack.MessageEvent:
				info := rtm.GetInfo()
				prefix := fmt.Sprintf("<@%s> ", info.User.ID)

				if ev.User != info.User.ID && strings.HasPrefix(ev.Text, prefix) {
					respond(rtm, ev)
				}

			case *slack.RTMError:
				fmt.Printf("Error: %s\n", ev.Error())

			case *slack.InvalidAuthEvent:
				fmt.Printf("Invalid credentials")
			}
		}
	}
}

func sendMessage(channel string, user string){
	time.AfterFunc(BREAKDOWN_TIME, func(){
		msg := fmt.Sprintf("<@%s> Time to get coffee!",user)
		_, _, err := api.PostMessage(channel, slack.MsgOptionText(msg, false))
		if err != nil{
			log.Error(err)
		}
	});
}

func respond(rtm *slack.RTM, msg *slack.MessageEvent) {
	go sendMessage(msg.Channel, msg.User)
	rtm.SendMessage(rtm.NewOutgoingMessage(MSG_ADDED_COFFEE, msg.Channel))
}
