package main

import (
	"fmt"
	"log"
	"net/http"
	_ "net/http"
	_ "net/http/pprof"
)

func main() {
	// tr := &http.Transport{
	// 	MaxIdleConns:       10,
	// 	IdleConnTimeout:    30 * time.Second,
	// 	DisableCompression: true,
	// }

	// client := &http.Client{
	// 	Transport: tr,
	// }

	// token := flag.String("token", "", "Telegram Bot Token")
	// flag.Parse()
	// if *token == "" {
	// 	log.Fatal("Token must be provided!\n")
	// }
	// bot, err := NewBotApi(*token, client)
	// if err != nil {
	// 	fmt.Printf("%v\n", err)
	// }
	// fmt.Printf("Token: %v\nUser: %v\n", bot.Token, bot.self)
	// config := UpdateConfig{
	// 	limit: 5,
	// }

	// bot.gettingUpdates(config)
	testTrReq()

}

func testTrReq() {

	tr := NewTrClient()
	TransmissionRPCServerName = "127.0.0.1"
	params := TrRequest{
		Method: "session-stats",
	}
	res, err := tr.makeRequest(params)
	if err != nil {
		fmt.Printf("Err -> %s\n", err)
	}
	fmt.Printf("Res -> %s\n", res.Result)
	fmt.Printf("ARGS -> %s\n", res.Arguments)

}

func pprofile() {
	go func() {
		log.Println(http.ListenAndServe("localhost:6060", nil))
	}()
}
