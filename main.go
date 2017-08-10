package main

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/nlopes/slack"
)

type Log struct {
	Days [][]Message
}

type Message struct {
	Text      string
	Timestamp time.Time
}

func main() {
	log := Log{Days: make([][]Message, 7, 7)}

	token := os.Getenv("SLACK_API_TOKEN")
	client := slack.New(token)

	channel := os.Getenv("SLACK_CHANNEL")
	params := slack.NewHistoryParameters()
	hist, err := client.GetChannelHistory(channel, params)
	if err != nil {
		fmt.Println(err)
		return
	}

	now := time.Now()
	day := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())

	d := 0
	for _, m := range hist.Messages {
		unixtime, err := strconv.ParseFloat(m.Msg.Timestamp, 64)
		if err != nil {
			fmt.Println(err)
			return
		}

		timestamp := time.Unix(int64(unixtime), 0)

		for day.Sub(timestamp) > time.Duration(0) {
			day = day.Add(time.Duration(-24) * time.Hour)
			d++
		}

		if d > 6 {
			break
		}

		if log.Days[d] == nil {
			log.Days[d] = []Message{}
		}

		nm := Message{
			Text:      m.Msg.Text,
			Timestamp: timestamp,
		}
		log.Days[d] = append(log.Days[d], nm)
	}

	day = time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	day = day.Add(time.Duration(-7*24) * time.Hour)
	for i := len(log.Days) - 1; i >= 0; i-- {
		day = day.Add(24 * time.Hour)
		fmt.Println("##", day.Format("2006/01/02 (Mon)"))

		if len(log.Days[i]) == 0 {
			continue
		}

		for j := len(log.Days[i]) - 1; j >= 0; j-- {
			fmt.Println(log.Days[i][j].Text)
			fmt.Println("")
		}
		fmt.Println("")
	}
}
