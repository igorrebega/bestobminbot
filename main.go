package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/gocolly/colly/v2"
	tb "gopkg.in/tucnak/telebot.v2"
)

func getUSDPrice() (string, error) {
	var usdPrice string

	c := colly.NewCollector()

	c.OnHTML("#opt > div:nth-child(1) > div:nth-child(4) > div > div > p", func(e *colly.HTMLElement) {
		usdPrice = e.Text
	})

	err := c.Visit("https://www.bestobmin.com.ua/")
	if err != nil {
		return "", err
	}

	return usdPrice, nil
}

func main() {
	token := os.Getenv("TELEGRAM_BOT_TOKEN")
	if token == "" {
		log.Fatal("TELEGRAM_BOT_TOKEN is required")
	}

	b, err := tb.NewBot(tb.Settings{
		Token:  token,
		Poller: &tb.LongPoller{Timeout: 10 * time.Second},
	})

	if err != nil {
		log.Fatal(err)
		return
	}

	b.Handle("/hello", func(m *tb.Message) {
		b.Send(m.Sender, "Hello, I'm your Telegram bot!")
	})

	b.Handle("/buyusd", func(m *tb.Message) {
		usdPrice, err := getUSDPrice()
		if err != nil {
			b.Send(m.Sender, "Sorry, I could not fetch the current USD price.")
			log.Println("Error getting USD price:", err)
		} else {
			b.Send(m.Sender, fmt.Sprintf("Current USD price: %s", usdPrice))
		}
	})

	log.Printf("Bot started on @%s", b.Me.Username)
	b.Start()
}
