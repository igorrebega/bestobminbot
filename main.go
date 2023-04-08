package main

import (
	"fmt"
	"github.com/gocolly/colly/v2"
	tb "gopkg.in/tucnak/telebot.v2"
	"log"
	"net/http"
	"os"
	"time"
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

func sendUSDPriceDaily(b *tb.Bot) {
	ticker := time.NewTicker(24 * time.Hour)
	for range ticker.C {
		usdPrice, err := getUSDPrice()
		if err != nil {
			log.Println("Error getting USD price:", err)
			continue
		}

		// Replace "YOUR_CHAT_ID" with your actual chat ID, which you can obtain by sending a message to your bot and checking the logs.
		chatID := tb.ChatID(123)
		b.Send(chatID, fmt.Sprintf("Daily USD price: %s", usdPrice))
	}
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

	b.Handle("/start", func(m *tb.Message) {
		helpText := `Welcome to the USD Price bot!
Here's a list of available commands:

/buyusd - Get the current USD price

If you have any questions, feel free to ask.`
		b.Send(m.Sender, helpText)
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

	b.Handle("/mychatid", func(m *tb.Message) {
		chatID := m.Chat.ID
		b.Send(m.Sender, fmt.Sprintf("Your chat ID is: %d", chatID))
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "3000" // Default port if not specified
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello, I'm your Telegram bot!")
	})

	go func() {
		log.Printf("Starting HTTP server on port %s", port)
		if err := http.ListenAndServe(":"+port, nil); err != nil {
			log.Fatalf("Error starting HTTP server: %v", err)
		}
	}()

	go sendUSDPriceDaily(b)
	log.Printf("Bot started on @%s", b.Me.Username)

	b.Start()
}
