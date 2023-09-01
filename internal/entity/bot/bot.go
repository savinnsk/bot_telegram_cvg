package entity

type Bot struct {
	ApiToken string
	ChatIds  []int64
}

func NewBot(apiToken string, chatId []int64) (*Bot, error) {
	return &Bot{
		ApiToken: apiToken,
		ChatIds:  chatId,
	}, nil
}
