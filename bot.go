package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

type TBot struct {
	// Telegram bot struct
	Token  string `json:"token"`
	Debug  bool   `json:"debug"`
	Buffer int    `json:"buffer"`

	closeChannel chan interface{}
	self         User
	client       *http.Client
}

func NewBotApi(token string, client *http.Client) (*TBot, error) {
	bot := &TBot{
		Token:        token,
		client:       client,
		Buffer:       100,
		closeChannel: make(chan interface{}),
	}

	user, err := bot.getMe()

	if err != nil {
		return nil, err
	}

	bot.self = user

	return bot, nil
}

func (bot *TBot) MakeRequest(endpoint string, params url.Values) (BotResponse, error) {
	method := fmt.Sprintf(APIEndpoint, bot.Token, endpoint)
	resp, err := bot.client.PostForm(method, params)

	if err != nil {
		return BotResponse{}, err
	}

	defer resp.Body.Close()

	r := BotResponse{}
	dec := json.NewDecoder(resp.Body)
	err = dec.Decode(&r)

	if err != nil {
		return BotResponse{}, err
	}

	return r, nil
}

func (bot *TBot) getMe() (User, error) {
	res, err := bot.MakeRequest("getMe", nil)

	if err != nil {
		return User{}, err
	}

	var user User
	err = json.Unmarshal(res.Result, &user)

	if err != nil {
		return User{}, err
	}

	return user, nil
}

func (bot *TBot) debugLog(context string, params url.Values, message interface{}) {
	if bot.Debug {
		log.Printf("%s req : %+v\n", context, params)
		log.Printf("%s resp: %+v\n", context, message)
	}
}

func (bot *TBot) getUpdate(config UpdateConfig) ([]Update, error) {
	params := url.Values{}
	if config.limit > 0 {
		params.Add("limit", strconv.Itoa(config.limit))
	} else {
		params.Add("limit", strconv.Itoa(10))
	}
	if config.offset != 0 {
		params.Add("offset", strconv.Itoa(config.offset))
	} else {
		params.Add("offset", strconv.Itoa(0))
	}
	if config.timeout > 0 {
		params.Add("timeout", strconv.Itoa(config.timeout))
	} else {
		params.Add("timeout", strconv.Itoa(30))
	}

	res, err := bot.MakeRequest("getUpdates", params)
	if err != nil {
		return []Update{}, err
	}
	var u []Update
	err = json.Unmarshal(res.Result, &u)
	if err != nil {
		return []Update{}, err
	}
	return u, nil
}

func (bot *TBot) getUpdatesChan(config *UpdateConfig) (chan Update, error) {
	updates := make(chan Update, bot.Buffer)

	go func() {
		for {

			select {
			case <-bot.closeChannel:
				return
			default:
			}
			fmt.Print("Getting updates...\n")
			fmt.Printf("Config: %v\n", config)
			upds, err := bot.getUpdate(*config)
			if err != nil {
				fmt.Printf("Ooooops something wrong... Wait few seconds\nOriginal Error: %v\n", err)
				time.Sleep(time.Second * 3)
				continue
			}
			for _, u := range upds {
				if u.UpdateID != config.offset {
					config.offset = u.UpdateID
					updates <- u
				} else {
					fmt.Print("No updates, sleeping 3 seconds...\n")
					time.Sleep(time.Second * 3)
				}
			}
		}
	}()

	return updates, nil
}

func (bot *TBot) gettingUpdates(config UpdateConfig) {
	updates, _ := bot.getUpdatesChan(&config)
	for u := range updates {
		u.Print()
	}
}
