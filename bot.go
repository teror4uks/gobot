package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
)

type TBot struct {
	Token  string `json:"token"`
	Debug  bool   `json:"debug"`
	Buffer int    `json:"buffer"`

	self   User
	client *http.Client
}

func NewBotApi(token string, client *http.Client) (*TBot, error) {
	bot := &TBot{
		Token:  token,
		client: client,
		Buffer: 100,
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
