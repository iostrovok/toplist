package queue

import (
	"context"
)

type Handler struct {
	Ch chan *Item
}

func NewHandler(ctx context.Context, save, insert SaveFunc, del DeleteFunc, build func()) *Handler {
	ch := make(chan *Item, 100)

	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case item, ok := <-ch:
				if !ok {
					return
				}

				item.Run(save, insert, del)
				build()
			}
		}
	}()

	return &Handler{
		Ch: ch,
	}
}

func (h *Handler) Run(item *Item) {
	h.Ch <- item
}
