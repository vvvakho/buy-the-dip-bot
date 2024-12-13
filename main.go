package main

import (
	"buy-the-dip-bot/telegram"
)

func main() {
	// telegram client
	telegram.InitBot()

	for {
		telegram.ListenForUpdates()
	}

	// database
	//
	// goroutine for api logic
	//
	// routine for listening for telegram

}
