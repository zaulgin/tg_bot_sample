package telegram

import "github.com/zaulgin/tg_bot_sample/clients/telegram"

type Processor struct {
	tg     *telegram.Client
	offset int
	//storage
}

func New(client *telegram.Client)
