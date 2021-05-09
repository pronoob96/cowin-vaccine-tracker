package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

const (
	TimeInterval = 6 * time.Second
	DistrictID   = "149"
	NumDays      = 2
	AgeGroup     = 18
	BotToken     = "YOUR TELEGRAM BOT TOKEN"
)

type Response struct {
	Sessions []Session `json:"sessions"`
}

type Session struct {
	CenterID          int32  `json:"center_id"`
	Name              string `json:"name"`
	Address           string `json:"address"`
	Date              string `json:"date"`
	AvailableCapacity int32  `json:"available_capacity"`
	MinAgeLimit       int32  `json:"min_age_limit"`
	Vaccine           string `json:"vaccine"`
	FeeType           string `json:"fee_type"`
	Pincode           int32  `json:"pincode"`
}

func getSlots() []Session {
	date := time.Now()
	var sessionList []Session
	for i := 0; i < NumDays; i++ {
		dateString := date.Format("02-01-2006")
		params := "district_id=" + url.QueryEscape(DistrictID) + "&" +
			"date=" + url.QueryEscape(dateString)
		path := fmt.Sprintf("https://cdn-api.co-vin.in/api/v2/appointment/sessions/public/findByDistrict?%s", params)

		client := &http.Client{}
		req, _ := http.NewRequest("GET", path, nil)
		req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 6.1; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/56.0.2924.76 Safari/537.36")

		response, err := client.Do(req)
		if err != nil {
			fmt.Print(err.Error())
			continue
		}
		var responseData Response

		err = json.NewDecoder(response.Body).Decode(&responseData)
		if err != nil {
			log.Fatal(err)
		}
		response.Body.Close()

		for _, session := range responseData.Sessions {
			if session.MinAgeLimit == AgeGroup {
				sessionList = append(sessionList, session)
			}
		}

		date = date.AddDate(0, 0, 1)
	}

	return sessionList
}

func notifSender(userSet map[int64]bool, bot *tgbotapi.BotAPI) {
	for {
		sessions := getSlots()

		for _, session := range sessions {
			msgString := fmt.Sprintf("Date: %s\nCenter Name: %s\nPincode: %d\nFee Type: %s\nAvailable Capacity: %d\nVaccine: %s\n", session.Date, session.Name, session.Pincode, session.FeeType, session.AvailableCapacity, session.Vaccine)
			for user := range userSet {
				log.Println(user)

				msg := tgbotapi.NewMessage(user, msgString)

				bot.Send(msg)
			}
		}
		time.Sleep(TimeInterval)
	}
}

func main() {

	bot, err := tgbotapi.NewBotAPI(BotToken)
	if err != nil {
		log.Panic(err)
	}

	bot.Debug = true

	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, err := bot.GetUpdatesChan(u)
	welcomeMsg := "Hi! I am a bot developed by Pranav Kumar. I will notify you if there is a slot available for vaccination for upcoming month in Thane for 18+. To stop getting updates, reply with /stop"
	goodbyeMsg := "Thank you for using the bot. For getting update reply with /start"

	userSet := make(map[int64]bool)

	go notifSender(userSet, bot)

	for update := range updates {
		if update.Message == nil { // ignore any non-Message Updates
			continue
		}

		if update.Message.Text == "/start" {
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, welcomeMsg)

			bot.Send(msg)
			exists := userSet[update.Message.Chat.ID]
			if !exists {
				userSet[update.Message.Chat.ID] = true
			}

			continue
		}

		if update.Message.Text == "/stop" {
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, goodbyeMsg)

			bot.Send(msg)
			exists := userSet[update.Message.Chat.ID]
			if exists {
				delete(userSet, update.Message.Chat.ID)
			}

			continue
		}

	}

	responseData := getSlots()
	fmt.Println(responseData)
}
