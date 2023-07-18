package main

import (
	"fmt"
	"github.com/gocolly/colly/v2"
	tb "gopkg.in/tucnak/telebot.v2"
	"log"
	"net/http"
	"os"
	"time"
	"io/ioutil"
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

func getSpots (string, error) {
	req, err := http.NewRequest("GET", "https://online.mfa.gov.ua/api/v1/queue/consulates/52/schedule?date=2023-06-23&dateEnd=2023-06-23&serviceId=530", nil)
	if err != nil {
		return "", err
	}

	req.Header.Set("authority", "online.mfa.gov.ua")
	req.Header.Set("accept", "application/json, text/plain, */*")
	req.Header.Set("accept-language", "en-GB,en-US;q=0.9,en;q=0.8")
	req.Header.Set("authorization", "Bearer YOUR_TOKEN_HERE")
	req.Header.Set("cache-control", "no-cache")
	req.Header.Set("cookie", "YOUR_COOKIE_HERE")
	req.Header.Set("pragma", "no-cache")
	req.Header.Set("referer", "https://online.mfa.gov.ua/application")
	req.Header.Set("sec-ch-ua", `"Not.A/Brand";v="8", "Chromium";v="114", "Google Chrome";v="114"`)
	req.Header.Set("sec-ch-ua-mobile", "?0")
	req.Header.Set("sec-ch-ua-platform", "macOS")
	req.Header.Set("sec-fetch-dest", "empty")
	req.Header.Set("sec-fetch-mode", "cors")
	req.Header.Set("sec-fetch-site", "same-origin")
	req.Header.Set("user-agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/114.0.0.0 Safari/537.36")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(body), nil
}

func sendUSDPriceDaily(b *tb.Bot) {
	ticker := time.NewTicker(24 * time.Hour)
	for range ticker.C {
		usdPrice, err := getUSDPrice()
		if err != nil {
			log.Println("Error getting USD price:", err)
			continue
		}

		chatID := tb.ChatID(381466119)
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
