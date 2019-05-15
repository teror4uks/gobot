package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"time"
)

func main() {
	tr := &http.Transport{
		MaxIdleConns:       10,
		IdleConnTimeout:    30 * time.Second,
		DisableCompression: true,
	}

	client := &http.Client{
		Transport: tr,
	}

	token := flag.String("token", "", "Telegram Bot Token")
	flag.Parse()
	if *token == "" {
		log.Fatal("Token must be provided!\n")
	}
	bot, err := NewBotApi(*token, client)
	if err != nil {
		fmt.Printf("%v\n", err)
	}
	fmt.Printf("Token: %v\nUser: %v\n", bot.Token, bot.self)
	res, err := bot.getUpdate(0, 100, 30)
	if err != nil {
		log.Fatalf("Get Update method return err: %v\n", err)
	}
	fmt.Printf("Update result: %v\n", res)
}
