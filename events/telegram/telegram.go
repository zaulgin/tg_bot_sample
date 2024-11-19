package telegram

import (
	"errors"

	"github.com/zaulgin/tg_bot_sample/clients/telegram"
	"github.com/zaulgin/tg_bot_sample/events"
	"github.com/zaulgin/tg_bot_sample/lib/e"
	"github.com/zaulgin/tg_bot_sample/storage"
)

type Processor struct {
	tg      *telegram.Client
	offset  int
	storage storage.Storage
}

type Meta struct {
	ChatID   int
	UserName string
}

func New(client *telegram.Client, storage storage.Storage) *Processor {
	return &Processor{
		tg:      client,
		storage: storage,
	}
}

var (
	ErrUnknownEventType = errors.New("unknown event type")
	ErrUnknownMetaType  = errors.New("unknown meta type")
)

func (p *Processor) Fetch(limit int) ([]events.Event, error) {
	updates, err := p.tg.Updates(p.offset, limit)
	if err != nil {
		return nil, e.Wrap("cant get events", err)
	}

	if len(updates) == 0 {
		return nil, nil
	}

	res := make([]events.Event, 0, len(updates))

	for _, u := range updates {
		res = append(res, event(u))
	}

	p.offset = updates[len(updates)-1].ID + 1

	return res, nil
}

func (p *Processor) Proccess(event events.Event) error {
	switch event.Type {
	case events.Message:
		return p.proccessMessage(event)
	default:
		return e.Wrap("cant process message", ErrUnknownEventType)
	}
}

func (p *Processor) proccessMessage(event events.Event) error {
	meta, err := meta(event)
	if err != nil {
		return e.Wrap("cant process message", err)
	}

	if err := p.doCmd(event.Text, meta.ChatID, meta.UserName); err != nil {
		return e.Wrap("cant process message", err)
	}

	return nil
}

func meta(event events.Event) (Meta, error) {
	res, ok := event.Meta.(Meta)
	if !ok {
		return Meta{}, e.Wrap("cant get meta", ErrUnknownMetaType)
	}

	return res, nil
}

func event(upd telegram.Update) events.Event {
	updType := fetchType(upd)

	res := events.Event{
		Type: updType,
		Text: fetchText(upd),
	}

	if updType == events.Message {
		res.Meta = Meta{
			ChatID:   upd.Message.Chat.ID,
			UserName: upd.Message.From.UserName,
		}
	}

	return res
}

func fetchType(upd telegram.Update) events.Type {
	if upd.Message == nil {
		return events.Unknown
	}
	return events.Message
}

func fetchText(upd telegram.Update) string {
	if upd.Message == nil {
		return ""
	}
	return upd.Message.Text
}
