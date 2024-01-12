package controllers

import (
	"fmt"

	"github.com/lalathealter/telegospel/db"
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
		sets, err := db.SelectSettingsFor(id)
		if err != nil {
			return c
		}

		return stuffContextWith(c, sets)
	}

	return transferSettings(saved, c)
}

func transferSettings(from tele.Context, to tele.Context) tele.Context {
	for _, k := range keys.SETTINGS_KEYS {
		to.Set(k, from.Get(k))
	}
	return to
}

func stuffContextWith(target tele.Context, settings map[string]any) tele.Context {
	for _, k := range keys.SETTINGS_KEYS {
		target.Set(k, settings[k])
	}
	return target
}

func saveResult(c tele.Context) {
	id := getSpecificID(c)
	userContexts[id] = c
	err := saveSettingsFor(id, c)
	if err != nil {
		fmt.Println(err)
	}
}

func saveSettingsFor(id int64, c tele.Context) error {
	data := map[string]any{}
	for _, k := range keys.SETTINGS_KEYS {
		data[k] = c.Get(k)
	}

	err := db.InsertSettingsFor(id, data)
	return err
}

func getSpecificID(c tele.Context) int64 {
	return c.Chat().ID
}
