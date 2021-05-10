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
	DistrictID   = "392"
	NumDays      = 2
	AgeGroup     = 18
	BotToken     = "YOUR TELEGRAM BOT TOKEN" // Enter Bot Token
	MessageID    = 12345                     //Enter messageId
)

type Response struct {
	Sessions []Session `json:"sessions"`
}

type Session struct {
	SessionID         string `json:"session_id"`
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

func main() {
	bot, err := tgbotapi.NewBotAPI(BotToken)
	if err != nil {
		log.Panic(err)
	}

	bot.Debug = true

	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	sessionSet := make(map[string]bool)

	for {
		sessions := getSlots()
		for _, session := range sessions {
			exists := sessionSet[session.SessionID]
			if !exists {
				msgString := fmt.Sprintf("Date: %s\nCenter Name: %s\nPincode: %d\nFee Type: %s\nAvailable Capacity: %d\nVaccine: %s\n", session.Date, session.Name, session.Pincode, session.FeeType, session.AvailableCapacity, session.Vaccine)
				msg := tgbotapi.NewMessage(MessageID, msgString)
				bot.Send(msg)
				sessionSet[session.SessionID] = true
			}
		}
		time.Sleep(TimeInterval)
	}
}
