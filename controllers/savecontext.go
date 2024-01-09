package controllers

import (
	"github.com/lalathealter/telegospel/keys"
	tele "gopkg.in/telebot.v3"
)

func ManageUserContext(next tele.HandlerFunc) tele.HandlerFunc {
	return func(ctx tele.Context) error {
		ctx = getPreviousData(ctx)
		defer saveResult(ctx)
		return next(ctx)
	}
}

var userContexts = map[int64]tele.Context{}

func getPreviousData(c tele.Context) tele.Context {
	id := getSpecificID(c)
	saved, ok := userContexts[id]
	if !ok {
		saved = c
	}
	return transferSettings(saved, c)
}

func transferSettings(from tele.Context, to tele.Context) tele.Context {
	for _, k := range keys.SETTINGS_KEYS {
		to.Set(k, from.Get(k))
	}
	return to
}

func saveResult(c tele.Context) {
	id := getSpecificID(c)
	userContexts[id] = c
}

func getSpecificID(c tele.Context) int64 {
	return c.Chat().ID
}
